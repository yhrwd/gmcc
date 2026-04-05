package mcclient

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	mcauth "gmcc/internal/auth/minecraft"
	authsession "gmcc/internal/auth/session"
	"gmcc/internal/config"
	"gmcc/internal/constants"
	"gmcc/internal/logx"
	"gmcc/internal/mcclient/packet"
	"gmcc/internal/mcclient/protocol"
)

type onlineSession struct {
	AccessToken string
	ProfileID   string
	ProfileName string
	ProfileUUID [16]byte
}

// Player 玩家状态
type Player struct {
	mu sync.RWMutex

	EntityID int32
	UUID     [16]byte
	Name     string

	X, Y, Z    float64
	Yaw, Pitch float32
	OnGround   bool

	// 状态
	Health     float32
	MaxHealth  float32
	Food       int32
	Saturation float32
	Level      int32
	Experience float32

	JoinTime time.Time
}

func NewPlayer() *Player {
	return &Player{
		JoinTime:   time.Now(),
		OnGround:   true,
		X:          0,
		Y:          64,
		Z:          0,
		Yaw:        0,
		Pitch:      0,
		Health:     20,
		MaxHealth:  20,
		Food:       20,
		Saturation: 5,
	}
}

func (p *Player) UpdatePosition(x, y, z float64, yaw, pitch float32, relative int8) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if relative&0x01 != 0 {
		x += p.X
	}
	if relative&0x02 != 0 {
		y += p.Y
	}
	if relative&0x04 != 0 {
		z += p.Z
	}
	if relative&0x08 != 0 {
		yaw += p.Yaw
	}
	if relative&0x10 != 0 {
		pitch += p.Pitch
	}
	p.X, p.Y, p.Z = x, y, z
	p.Yaw, p.Pitch = yaw, pitch
}

func (p *Player) GetPosition() (float64, float64, float64) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.X, p.Y, p.Z
}

func (p *Player) GetMovementState() (float64, float64, float64, float32, float32, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.X, p.Y, p.Z, p.Yaw, p.Pitch, p.OnGround
}

func (p *Player) SetEntityID(id int32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.EntityID = id
}

func (p *Player) SetUUID(uuid [16]byte) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.UUID = uuid
}

func (p *Player) SetName(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Name = name
}

func (p *Player) UpdateHealth(health, maxHealth float32, food int32, saturation float32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Health = health
	p.MaxHealth = maxHealth
	p.Food = food
	p.Saturation = saturation
}

func (p *Player) UpdateExperience(level int32, experience float32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Level = level
	p.Experience = experience
}

func (p *Player) GetHealth() (float32, float32, int32, float32) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Health, p.MaxHealth, p.Food, p.Saturation
}

func (p *Player) GetExperience() (int32, float32) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Level, p.Experience
}

// GetOnlineDuration 获取在线时长
func (p *Player) GetOnlineDuration() time.Duration {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return time.Since(p.JoinTime)
}

// Client 是协议 774 (MC 1.21.11) 的客户端
type Client struct {
	cfg *config.Config

	username string
	uuid     [16]byte

	online *onlineSession

	state  protocol.State
	inPlay bool

	conn *packet.PacketConn

	lastAFKPacket time.Time

	Player *Player

	accountID   string
	authManager *authsession.AuthManager
}

func New(cfg *config.Config) *Client {
	return &Client{
		cfg:         cfg,
		Player:      NewPlayer(),
		accountID:   strings.TrimSpace(cfg.ClusterRuntime.AccountID),
		authManager: cfg.ClusterRuntime.AuthManager,
	}
}

func (c *Client) IsReady() bool {
	return c.inPlay
}

func (c *Client) Run(ctx context.Context) error {
	logx.Infof("客户端协议: %s", protocol.Label)
	host, port, err := packet.ParseAddress(c.cfg.Server.Address)
	if err != nil {
		return err
	}

	if err := c.prepareOnlineSession(); err != nil {
		return err
	}
	return c.connectAndLoop(ctx, host, port, true)
}

