package mcclient

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"gmcc/internal/auth/microsoft"
	mcauth "gmcc/internal/auth/minecraft"
	"gmcc/internal/config"
	"gmcc/internal/logx"
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

	state  connState
	inPlay bool

	conn *packetConn

	lastAFKPacket time.Time
	chatHandler   func(ChatMessage)
	chatSessionOK bool
	chatSession   *secureChatSession
	commandSign   map[string]signableCommandTarget
	chatSignMu    sync.Mutex
}

func New(cfg *config.Config) *Client {
	name := strings.TrimSpace(cfg.Account.PlayerID)
	client := &Client{
		cfg:         cfg,
		offlineName: name,
		offlineUUID: offlineUUID(name),
		username:    name,
		uuid:        offlineUUID(name),
		commandSign: map[string]signableCommandTarget{},
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
	logx.Infof("客户端协议: %s", protocolLabel)
	host, port, err := parseAddress(c.cfg.Server.Address)
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

	c.state = stateLogin
	c.inPlay = false
	c.chatSessionOK = false
	c.chatSession = nil
	c.commandSign = map[string]signableCommandTarget{}
	c.lastAFKPacket = time.Now()

	addr := net.JoinHostPort(host, strconv.Itoa(int(port)))
	dialer := net.Dialer{Timeout: 10 * time.Second}
	rawConn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("连接服务器失败: %w", err)
	}

	c.conn = newPacketConn(rawConn)
	defer func() {
		_ = c.conn.Close()
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

		_ = c.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		pkt, err := c.conn.ReadPacket()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				if c.state == statePlay {
					if err := c.sendAFKHeartbeatIfNeeded(); err != nil {
						return err
					}
				}
				continue
			}
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				return fmt.Errorf("服务器已关闭连接 (state=%s): %w", stateName(c.state), err)
			}
			return fmt.Errorf("读取数据包失败 (state=%s): %w", stateName(c.state), err)
		}

		logx.PacketLogf(
			"收到数据包: state=%s id=0x%02X (%s) len=%d",
			stateName(c.state),
			pkt.ID,
			packetName(c.state, pkt.ID),
			len(pkt.Data),
		)
		if err := c.handlePacket(pkt); err != nil {
			return err
		}

		if c.state == statePlay {
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
			logx.Infof("使用缓存的 Minecraft token: %s (%s)", c.online.ProfileName, formatUUID(c.online.ProfileUUID))
			return nil
		}
		logx.Warnf("检测到缓存的 Minecraft token 数据损坏，准备重新认证")
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

	logx.Infof("正版认证成功: %s (%s)", c.online.ProfileName, formatUUID(c.online.ProfileUUID))
	return nil
}

