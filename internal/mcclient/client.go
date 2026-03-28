package mcclient

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"gmcc/internal/auth/microsoft"
	mcauth "gmcc/internal/auth/minecraft"
	"gmcc/internal/config"
	"gmcc/internal/constants"
	"gmcc/internal/entity"
	"gmcc/internal/logx"
	"gmcc/internal/mcclient/packet"
	"gmcc/internal/mcclient/protocol"
	"gmcc/internal/player"
	"gmcc/internal/session"
)

var errOnlineAuthRequired = errors.New("online auth required")

type onlineSession struct {
	AccessToken string
	ProfileID   string
	ProfileName string
	ProfileUUID [16]byte
}

// Client 是协议 774 (MC 1.21.11) 的最小登录/挂机客户端。
type Client struct {
	cfg *config.Config

	offlineName string
	offlineUUID [16]byte

	username string
	uuid     [16]byte

	online *onlineSession

	state  protocol.State
	inPlay bool

	conn *packet.PacketConn

	lastAFKPacket time.Time
	chatHandler   func(ChatMessage)
	chatSessionOK bool
	chatSession   *secureChatSession
	commandSign   map[string]signableCommandTarget
	chatSignMu    sync.Mutex

	Player    *player.Player
	players   map[string]playerInfo
	playersMu sync.RWMutex

	// 实体跟踪
	entityTracker *entity.Tracker
	NearbyPlayers *player.NearbyTracker
}

type playerInfo struct {
	uuid [16]byte
}

func New(cfg *config.Config) *Client {
	name := strings.TrimSpace(cfg.Account.PlayerID)
	client := &Client{
		cfg:         cfg,
		offlineName: name,
		offlineUUID: packet.OfflineUUID(name),
		username:    name,
		uuid:        packet.OfflineUUID(name),
		commandSign: map[string]signableCommandTarget{},
		Player:      player.NewPlayer(),
		players:     make(map[string]playerInfo),
	}

	return client
}

func (c *Client) IsReady() bool {
	return c.inPlay
}

func (c *Client) Run(ctx context.Context) error {
	if strings.TrimSpace(c.offlineName) == "" {
		return fmt.Errorf("account.player_id 不能为空")
	}
	logx.Infof("客户端协议: %s", protocol.Label)
	host, port, err := packet.ParseAddress(c.cfg.Server.Address)
	if err != nil {
		return err
	}

	if c.cfg.Account.UseOfficialAuth {
		logx.Infof("已启用正版认证 (account.use_official_auth=true)")
		if err := c.prepareOnlineSession(); err != nil {
			return err
		}
		return c.connectAndLoop(ctx, host, port, true)
	}

	logx.Infof("使用离线模式 (account.use_official_auth=false)")
	err = c.connectAndLoop(ctx, host, port, false)
	if errors.Is(err, errOnlineAuthRequired) {
		return fmt.Errorf("服务器要求正版会话认证，请将 account.use_official_auth 设为 true")
	}
	return err
}

func (c *Client) connectAndLoop(ctx context.Context, host string, port uint16, useOnline bool) error {
	if useOnline {
		if c.online == nil {
			return fmt.Errorf("online session 不存在")
		}
		c.username = c.online.ProfileName
		c.uuid = c.online.ProfileUUID
	} else {
		c.username = c.offlineName
		c.uuid = c.offlineUUID
	}

	c.state = protocol.StateLogin
	c.inPlay = false
	c.chatSessionOK = false
	c.chatSession = nil
	c.commandSign = map[string]signableCommandTarget{}
	c.lastAFKPacket = time.Now()

	addr := net.JoinHostPort(host, strconv.Itoa(int(port)))
	dialer := net.Dialer{Timeout: constants.DialTimeout}
	rawConn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("连接服务器失败: %w", err)
	}

	c.conn = packet.NewPacketConn(rawConn)
	defer func() {
		if c.conn != nil {
			_ = c.conn.Close()
		}
	}()

	if useOnline {
		logx.Infof("已连接服务器: %s (online-mode)", addr)
	} else {
		logx.Infof("已连接服务器: %s (offline-mode attempt)", addr)
	}

	if err := c.sendHandshake(host, port); err != nil {
		return err
	}
	if err := c.sendLoginStart(); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		_ = c.conn.SetReadDeadline(time.Now().Add(constants.ReadTimeout))
		pkt, err := c.conn.ReadPacket()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				if c.state == protocol.StatePlay {
					if err := c.sendAFKHeartbeatIfNeeded(); err != nil {
						return err
					}
				}
				continue
			}
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				return fmt.Errorf("服务器已关闭连接 (state=%s): %w", c.state.String(), err)
			}
			return fmt.Errorf("读取数据包失败 (state=%s): %w", c.state.String(), err)
		}

		logx.PacketLogf(
			"收到数据包: state=%s id=0x%02X (%s) len=%d",
			c.state.String(),
			pkt.ID,
			protocol.PacketName(c.state, pkt.ID),
			len(pkt.Data),
		)
		if err := c.handlePacket(pkt); err != nil {
			return err
		}

		if c.state == protocol.StatePlay {
			if err := c.sendAFKHeartbeatIfNeeded(); err != nil {
				return err
			}
		}
	}
}