func (c *Client) connectAndLoop(ctx context.Context, host string, port uint16, useOnline bool) error {
	if c.online == nil {
		return fmt.Errorf("online session 不存在")
	}
	c.username = c.online.ProfileName
	c.uuid = c.online.ProfileUUID

	c.state = protocol.StateLogin
	c.inPlay = false
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

func (c *Client) sendAFKHeartbeatIfNeeded() error {
	if time.Since(c.lastAFKPacket) < constants.AFKCheckInterval {
		return nil
	}
	c.lastAFKPacket = time.Now()

	x, y, z, _, _, onGround := c.Player.GetMovementState()
	flags := byte(0)
	if onGround {
		flags |= 0x01
	}

	payload := make([]byte, 0, 25)
	payload = append(payload, packet.EncodeFloat64(x)...)
	payload = append(payload, packet.EncodeFloat64(y)...)
	payload = append(payload, packet.EncodeFloat64(z)...)
	payload = append(payload, flags)

	if err := c.conn.WritePacket(protocol.PlayServerMovePlayerPos, payload); err != nil {
		return err
	}

	return c.conn.WritePacket(protocol.PlayServerClientTickEnd, nil)
}

func (c *Client) prepareOnlineSession() error {
	if c.authManager == nil {
		return fmt.Errorf("auth manager 未配置")
	}
	if c.accountID == "" {
		return fmt.Errorf("cluster.account_id 不能为空")
	}
	authSession, err := c.authManager.GetSession(context.Background(), c.accountID)
	if err != nil {
		return fmt.Errorf("获取账号认证会话失败: %w", err)
	}
	if err := c.setOnlineSession(authSession.MinecraftAccessToken, authSession.ProfileID, authSession.ProfileName); err != nil {
		return err
	}
	logx.Summaryf("info", "正版认证成功: %s (%s)", c.online.ProfileName, packet.FormatUUID(c.online.ProfileUUID))
	return nil
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

// SendCommand 发送命令（无签名）
func (c *Client) SendCommand(cmd string) error {
	cmd = strings.TrimSpace(cmd)
	cmd = strings.TrimPrefix(cmd, "/")
	if cmd == "" {
		return fmt.Errorf("command 不能为空")
	}
	if c.state != protocol.StatePlay {
		return fmt.Errorf("当前状态不是 Play，无法发送命令")
	}

	payload := packet.EncodeString(cmd)
	return c.conn.WritePacket(protocol.PlayServerChatCommand, payload)
}

// SendClientCommand 发送客户端命令 (0=重生)
func (c *Client) SendClientCommand(action int32) error {
	if c.state != protocol.StatePlay {
		return fmt.Errorf("当前状态不是 Play")
	}
	payload := packet.EncodeVarInt(action)
	return c.conn.WritePacket(protocol.PlayServerClientCommand, payload)
}

// handleLoginPacket 处理登录阶段的数据包
func (c *Client) handleLoginPacket(pkt packet.Packet) error {
	switch pkt.ID {
	case protocol.LoginClientDisconnect:
		return fmt.Errorf("登录阶段被服务器断开: %s", disconnectReasonFromJSON(pkt.Data))

	case protocol.LoginClientCompression:
		r := bytes.NewReader(pkt.Data)
		threshold, err := packet.ReadVarInt(r)
		if err != nil {
			return fmt.Errorf("读取压缩阈值失败: %w", err)
		}
		c.conn.SetCompressionThreshold(int(threshold))
		logx.Infof("启用压缩, threshold=%d", threshold)
		return nil

	case protocol.LoginClientHello:
		return c.handleEncryptionRequest(pkt.Data)

	case protocol.LoginClientFinished:
		profileUUID, profileName, err := parseGameProfile(pkt.Data)
		if err != nil {
			return fmt.Errorf("解析 login_finished 失败: %w", err)
		}
		c.Player.SetUUID(profileUUID)
		c.Player.SetName(profileName)
		c.username = profileName

		if profileName != "" {
			logx.Infof("登录成功: %s (%s)", profileName, packet.FormatUUID(profileUUID))
		} else {
			logx.Infof("登录成功: %s", packet.FormatUUID(profileUUID))
		}

		if err := c.conn.WritePacket(protocol.LoginServerAck, nil); err != nil {
			return fmt.Errorf("发送 login_acknowledged 失败: %w", err)
		}
		c.state = protocol.StateConfiguration
		if err := c.sendClientInformation(protocol.CfgServerClientInfo); err != nil {
			return fmt.Errorf("发送 configuration client_information 失败: %w", err)
		}
		return nil

	case protocol.LoginClientCustomQuery:
		r := bytes.NewReader(pkt.Data)
		txID, err := packet.ReadVarInt(r)
		if err != nil {
			return fmt.Errorf("读取 custom_query transactionId 失败: %w", err)
		}
		payload := make([]byte, 0, 8)
		payload = append(payload, packet.EncodeVarInt(txID)...)
		payload = append(payload, packet.EncodeBool(false)...)
		if err := c.conn.WritePacket(protocol.LoginServerCustomAnswer, payload); err != nil {
			return fmt.Errorf("发送 custom_query_answer 失败: %w", err)
		}
		return nil

	case protocol.LoginClientCookieReq:
		r := bytes.NewReader(pkt.Data)
		key, err := packet.ReadString(r, r)
		if err != nil {
			return fmt.Errorf("读取 login cookie key 失败: %w", err)
		}
		return c.sendCookieResponse(protocol.LoginServerCookieResp, key)

	default:
		logx.PacketLogf("未处理的 Login 数据包: id=0x%02X (%s) len=%d", pkt.ID, protocol.PacketName(protocol.StateLogin, pkt.ID), len(pkt.Data))
		return nil
	}
}

// handleConfigurationPacket 处理配置阶段的数据包
func (c *Client) handleConfigurationPacket(pkt packet.Packet) error {
	switch pkt.ID {
	case protocol.CfgClientDisconnect:
		return fmt.Errorf("配置阶段被服务器断开: %s", disconnectReasonFromNBT(pkt.Data))

	case protocol.CfgClientKeepAlive:
		r := bytes.NewReader(pkt.Data)
		id, err := packet.ReadInt64(r)
		if err != nil {
			return fmt.Errorf("读取配置阶段 keep_alive 失败: %w", err)
		}
		return c.conn.WritePacket(protocol.CfgServerKeepAlive, packet.EncodeInt64(id))

	case protocol.CfgClientPing:
		r := bytes.NewReader(pkt.Data)
		id, err := packet.ReadInt32(r)
		if err != nil {
			return fmt.Errorf("读取配置阶段 ping 失败: %w", err)
		}
		return c.conn.WritePacket(protocol.CfgServerPong, packet.EncodeInt32(id))

	case protocol.CfgClientCookieReq:
		r := bytes.NewReader(pkt.Data)
		key, err := packet.ReadString(r, r)
		if err != nil {
			return fmt.Errorf("读取配置阶段 cookie key 失败: %w", err)
		}
		return c.sendCookieResponse(protocol.CfgServerCookieResp, key)

	case protocol.CfgClientSelectPacks:
		payload := packet.EncodeVarInt(0)
		return c.conn.WritePacket(protocol.CfgServerSelectPacks, payload)

	case protocol.CfgClientPackPush:
		id, err := readPacketUUID(pkt.Data)
		if err != nil {
			return fmt.Errorf("读取配置阶段 resource_pack_push 失败: %w", err)
		}
		return c.sendResourcePackResponses(protocol.CfgServerResource, id)

	case protocol.CfgClientCodeOfConduct:
		return c.conn.WritePacket(protocol.CfgServerAcceptCode, nil)

	case protocol.CfgClientFinish:
		if err := c.conn.WritePacket(protocol.CfgServerFinish, nil); err != nil {
			return fmt.Errorf("发送 finish_configuration 失败: %w", err)
		}
		c.state = protocol.StatePlay
		if err := c.sendClientInformation(protocol.PlayServerClientInfo); err != nil {
			return fmt.Errorf("发送 play client_information 失败: %w", err)
		}
		logx.Infof("进入 Play 阶段")
		return nil

	case protocol.CfgClientCustom, protocol.CfgClientRegistry, protocol.CfgClientPackPop:
		return nil

	default:
		logx.PacketLogf("未处理的 Configuration 数据包: id=0x%02X (%s) len=%d", pkt.ID, protocol.PacketName(protocol.StateConfiguration, pkt.ID), len(pkt.Data))
		return nil
	}
}

// handlePlayPacket 处理游戏阶段的数据包
func (c *Client) handlePlayPacket(pkt packet.Packet) error {
	switch pkt.ID {
	case protocol.PlayClientDisconnect:
		return fmt.Errorf("Play 阶段被服务器断开: %s", disconnectReasonFromNBT(pkt.Data))

	case protocol.PlayClientKeepAlive:
		r := bytes.NewReader(pkt.Data)
		id, err := packet.ReadInt64(r)
		if err != nil {
			return fmt.Errorf("读取 play keep_alive 失败: %w", err)
		}
		return c.conn.WritePacket(protocol.PlayServerKeepAlive, packet.EncodeInt64(id))

	case protocol.PlayClientPing:
		r := bytes.NewReader(pkt.Data)
		id, err := packet.ReadInt32(r)
		if err != nil {
			return fmt.Errorf("读取 play ping 失败: %w", err)
		}
		return c.conn.WritePacket(protocol.PlayServerPong, packet.EncodeInt32(id))

	case protocol.PlayClientDeclareCommands:
		return nil

	case protocol.PlayClientProfilelessChat:
		return nil

	case protocol.PlayClientPlayerChat:
		return nil

	case protocol.PlayClientPosition:
		r := bytes.NewReader(pkt.Data)
		teleportID, err := packet.ReadVarInt(r)
		if err != nil {
			return fmt.Errorf("读取 player_position teleport id 失败: %w", err)
		}
		x := packet.MustReadFloat64(r, "player_position.x")
		y := packet.MustReadFloat64(r, "player_position.y")
		z := packet.MustReadFloat64(r, "player_position.z")
		_ = packet.MustReadFloat64(r, "player_position.deltaX")
		_ = packet.MustReadFloat64(r, "player_position.deltaY")
		_ = packet.MustReadFloat64(r, "player_position.deltaZ")
		yRotBytes := packet.MustReadBytes(r, 4, "player_position.yRot")
		yRot := packet.ReadFloat32FromBytes(yRotBytes)
		xRotBytes := packet.MustReadBytes(r, 4, "player_position.xRot")
		xRot := packet.ReadFloat32FromBytes(xRotBytes)
		relBits := packet.MustReadVarInt(r, "player_position.relBits")

		px, py, pz := c.Player.GetPosition()
		if relBits&0x01 != 0 {
			x = px + x
		}
		if relBits&0x02 != 0 {
			y = py + y
		}
		if relBits&0x04 != 0 {
			z = pz + z
		}

		c.Player.UpdatePosition(x, y, z, yRot, xRot, 0)

		if err := c.conn.WritePacket(protocol.PlayServerAcceptTeleport, packet.EncodeVarInt(teleportID)); err != nil {
			return fmt.Errorf("发送 accept_teleportation 失败: %w", err)
		}
		return nil

	case protocol.PlayClientActionBar:
		return nil

	case protocol.PlayClientCookieReq:
		r := bytes.NewReader(pkt.Data)
		key, err := packet.ReadString(r, r)
		if err != nil {
			return fmt.Errorf("读取 play cookie key 失败: %w", err)
		}
		return c.sendCookieResponse(protocol.PlayServerCookieResp, key)

	case protocol.PlayClientPackPush:
		id, err := readPacketUUID(pkt.Data)
		if err != nil {
			return fmt.Errorf("读取 play resource_pack_push 失败: %w", err)
		}
		return c.sendResourcePackResponses(protocol.PlayServerResource, id)

	case protocol.PlayClientLogin:
		r := bytes.NewReader(pkt.Data)
		entityID, err := packet.ReadInt32(r)
		if err != nil {
			return fmt.Errorf("读取 entity_id 失败: %w", err)
		}
		c.Player.SetEntityID(entityID)
		if !c.inPlay {
			c.inPlay = true
			logx.Infof("已进入服务器, 开始挂机: %s (%s)", c.username, packet.FormatUUID(c.uuid))
		}
		return nil

	case protocol.PlayClientSystemChat:
		return nil

	case protocol.PlayClientPackPop:
		return nil

	case protocol.PlayClientSetHealth:
		r := bytes.NewReader(pkt.Data)
		health, _ := packet.ReadFloat32(r)
		food := packet.MustReadVarInt(r, "food")
		saturation, _ := packet.ReadFloat32(r)
		c.Player.UpdateHealth(health, c.Player.MaxHealth, food, saturation)

		// 检测死亡并自动重生
		if health <= 0 {
			logx.Infof("检测到玩家死亡，自动重生")
			go func() {
				time.Sleep(500 * time.Millisecond)
				if err := c.SendClientCommand(0); err != nil {
					logx.Warnf("发送重生数据包失败: %v", err)
				}
			}()
		}
		return nil

	case protocol.PlayClientSetExperience:
		r := bytes.NewReader(pkt.Data)
		experience, _ := packet.ReadFloat32(r)
		level := packet.MustReadVarInt(r, "level")
		c.Player.UpdateExperience(level, experience)
		return nil

	case protocol.PlayClientSetHeldSlot:
		return nil

	case protocol.PlayClientContainerContent:
		return nil

	case protocol.PlayClientContainerSlot:
		return nil

	case protocol.PlayClientContainerClose:
		return nil

	case protocol.PlayClientContainerSetData:
		return nil

	case protocol.PlayClientOpenScreen:
		return nil

	case protocol.PlayClientPlayerAbilities:
		return nil

	case protocol.PlayClientGameEvent:
		return nil

	case protocol.PlayClientEntityData:
		return nil

	case protocol.PlayClientPlayerInfoUpdate:
		return nil

	case protocol.PlayClientPlayerInfoRemove:
		return nil

	case protocol.PlayClientAddEntity:
		return nil

	case protocol.PlayClientTeleportEntity:
		return nil

	case protocol.PlayClientMoveEntityPos:
		return nil

	case protocol.PlayClientRemoveEntities:
		return nil

	default:
		logx.PacketLogf("未处理的 Play 数据包: id=0x%02X (%s) len=%d", pkt.ID, protocol.PacketName(protocol.StatePlay, pkt.ID), len(pkt.Data))
		return nil
	}
}

// 加密处理
func (c *Client) handleEncryptionRequest(data []byte) error {
	r := bytes.NewReader(data)
	serverID, err := packet.ReadString(r, r)
	if err != nil {
		return fmt.Errorf("读取 encrypt serverId 失败: %w", err)
	}
	publicKeyDER, err := packet.ReadByteArray(r, r)
	if err != nil {
		return fmt.Errorf("读取 encrypt publicKey 失败: %w", err)
	}
	challenge, err := packet.ReadByteArray(r, r)
	if err != nil {
		return fmt.Errorf("读取 encrypt challenge 失败: %w", err)
	}
	shouldAuthenticate, err := packet.ReadBool(r)
	if err != nil {
		return fmt.Errorf("读取 shouldAuthenticate 失败: %w", err)
	}
	logx.Debugf("收到 Encryption Request: serverId=%q publicKeyLen=%d verifyTokenLen=%d shouldAuthenticate=%t",
		serverID, len(publicKeyDER), len(challenge), shouldAuthenticate)

	if shouldAuthenticate && c.online == nil {
		return fmt.Errorf("未获取正版认证会话，无法完成加密握手")
	}

	pubAny, err := x509.ParsePKIXPublicKey(publicKeyDER)
	if err != nil {
		return fmt.Errorf("解析 RSA 公钥失败: %w", err)
	}
	pubKey, ok := pubAny.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("服务器公钥不是 RSA")
	}

	sharedSecret := make([]byte, 16)
	if _, err := rand.Read(sharedSecret); err != nil {
		return fmt.Errorf("生成 shared secret 失败: %w", err)
	}

	if shouldAuthenticate {
		serverHash := packet.MinecraftServerHash(serverID, sharedSecret, publicKeyDER)
		logx.Debugf("准备调用 session join: profile=%s serverHash=%s", c.online.ProfileID, serverHash)
		if err := mcauth.JoinServer(c.online.AccessToken, c.online.ProfileID, serverHash); err != nil {
			return fmt.Errorf("正版会话认证失败: %w", err)
		}
	}

	encryptedSecret, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, sharedSecret)
	if err != nil {
		return fmt.Errorf("加密 shared secret 失败: %w", err)
	}
	encryptedChallenge, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, challenge)
	if err != nil {
		return fmt.Errorf("加密 challenge 失败: %w", err)
	}

	payload := make([]byte, 0, len(encryptedSecret)+len(encryptedChallenge)+10)
	payload = append(payload, packet.EncodeByteArray(encryptedSecret)...)
	if protocol.Features774.SignatureEncryption {
		return fmt.Errorf("当前协议实现未启用 signatureEncryption 分支")
	}
	payload = append(payload, packet.EncodeByteArray(encryptedChallenge)...)
	logx.Debugf("准备发送 Encryption Response(legacy): encryptedSecretLen=%d encryptedVerifyTokenLen=%d",
		len(encryptedSecret), len(encryptedChallenge))

	if err := c.conn.WritePacket(protocol.LoginServerKey, payload); err != nil {
		return fmt.Errorf("发送 login key response 失败: %w", err)
	}
	if err := c.conn.EnableEncryption(sharedSecret); err != nil {
		return fmt.Errorf("启用 AES/CFB8 加密失败: %w", err)
	}
	return nil
}

