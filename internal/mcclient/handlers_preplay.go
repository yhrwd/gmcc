package mcclient

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"fmt"

	mcauth "gmcc/internal/auth/minecraft"
	"gmcc/internal/logx"
)

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