func (c *Client) resolveXSTSToken(cache *session.TokenCache, now time.Time) (*microsoft.XSTSResponse, *microsoft.TokenResponse, error) {
	if cache != nil && cache.HasValidMicrosoftAccess(now) {
		xstsResp, err := microsoft.GetXSTSTokenFromAccessToken(cache.Microsoft.AccessToken)
		if err == nil {
			logx.Infof("使用缓存的 Microsoft access_token")
			return xstsResp, nil, nil
		}
		logx.Warnf("缓存 Microsoft access_token 已失效，将尝试 refresh_token: %v", err)
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

	profileUUID, err := parseUUID(cleanID)
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

func (c *Client) handlePacket(pkt packet) error {
	switch c.state {
	case stateLogin:
		return c.handleLoginPacket(pkt)
	case stateConfiguration:
		return c.handleConfigurationPacket(pkt)
	case statePlay:
		return c.handlePlayPacket(pkt)
	default:
		return fmt.Errorf("未知连接状态: %s", stateName(c.state))
	}
}

func (c *Client) handleLoginPacket(pkt packet) error {
	switch pkt.ID {
	case loginClientDisconnect:
		r := bytes.NewReader(pkt.Data)
		reason, err := readString(r, r)
		if err != nil {
			reason = rawPreview(pkt.Data)
		}
		return fmt.Errorf("登录阶段被服务器断开: %s", reason)

	case loginClientCompression:
		r := bytes.NewReader(pkt.Data)
		threshold, err := readVarInt(r)
		if err != nil {
			return fmt.Errorf("读取压缩阈值失败: %w", err)
		}
		c.conn.SetCompressionThreshold(int(threshold))
		logx.Infof("启用压缩, threshold=%d", threshold)
		return nil

	case loginClientHello:
		return c.handleEncryptionRequest(pkt.Data)

	case loginClientFinished:
		profileUUID, profileName, err := parseGameProfile(pkt.Data)
		if err != nil {
			return fmt.Errorf("解析 login_finished 失败: %w", err)
		}
		if profileName != "" {
			logx.Infof("登录成功: %s (%s)", profileName, formatUUID(profileUUID))
		} else {
			logx.Infof("登录成功: %s", formatUUID(profileUUID))
		}

		if err := c.conn.WritePacket(loginServerAck, nil); err != nil {
			return fmt.Errorf("发送 login_acknowledged 失败: %w", err)
		}
		c.state = stateConfiguration
		if err := c.sendClientInformation(cfgServerClientInfo); err != nil {
			return fmt.Errorf("发送 configuration client_information 失败: %w", err)
		}
		return nil

	case loginClientCustomQuery:
		r := bytes.NewReader(pkt.Data)
		txID, err := readVarInt(r)
		if err != nil {
			return fmt.Errorf("读取 custom_query transactionId 失败: %w", err)
		}
		payload := make([]byte, 0, 8)
		payload = append(payload, encodeVarInt(txID)...)
		payload = append(payload, encodeBool(false)...)
		if err := c.conn.WritePacket(loginServerCustomAnswer, payload); err != nil {
			return fmt.Errorf("发送 custom_query_answer 失败: %w", err)
		}
		return nil

	case loginClientCookieReq:
		r := bytes.NewReader(pkt.Data)
		key, err := readString(r, r)
		if err != nil {
			return fmt.Errorf("读取 login cookie key 失败: %w", err)
		}
		return c.sendCookieResponse(loginServerCookieResp, key)

	default:
		logx.PacketLogf("未处理的 Login 数据包: id=0x%02X (%s) len=%d", pkt.ID, packetName(stateLogin, pkt.ID), len(pkt.Data))
		return nil
	}
}

func (c *Client) handleConfigurationPacket(pkt packet) error {
	switch pkt.ID {
	case cfgClientDisconnect:
		return fmt.Errorf("配置阶段被服务器断开: %s", rawPreview(pkt.Data))

	case cfgClientKeepAlive:
		r := bytes.NewReader(pkt.Data)
		id, err := readInt64(r)
		if err != nil {
			return fmt.Errorf("读取配置阶段 keep_alive 失败: %w", err)
		}
		return c.conn.WritePacket(cfgServerKeepAlive, encodeInt64(id))

	case cfgClientPing:
		r := bytes.NewReader(pkt.Data)
		id, err := readInt32(r)
		if err != nil {
			return fmt.Errorf("读取配置阶段 ping 失败: %w", err)
		}
		return c.conn.WritePacket(cfgServerPong, encodeInt32(id))

	case cfgClientCookieReq:
		r := bytes.NewReader(pkt.Data)
		key, err := readString(r, r)
		if err != nil {
			return fmt.Errorf("读取配置阶段 cookie key 失败: %w", err)
		}
		return c.sendCookieResponse(cfgServerCookieResp, key)

	case cfgClientSelectPacks:
		payload := encodeVarInt(0)
		return c.conn.WritePacket(cfgServerSelectPacks, payload)

	case cfgClientPackPush:
		id, err := readPacketUUID(pkt.Data)
		if err != nil {
			return fmt.Errorf("读取配置阶段 resource_pack_push 失败: %w", err)
		}
		return c.sendResourcePackResponses(cfgServerResource, id)

	case cfgClientCodeOfConduct:
		return c.conn.WritePacket(cfgServerAcceptCode, nil)

	case cfgClientFinish:
		if err := c.conn.WritePacket(cfgServerFinish, nil); err != nil {
			return fmt.Errorf("发送 finish_configuration 失败: %w", err)
		}
		c.state = statePlay
		if err := c.sendClientInformation(playServerClientInfo); err != nil {
			return fmt.Errorf("发送 play client_information 失败: %w", err)
		}
		logx.Infof("进入 Play 阶段")
		return nil

	case cfgClientCustom, cfgClientRegistry, cfgClientPackPop:
		return nil

	default:
		logx.PacketLogf("未处理的 Configuration 数据包: id=0x%02X (%s) len=%d", pkt.ID, packetName(stateConfiguration, pkt.ID), len(pkt.Data))
		return nil
	}
}

func (c *Client) handlePlayPacket(pkt packet) error {
	switch pkt.ID {
	case playClientDisconnect:
		return fmt.Errorf("Play 阶段被服务器断开: %s", rawPreview(pkt.Data))

	case playClientKeepAlive:
		r := bytes.NewReader(pkt.Data)
		id, err := readInt64(r)
		if err != nil {
			return fmt.Errorf("读取 play keep_alive 失败: %w", err)
		}
		return c.conn.WritePacket(playServerKeepAlive, encodeInt64(id))

	case playClientPing:
		r := bytes.NewReader(pkt.Data)
		id, err := readInt32(r)
		if err != nil {
			return fmt.Errorf("读取 play ping 失败: %w", err)
		}
		return c.conn.WritePacket(playServerPong, encodeInt32(id))

	case playClientDeclareCommands:
		return c.handleDeclareCommandsPacket(pkt.Data)

	case playClientProfilelessChat:
		return c.handleProfilelessChatPacket(pkt.Data)

	case playClientPlayerChat:
		return c.handlePlayerChatPacket(pkt.Data)

	case playClientPosition:
		r := bytes.NewReader(pkt.Data)
		teleportID, err := readVarInt(r)
		if err != nil {
			return fmt.Errorf("读取 player_position teleport id 失败: %w", err)
		}
		if err := c.conn.WritePacket(playServerAcceptTeleport, encodeVarInt(teleportID)); err != nil {
			return fmt.Errorf("发送 accept_teleportation 失败: %w", err)
		}
		return nil

	case playClientActionBar:
		return c.handleActionBarPacket(pkt.Data)

	case playClientCookieReq:
		r := bytes.NewReader(pkt.Data)
		key, err := readString(r, r)
		if err != nil {
			return fmt.Errorf("读取 play cookie key 失败: %w", err)
		}
		return c.sendCookieResponse(playServerCookieResp, key)

	case playClientPackPush:
		id, err := readPacketUUID(pkt.Data)
		if err != nil {
			return fmt.Errorf("读取 play resource_pack_push 失败: %w", err)
		}
		return c.sendResourcePackResponses(playServerResource, id)

	case playClientLogin:
		if !c.inPlay {
			c.inPlay = true
			logx.Infof("已进入服务器, 开始挂机: %s (%s)", c.username, formatUUID(c.uuid))
			if err := c.initSecureChatSession(); err != nil {
				logx.Warnf("初始化 secure chat 会话失败: %v", err)
			}
			c.runOnJoinActions()
		}
		return nil

	case playClientSystemChat:
		return c.handleSystemChatPacket(pkt.Data)

	case playClientPackPop:
		return nil

	default:
		logx.PacketLogf("未处理的 Play 数据包: id=0x%02X (%s) len=%d", pkt.ID, packetName(statePlay, pkt.ID), len(pkt.Data))
		return nil
	}
}

func (c *Client) sendHandshake(host string, port uint16) error {
	payload := make([]byte, 0, 128)
	payload = append(payload, encodeVarInt(protocolVersion)...)
	payload = append(payload, encodeString(host)...)

	var p [2]byte
	binary.BigEndian.PutUint16(p[:], port)
	payload = append(payload, p[:]...)
	payload = append(payload, encodeVarInt(intentionLogin)...)
	logx.PacketLogf("准备发送 Handshake: host=%s port=%d protocol=%d nextState=login", host, port, protocolVersion)

	return c.conn.WritePacket(handshakingServerIntention, payload)
}

func (c *Client) sendLoginStart() error {
	payload := make([]byte, 0, 96)
	payload = append(payload, encodeString(c.username)...)
	payload = append(payload, c.uuid[:]...)
	logx.PacketLogf("准备发送 Login Start: username=%s uuid=%s", c.username, formatUUID(c.uuid))
	return c.conn.WritePacket(loginServerHello, payload)
}

func (c *Client) sendCookieResponse(packetID int32, key string) error {
	payload := make([]byte, 0, len(key)+16)
	payload = append(payload, encodeString(key)...)
	payload = append(payload, encodeBool(false)...)
	return c.conn.WritePacket(packetID, payload)
}

func (c *Client) sendClientInformation(packetID int32) error {
	payload := make([]byte, 0, 64)
	payload = append(payload, encodeString("zh_cn")...)
	payload = append(payload, byte(8))
	payload = append(payload, encodeVarInt(0)...)
	payload = append(payload, encodeBool(true)...)
	payload = append(payload, byte(0x7F))
	payload = append(payload, encodeVarInt(1)...)
	payload = append(payload, encodeBool(false)...)
	payload = append(payload, encodeBool(true)...)
	payload = append(payload, encodeVarInt(0)...)
	return c.conn.WritePacket(packetID, payload)
}

func (c *Client) sendResourcePackResponses(packetID int32, id [16]byte) error {
	for _, action := range []int32{resourcePackAccepted, resourcePackDownloaded, resourcePackLoaded} {
		payload := make([]byte, 0, 24)
		payload = append(payload, id[:]...)
		payload = append(payload, encodeVarInt(action)...)
		if err := c.conn.WritePacket(packetID, payload); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) sendAFKHeartbeatIfNeeded() error {
	if time.Since(c.lastAFKPacket) < 15*time.Second {
		return nil
	}
	c.lastAFKPacket = time.Now()
	return c.conn.WritePacket(playServerMoveStatus, []byte{0x01})
}

func (c *Client) handleEncryptionRequest(data []byte) error {
	r := bytes.NewReader(data)
	serverID, err := readString(r, r)
	if err != nil {
		return fmt.Errorf("读取 encrypt serverId 失败: %w", err)
	}
	publicKeyDER, err := readByteArray(r, r)
	if err != nil {
		return fmt.Errorf("读取 encrypt publicKey 失败: %w", err)
	}
	challenge, err := readByteArray(r, r)
	if err != nil {
		return fmt.Errorf("读取 encrypt challenge 失败: %w", err)
	}
	shouldAuthenticate, err := readBool(r)
	if err != nil {
		return fmt.Errorf("读取 shouldAuthenticate 失败: %w", err)
	}
	logx.Debugf(
		"收到 Encryption Request: serverId=%q publicKeyLen=%d verifyTokenLen=%d shouldAuthenticate=%t",
		serverID,
		len(publicKeyDER),
		len(challenge),
		shouldAuthenticate,
	)

	if shouldAuthenticate && c.online == nil {
		return errOnlineAuthRequired
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
		serverHash := minecraftServerHash(serverID, sharedSecret, publicKeyDER)
		logx.Debugf("准备调用 session join: profile=%s serverHash=%s", c.online.ProfileID, serverHash)
		if err := mcauth.JoinServer(c.online.AccessToken, c.online.ProfileID, serverHash); err != nil {
			return fmt.Errorf("正版会话认证失败: %w", err)
		}
		logx.Infof("正版会话认证成功")
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
	payload = append(payload, encodeByteArray(encryptedSecret)...)
	if features774.SignatureEncryption {
		return fmt.Errorf("当前协议实现未启用 signatureEncryption 分支")
	}
	payload = append(payload, encodeByteArray(encryptedChallenge)...)
	logx.Debugf("准备发送 Encryption Response(legacy): encryptedSecretLen=%d encryptedVerifyTokenLen=%d", len(encryptedSecret), len(encryptedChallenge))

	if err := c.conn.WritePacket(loginServerKey, payload); err != nil {
		return fmt.Errorf("发送 login key response 失败: %w", err)
	}
	if err := c.conn.EnableEncryption(sharedSecret); err != nil {
		return fmt.Errorf("启用 AES/CFB8 加密失败: %w", err)
	}
	if shouldAuthenticate {
		logx.Infof("已完成正版加密握手, serverId=%q", serverID)
	} else {
		logx.Infof("已完成非会话加密握手, serverId=%q", serverID)
	}
	return nil
}

func parseGameProfile(data []byte) ([16]byte, string, error) {
	r := bytes.NewReader(data)

	uuid, err := readUUID(r)
	if err != nil {
		return [16]byte{}, "", err
	}
	name, err := readString(r, r)
	if err != nil {
		return [16]byte{}, "", err
	}
	propCount, err := readVarInt(r)
	if err != nil {
		return [16]byte{}, "", err
	}
	for i := 0; i < int(propCount); i++ {
		if _, err := readString(r, r); err != nil {
			return [16]byte{}, "", err
		}
		if _, err := readString(r, r); err != nil {
			return [16]byte{}, "", err
		}
		hasSig, err := readBool(r)
		if err != nil {
			return [16]byte{}, "", err
		}
		if hasSig {
			if _, err := readString(r, r); err != nil {
				return [16]byte{}, "", err
			}
		}
	}

	return uuid, name, nil
}

func readPacketUUID(data []byte) ([16]byte, error) {
	r := bytes.NewReader(data)
	return readUUID(r)
}

func parseAddress(addr string) (string, uint16, error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return "", 0, fmt.Errorf("server.address 不能为空")
	}

	if !strings.Contains(addr, ":") {
		return addr, 25565, nil
	}

	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		idx := strings.LastIndex(addr, ":")
		if idx <= 0 {
			return "", 0, fmt.Errorf("无效地址: %s", addr)
		}
		host = addr[:idx]
		portStr = addr[idx+1:]
	}

	portInt, err := strconv.Atoi(portStr)
	if err != nil || portInt <= 0 || portInt > 65535 {
		return "", 0, fmt.Errorf("无效端口: %s", portStr)
	}
	return host, uint16(portInt), nil
}

func offlineUUID(name string) [16]byte {
	hash := md5.Sum([]byte("OfflinePlayer:" + name))
	hash[6] = (hash[6] & 0x0F) | 0x30
	hash[8] = (hash[8] & 0x3F) | 0x80
	return hash
}

func parseUUID(raw string) ([16]byte, error) {
	clean := strings.ReplaceAll(strings.TrimSpace(raw), "-", "")
	if len(clean) != 32 {
		return [16]byte{}, fmt.Errorf("uuid 长度无效")
	}
	b, err := hex.DecodeString(clean)
	if err != nil {
		return [16]byte{}, err
	}
	var id [16]byte
	copy(id[:], b)
	return id, nil
}

func formatUUID(id [16]byte) string {
	hexStr := hex.EncodeToString(id[:])
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hexStr[0:8],
		hexStr[8:12],
		hexStr[12:16],
		hexStr[16:20],
		hexStr[20:32],
	)
}

func minecraftServerHash(serverID string, sharedSecret []byte, publicKey []byte) string {
	h := sha1.New()
	_, _ = h.Write([]byte(serverID))
	_, _ = h.Write(sharedSecret)
	_, _ = h.Write(publicKey)
	sum := h.Sum(nil)

	if len(sum) == 0 {
		return "0"
	}
	if sum[0]&0x80 != 0 {
		neg := make([]byte, len(sum))
		for i := range sum {
			neg[i] = ^sum[i]
		}
		for i := len(neg) - 1; i >= 0; i-- {
			neg[i]++
			if neg[i] != 0 {
				break
			}
		}
		n := new(big.Int).SetBytes(neg)
		if n.Sign() == 0 {
			return "0"
		}
		return "-" + n.Text(16)
	}

	n := new(big.Int).SetBytes(sum)
	if n.Sign() == 0 {
		return "0"
	}
	return n.Text(16)
}

func rawPreview(b []byte) string {
	if len(b) == 0 {
		return "<empty>"
	}
	const max = 120
	if len(b) <= max {
		return hex.EncodeToString(b)
	}
	return hex.EncodeToString(b[:max]) + "..."
}