func (c *Client) prepareOnlineSession() error {
	cache, err := session.Load(c.offlineName)
	if err != nil {
		return err
	}

	now := time.Now()
	if cache.HasValidMinecraftToken(now) {
		if err := c.setOnlineSession(
			cache.Minecraft.AccessToken,
			cache.Minecraft.ProfileID,
			cache.Minecraft.ProfileName,
		); err == nil {
			logx.LogTokenCache("minecraft", c.online.ProfileName, packet.FormatUUID(c.online.ProfileUUID))
			return nil
		}
		logx.Warnf("检测 detections 缓存的 Minecraft token 数据损坏，准备重新认证")
	}

	xstsResp, msToken, err := c.resolveXSTSToken(cache, now)
	if err != nil {
		return fmt.Errorf("Microsoft 登录失败: %w", err)
	}

	mcToken, err := mcauth.GetMinecraftToken(xstsResp)
	if err != nil {
		return fmt.Errorf("获取 Minecraft Token 失败: %w", err)
	}
	if err := mcauth.VerifyGameOwnership(mcToken.AccessToken); err != nil {
		return err
	}

	profile, err := mcauth.GetProfile(mcToken.AccessToken)
	if err != nil {
		return err
	}
	if err := c.setOnlineSession(mcToken.AccessToken, profile.ID, profile.Name); err != nil {
		return err
	}

	if msToken != nil {
		cache.Microsoft.AccessToken = strings.TrimSpace(msToken.AccessToken)
		if refresh := strings.TrimSpace(msToken.RefreshToken); refresh != "" {
			cache.Microsoft.RefreshToken = refresh
		}
		cache.Microsoft.ExpiresAt = tokenExpiresAt(msToken.ExpiresIn)
	}
	cache.Minecraft.AccessToken = mcToken.AccessToken
	cache.Minecraft.ExpiresAt = tokenExpiresAt(mcToken.ExpiresIn)
	cache.Minecraft.ProfileID = c.online.ProfileID
	cache.Minecraft.ProfileName = c.online.ProfileName

	if err := session.Save(c.offlineName, cache); err != nil {
		logx.Warnf("写入 token 缓存失败: %v", err)
	}

	logx.Infof("正版认证成功: %s (%s)", c.online.ProfileName, packet.FormatUUID(c.online.ProfileUUID))
	return nil
}

func (c *Client) resolveXSTSToken(cache *session.TokenCache, now time.Time) (*microsoft.XSTSResponse, *microsoft.TokenResponse, error) {
	if cache != nil && cache.HasValidMicrosoftAccess(now) {
		xstsResp, err := microsoft.GetXSTSTokenFromAccessToken(cache.Microsoft.AccessToken)
		if err == nil {
			logx.LogTokenCache("microsoft", "", "")
			return xstsResp, nil, nil
		}
		logx.LogTokenExpired("microsoft", err)
	}

	if cache != nil && cache.HasMicrosoftRefreshToken() {
		msToken, err := microsoft.RefreshMicrosoftToken(cache.Microsoft.RefreshToken)
		if err == nil {
			xstsResp, xstsErr := microsoft.GetXSTSTokenFromMicrosoftToken(msToken)
			if xstsErr == nil {
				logx.Infof("已通过 refresh_token 刷新 Microsoft 令牌")
				return xstsResp, msToken, nil
			}
			logx.Warnf("刷新后换取 XSTS 失败，回退设备码授权: %v", xstsErr)
		} else {
			logx.Warnf("refresh_token 刷新失败，回退设备码授权: %v", err)
		}
	}

	msToken, err := microsoft.GetMicrosoftToken()
	if err != nil {
		return nil, nil, err
	}
	xstsResp, err := microsoft.GetXSTSTokenFromMicrosoftToken(msToken)
	if err != nil {
		return nil, nil, err
	}
	logx.Infof("Microsoft 设备码登录成功")
	return xstsResp, msToken, nil
}