func parseGameProfile(data []byte) ([16]byte, string, error) {
	r := bytes.NewReader(data)
	uuid, err := packet.ReadUUID(r)
	if err != nil {
		return [16]byte{}, "", err
	}
	name, err := packet.ReadString(r, r)
	if err != nil {
		return [16]byte{}, "", err
	}
	propCount, err := packet.ReadVarInt(r)
	if err != nil {
		return [16]byte{}, "", err
	}
	for i := 0; i < int(propCount); i++ {
		if _, err := packet.ReadString(r, r); err != nil {
			return [16]byte{}, "", err
		}
		if _, err := packet.ReadString(r, r); err != nil {
			return [16]byte{}, "", err
		}
		hasSig, err := packet.ReadBool(r)
		if err != nil {
			return [16]byte{}, "", err
		}
		if hasSig {
			if _, err := packet.ReadString(r, r); err != nil {
				return [16]byte{}, "", err
			}
		}
	}
	return uuid, name, nil
}

// 发送客户端信息
func (c *Client) sendClientInformation(packetID int32) error {
	payload := make([]byte, 0, 64)
	payload = append(payload, packet.EncodeString("zh_cn")...)
	payload = append(payload, byte(8))
	payload = append(payload, packet.EncodeVarInt(0)...)
	payload = append(payload, packet.EncodeBool(true)...)
	payload = append(payload, byte(0x7F))
	payload = append(payload, packet.EncodeVarInt(1)...)
	payload = append(payload, packet.EncodeBool(false)...)
	payload = append(payload, packet.EncodeBool(true)...)
	payload = append(payload, packet.EncodeVarInt(0)...)
	return c.conn.WritePacket(packetID, payload)
}

