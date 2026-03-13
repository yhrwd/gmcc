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
	"gmcc/internal/mcclient/packet"
	"gmcc/internal/mcclient/protocol"
)

func (c *Client) handleLoginPacket(pkt packet.Packet) error {
	switch pkt.ID {
	case protocol.LoginClientDisconnect:
		r := bytes.NewReader(pkt.Data)
		reason, err := packet.ReadString(r, r)
		if err != nil {
			reason = packet.RawPreview(pkt.Data)
		}
		return fmt.Errorf("登录阶段被服务器断开: %s", reason)

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

func (c *Client) handleConfigurationPacket(pkt packet.Packet) error {
	switch pkt.ID {
	case protocol.CfgClientDisconnect:
		return fmt.Errorf("配置阶段被服务器断开: %s", packet.RawPreview(pkt.Data))

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
		serverHash := packet.MinecraftServerHash(serverID, sharedSecret, publicKeyDER)
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
	payload = append(payload, packet.EncodeByteArray(encryptedSecret)...)
	if protocol.Features774.SignatureEncryption {
		return fmt.Errorf("当前协议实现未启用 signatureEncryption 分支")
	}
	payload = append(payload, packet.EncodeByteArray(encryptedChallenge)...)
	logx.Debugf("准备发送 Encryption Response(legacy): encryptedSecretLen=%d encryptedVerifyTokenLen=%d", len(encryptedSecret), len(encryptedChallenge))

	if err := c.conn.WritePacket(protocol.LoginServerKey, payload); err != nil {
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
	payload = append(payload, packet.EncodeVarInt(protocol.Version)...)
	payload = append(payload, packet.EncodeString(host)...)

	var p [2]byte
	binary.BigEndian.PutUint16(p[:], port)
	payload = append(payload, p[:]...)
	payload = append(payload, packet.EncodeVarInt(protocol.IntentionLogin)...)
	logx.PacketLogf("准备发送 Handshake: host=%s port=%d protocol=%d nextState=login", host, port, protocol.Version)

	return c.conn.WritePacket(protocol.HandshakingServerIntention, payload)
}

func (c *Client) sendLoginStart() error {
	payload := make([]byte, 0, 96)
	payload = append(payload, packet.EncodeString(c.username)...)
	payload = append(payload, c.uuid[:]...)
	logx.PacketLogf("准备发送 Login Start: username=%s uuid=%s", c.username, packet.FormatUUID(c.uuid))
	return c.conn.WritePacket(protocol.LoginServerHello, payload)
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