func (c *Client) setOnlineSession(accessToken, profileID, profileName string) error {
	token := strings.TrimSpace(accessToken)
	cleanID := strings.ReplaceAll(strings.TrimSpace(profileID), "-", "")
	name := strings.TrimSpace(profileName)
	if token == "" {
		return fmt.Errorf("Minecraft access token 为空")
	}
	if cleanID == "" || name == "" {
		return fmt.Errorf("Minecraft profile 信息不完整")
	}

	profileUUID, err := packet.ParseUUID(cleanID)
	if err != nil {
		return fmt.Errorf("解析 profile UUID 失败: %w", err)
	}

	c.online = &onlineSession{
		AccessToken: token,
		ProfileID:   cleanID,
		ProfileName: name,
		ProfileUUID: profileUUID,
	}
	return nil
}

func tokenExpiresAt(expiresIn int) time.Time {
	if expiresIn <= 0 {
		return time.Time{}
	}
	return time.Now().Add(time.Duration(expiresIn) * time.Second).UTC()
}

// initializeTrackers 初始化实体和玩家跟踪器
func (c *Client) initializeTrackers() {
	// 创建实体跟踪器
	c.entityTracker = entity.NewTracker()

	// 创建玩家跟踪器
	c.NearbyPlayers = player.NewNearbyTracker(c.entityTracker, c.getPlayerInfoByUUID)
}

// getPlayerInfoByUUID 通过UUID查找玩家用户名
func (c *Client) getPlayerInfoByUUID(uuid [16]byte) (username string, found bool) {
	c.playersMu.RLock()
	defer c.playersMu.RUnlock()

	for name, info := range c.players {
		if info.uuid == uuid {
			return name, true
		}
	}
	return "", false
}

func (c *Client) handlePacket(pkt packet.Packet) error {
	switch c.state {
	case protocol.StateLogin:
		return c.handleLoginPacket(pkt)
	case protocol.StateConfiguration:
		return c.handleConfigurationPacket(pkt)
	case protocol.StatePlay:
		return c.handlePlayPacket(pkt)
	default:
		return fmt.Errorf("未知连接状态: %s", c.state.String())
	}
}

func (c *Client) SendPlayerRotation(yaw, pitch float32, onGround bool) error {
	if c.state != protocol.StatePlay {
		return fmt.Errorf("当前状态不是 Play，无法发送视角数据包")
	}
	if c.conn == nil {
		return fmt.Errorf("连接未初始化")
	}

	payload := make([]byte, 0, 9)
	payload = append(payload, packet.EncodeFloat32(yaw)...)
	payload = append(payload, packet.EncodeFloat32(pitch)...)
	payload = append(payload, packet.EncodeBool(onGround)...)

	return c.conn.WritePacket(protocol.PlayServerMovePlayerRot, payload)
}

// SendSetCarriedItem 切换快捷栏槽位 (0-8 对应快捷栏 1-9)
// 注: 包ID需要根据实际协议版本确认，目前代码暂未连接实际服务器
func (c *Client) SendSetCarriedItem(slot int16) error {
	if c.state != protocol.StatePlay {
		return fmt.Errorf("当前状态不是 Play，无法发送快捷栏数据包")
	}
	if c.conn == nil {
		return fmt.Errorf("连接未初始化")
	}
	if slot < 0 || slot > 8 {
		return fmt.Errorf("slot 必须在 0-8 范围内")
	}

	// Set Carried Item: slot (Short)
	// 使用 big-endian 编码
	payload := make([]byte, 0, 2)
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(slot))
	payload = append(payload, buf...)

	return c.conn.WritePacket(protocol.PlayServerSetCarriedItem, payload)
}

// SendInteract 发送实体交互包 (右键点击实体)
// entityID: 目标实体ID
// action: InteractActionInteract (0), InteractActionAttack (1), InteractActionInteractAt (2)
// hand: HandMainHand (0), HandOffHand (1)
// sneaking: 是否潜行
func (c *Client) SendInteract(entityID int32, action int32, hand int32, sneaking bool) error {
	if c.state != protocol.StatePlay {
		return fmt.Errorf("当前状态不是 Play，无法发送交互数据包")
	}
	if c.conn == nil {
		return fmt.Errorf("连接未初始化")
	}

	payload := make([]byte, 0, 16)
	payload = append(payload, packet.EncodeVarInt(entityID)...)
	payload = append(payload, packet.EncodeVarInt(action)...)

	// Interact (action=0) 和 InteractAt (action=2) 需要发送 hand
	if action == protocol.InteractActionInteract || action == protocol.InteractActionInteractAt {
		// 对于 InteractAt，还需要发送目标位置
		if action == protocol.InteractActionInteractAt {
			// INTERACT_AT 需要额外的 targetX/y/z (float32)
			// 我们使用简化的 INTERACT (action=0) 进行右键交互
			return fmt.Errorf("暂不支持 INTERACT_AT，请使用 INTERACT")
		}
		payload = append(payload, packet.EncodeVarInt(hand)...)
	}

	payload = append(payload, packet.EncodeBool(sneaking)...)

	return c.conn.WritePacket(protocol.PlayServerInteract, payload)
}