// 辅助函数
func (c *Client) sendCookieResponse(packetID int32, key string) error {
	payload := make([]byte, 0, len(key)+16)
	payload = append(payload, packet.EncodeString(key)...)
	payload = append(payload, packet.EncodeBool(false)...)
	return c.conn.WritePacket(packetID, payload)
}

func (c *Client) sendResourcePackResponses(packetID int32, id [16]byte) error {
	for _, action := range []int32{protocol.ResourcePackAccepted, protocol.ResourcePackDownloaded, protocol.ResourcePackLoaded} {
		payload := make([]byte, 0, 24)
		payload = append(payload, id[:]...)
		payload = append(payload, packet.EncodeVarInt(action)...)
		if err := c.conn.WritePacket(packetID, payload); err != nil {
			return err
		}
	}
	return nil
}

func readPacketUUID(data []byte) ([16]byte, error) {
	r := bytes.NewReader(data)
	return packet.ReadUUID(r)
}

func (c *Client) sendHandshake(host string, port uint16) error {
	payload := make([]byte, 0, 128)
	payload = append(payload, packet.EncodeVarInt(protocol.Version)...)
	payload = append(payload, packet.EncodeString(host)...)
	var p [2]byte
	binary.BigEndian.PutUint16(p[:], port)
	payload = append(payload, p[:]...)
	payload = append(payload, packet.EncodeVarInt(protocol.IntentionLogin)...)
	return c.conn.WritePacket(protocol.HandshakingServerIntention, payload)
}

func (c *Client) sendLoginStart() error {
	payload := make([]byte, 0, 96)
	payload = append(payload, packet.EncodeString(c.username)...)
	payload = append(payload, c.uuid[:]...)
	logx.Debugf("准备发送 Login Start: username=%s uuid=%s", c.username, packet.FormatUUID(c.uuid))
	return c.conn.WritePacket(protocol.LoginServerHello, payload)
}

func disconnectReasonFromNBT(data []byte) string {
	rawJSON, err := packet.ReadAnonymousNBTJSON(bytes.NewReader(data))
	if err != nil {
		return packet.RawPreview(data)
	}
	return rawJSON
}

func disconnectReasonFromJSON(data []byte) string {
	r := bytes.NewReader(data)
	rawJSON, err := packet.ReadString(r, r)
	if err != nil {
		return packet.RawPreview(data)
	}
	return rawJSON
}
