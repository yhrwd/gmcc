package mcclient

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"strings"
	"time"

	mcauth "gmcc/internal/auth/minecraft"
	"gmcc/internal/logx"
	"gmcc/internal/nbt"
)

// ChatMessage 是统一的聊天事件结构，方便后续接入机器人/插件解析。
type ChatMessage struct {
	Type        string
	PlainText   string
	RawJSON     string
	IsActionBar bool
	SenderUUID  string
	ReceivedAt  time.Time
}

type secureChatSession struct {
	sessionID    [16]byte
	privateKey   *rsa.PrivateKey
	messageIndex int32
}

type commandArgumentSignature struct {
	name      string
	signature []byte
}

type signableCommandTarget struct {
	ArgumentName string
	SliceIndex   int
}

func (c *Client) SetChatHandler(handler func(ChatMessage)) {
	c.chatHandler = handler
}

func (c *Client) SendCommand(command string) error {
	cmd := strings.TrimSpace(command)
	cmd = strings.TrimPrefix(cmd, "/")
	if cmd == "" {
		return fmt.Errorf("command 不能为空")
	}
	if c.state != statePlay {
		return fmt.Errorf("当前状态不是 Play，无法发送命令")
	}

	logx.Debugf("发送命令原始内容: %q", cmd)

	if err := c.sendUnsignedCommand(cmd); err != nil {
		return fmt.Errorf("发送命令失败: %w", err)
	}
	logx.Debugf("已发送命令: /%s", cmd)
	return nil
}

func (c *Client) SendCommandSigned(command string) error {
	cmd := strings.TrimSpace(command)
	cmd = strings.TrimPrefix(cmd, "/")
	if cmd == "" {
		return fmt.Errorf("command 不能为空")
	}
	if c.state != statePlay {
		return fmt.Errorf("当前状态不是 Play，无法发送命令")
	}

	if err := c.sendSignedCommand(cmd); err != nil {
		logx.Debugf("签名命令失败，尝试无签名: %v", err)
		if err2 := c.sendUnsignedCommand(cmd); err2 != nil {
			return fmt.Errorf("发送命令失败: %w", err2)
		}
	}
	logx.Debugf("已发送命令: /%s", cmd)
	return nil
}

// SendMessage 发送未签名聊天消息。
// 注：部分开启 secure chat 的服务器会拒绝未签名消息。
func (c *Client) SendMessage(message string) error {
	msg := strings.TrimSpace(message)
	if msg == "" {
		return fmt.Errorf("message 不能为空")
	}
	if c.state != statePlay {
		return fmt.Errorf("当前状态不是 Play，无法发送聊天消息")
	}

	salt, err := randomInt64()
	if err != nil {
		return err
	}
	timestamp := time.Now().UnixMilli()

	signature, signErr := c.signChatBody(msg, timestamp, salt, nil)
	hasSignature := signErr == nil && len(signature) == 256
	if signErr != nil {
		logx.Debugf("聊天消息签名不可用，按无签名发送: %v", signErr)
	}

	payload := make([]byte, 0, len(msg)+64)
	payload = append(payload, encodeString(msg)...)
	payload = append(payload, encodeInt64(timestamp)...)
	payload = append(payload, encodeInt64(salt)...)
	payload = append(payload, encodeBool(hasSignature)...)
	if hasSignature {
		payload = append(payload, signature...)
	}
	payload = append(payload, encodeVarInt(0)...)
	payload = append(payload, 0x00, 0x00, 0x00)
	payload = append(payload, 0x01)

	if err := c.conn.WritePacket(playServerChatMessage, payload); err != nil {
		return fmt.Errorf("发送聊天消息失败: %w", err)
	}
	logx.Infof("已发送聊天消息: %s", msg)
	return nil
}

func (c *Client) runOnJoinActions() {
	go c.runOnJoinActionsAsync()
}

func (c *Client) runOnJoinActionsAsync() {
	delay := c.cfg.Actions.DelayMs
	if delay <= 0 {
		delay = 1200
	}
	time.Sleep(time.Duration(delay) * time.Millisecond)

	signCmds := c.cfg.Actions.SignCommands

	for _, cmd := range c.cfg.Actions.OnJoinCommands {
		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}
		cmd = strings.TrimPrefix(cmd, "/")

		var err error
		if signCmds {
			err = c.SendCommandSigned(cmd)
		} else {
			err = c.SendCommand(cmd)
		}
		if err != nil {
			logx.Warnf("执行命令失败 /%s: %v", cmd, err)
		}
	}
	for _, msg := range c.cfg.Actions.OnJoinMessages {
		if strings.TrimSpace(msg) == "" {
			continue
		}
		if err := c.SendMessage(msg); err != nil {
			logx.Warnf("发送消息失败: %v", err)
		}
	}

	c.syncPlayerInfo()
}

func (c *Client) syncPlayerInfo() {
	time.Sleep(5 * time.Second)

	info := c.Player.GetInfo()
	logx.Infof("=== 玩家信息同步 ===")
	logx.Infof("名称: %v", info["name"])
	logx.Infof("UUID: %v", info["uuid"])
	logx.Infof("实体ID: %v", info["entity_id"])
	logx.Infof("游戏模式: %v", info["gamemode"])
	logx.Infof("维度: %v", info["dimension"])
	pos := info["position"].([]float64)
	rot := info["rotation"].([]float32)
	logx.Infof("位置: X=%.2f, Y=%.2f, Z=%.2f", pos[0], pos[1], pos[2])
	logx.Infof("朝向: Yaw=%.2f, Pitch=%.2f", rot[0], rot[1])
	logx.Infof("生命值: %.1f / %.1f", info["health"], info["max_health"])
	logx.Infof("饥饿值: %d (饱和度: %.1f)", info["food"], info["saturation"])
	logx.Infof("氧气值: %d / 300", info["air"])
	logx.Infof("实体生命: %.1f", info["entity_health"])
	logx.Infof("等级: %d, 经验: %.1f%%", info["level"], info["experience"])
	logx.Infof("手持槽位: %d", info["held_slot"])
	logx.Infof("飞行中: %v, 可飞行: %v", info["flying"], info["can_fly"])
	logx.Infof("在线时长: %v", info["duration"])

	inventory := c.Player.Inventory.GetAll()
	logx.Infof("=== 背包内容 ===")
	if len(inventory) == 0 {
		logx.Infof("背包为空（等待3秒后重试...）")
		time.Sleep(3 * time.Second)
		inventory = c.Player.Inventory.GetAll()
		if len(inventory) == 0 {
			logx.Infof("背包确实为空")
		} else {
			for slot, item := range inventory {
				if item != nil {
					logx.Infof("槽位 %d: %s x%d", slot, item.ID, item.Count)
				}
			}
		}
	} else {
		for slot, item := range inventory {
			if item != nil {
				logx.Infof("槽位 %d: %s x%d", slot, item.ID, item.Count)
			}
		}
	}
	logx.Infof("==================")
}

func (c *Client) initSecureChatSession() error {
	if c.chatSessionOK {
		return nil
	}
	if c.online == nil || strings.TrimSpace(c.online.AccessToken) == "" {
		return nil
	}

	cert, err := mcauth.GetPlayerCertificates(c.online.AccessToken)
	if err != nil {
		return err
	}

	privateKey, err := parsePrivateKeyPEM(cert.KeyPair.PrivateKey)
	if err != nil {
		return fmt.Errorf("解析证书 private key 失败: %w", err)
	}

	pubDER, err := parsePublicKeyPEMToSPKIDER(cert.KeyPair.PublicKey)
	if err != nil {
		return fmt.Errorf("解析证书 public key 失败: %w", err)
	}

	sigB64 := strings.TrimSpace(cert.PublicKeySignatureV2)
	if sigB64 == "" {
		sigB64 = strings.TrimSpace(cert.PublicKeySignature)
	}
	if sigB64 == "" {
		return fmt.Errorf("证书缺少 public key signature")
	}
	sig, err := base64.StdEncoding.DecodeString(sigB64)
	if err != nil {
		return fmt.Errorf("解析 public key signature 失败: %w", err)
	}

	sessionID, err := randomUUIDv4()
	if err != nil {
		return err
	}
	expireMillis := cert.ExpiresAt.UnixMilli()
	if cert.ExpiresAt.IsZero() {
		expireMillis = time.Now().Add(2 * time.Hour).UnixMilli()
	}

	payload := make([]byte, 0, len(pubDER)+len(sig)+64)
	payload = append(payload, sessionID[:]...)
	payload = append(payload, encodeInt64(expireMillis)...)
	payload = append(payload, encodeByteArray(pubDER)...)
	payload = append(payload, encodeByteArray(sig)...)

	if err := c.conn.WritePacket(playServerChatSession, payload); err != nil {
		return fmt.Errorf("发送 chat_session_update 失败: %w", err)
	}

	c.chatSignMu.Lock()
	c.chatSession = &secureChatSession{
		sessionID:    sessionID,
		privateKey:   privateKey,
		messageIndex: 0,
	}
	c.chatSignMu.Unlock()

	c.chatSessionOK = true
	logx.Infof("已发送 chat_session_update，启用 secure chat 公钥会话")
	return nil
}

func (c *Client) emitChat(chat ChatMessage) {
	if strings.TrimSpace(chat.RawJSON) != "" {
		logx.Infof("[聊天] %s", chat.RawJSON)
	}

	if c.chatHandler != nil {
		c.chatHandler(chat)
	}
}

func randomInt64() (int64, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0, fmt.Errorf("生成随机 salt 失败: %w", err)
	}
	return int64(binary.BigEndian.Uint64(b[:])), nil
}

func randomUUIDv4() ([16]byte, error) {
	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		return [16]byte{}, fmt.Errorf("生成 session UUID 失败: %w", err)
	}
	id[6] = (id[6] & 0x0F) | 0x40
	id[8] = (id[8] & 0x3F) | 0x80
	return id, nil
}

func (c *Client) sendUnsignedCommand(cmd string) error {
	payload := encodeString(cmd)
	logx.Debugf("命令包内容 (hex): %x", payload)
	return c.conn.WritePacket(playServerChatCommand, payload)
}

func (c *Client) sendSignedCommand(cmd string) error {
	timestamp := time.Now().UnixMilli()
	salt, err := randomInt64()
	if err != nil {
		return fmt.Errorf("生成命令签名 salt 失败: %w", err)
	}

	argSignatures, err := c.buildCommandArgumentSignatures(cmd, timestamp, salt)
	if err != nil {
		return err
	}

	payload := make([]byte, 0, len(cmd)+64)
	payload = append(payload, encodeString(cmd)...)
	payload = append(payload, encodeInt64(timestamp)...)
	payload = append(payload, encodeInt64(salt)...)
	payload = append(payload, encodeVarInt(int32(len(argSignatures)))...)
	for _, sig := range argSignatures {
		if len(sig.signature) != 256 {
			return fmt.Errorf("命令参数签名长度无效: %d", len(sig.signature))
		}
		payload = append(payload, encodeString(sig.name)...)
		payload = append(payload, sig.signature...)
	}
	payload = append(payload, encodeVarInt(0)...) // messageCount
	payload = append(payload, 0x00, 0x00, 0x00)   // acknowledged (3 bytes)
	payload = append(payload, 0x01)               // checksum (u8, 1 when no seen signatures)

	return c.conn.WritePacket(playServerChatCommandSign, payload)
}

func (c *Client) buildCommandArgumentSignatures(cmd string, timestampMillis int64, salt int64) ([]commandArgumentSignature, error) {
	verb, _ := splitCommand(cmd)
	target, ok := c.commandSignTarget(verb)
	if !ok {
		return nil, nil
	}

	parts := strings.Split(cmd, " ")
	if len(parts) <= target.SliceIndex {
		return nil, nil
	}
	arg := strings.Join(parts[target.SliceIndex:], " ")
	if strings.TrimSpace(arg) == "" {
		return nil, nil
	}

	sig, err := c.signChatBody(arg, timestampMillis, salt, nil)
	if err != nil {
		return nil, fmt.Errorf("生成命令参数签名失败: %w", err)
	}
	return []commandArgumentSignature{
		{
			name:      target.ArgumentName,
			signature: sig,
		},
	}, nil
}

func (c *Client) commandSignTarget(verb string) (signableCommandTarget, bool) {
	c.chatSignMu.Lock()
	defer c.chatSignMu.Unlock()
	target, ok := c.commandSign[verb]
	return target, ok
}

func splitCommand(cmd string) (verb string, arg string) {
	parts := strings.SplitN(strings.TrimSpace(cmd), " ", 2)
	if len(parts) == 0 {
		return "", ""
	}
	verb = strings.ToLower(strings.TrimSpace(parts[0]))
	if len(parts) == 1 {
		return verb, ""
	}
	return verb, strings.TrimSpace(parts[1])
}

func (c *Client) signChatBody(content string, timestampMillis int64, salt int64, acknowledgements [][]byte) ([]byte, error) {
	c.chatSignMu.Lock()
	defer c.chatSignMu.Unlock()

	if c.chatSession == nil || c.chatSession.privateKey == nil {
		return nil, fmt.Errorf("secure chat 会话未初始化")
	}

	contentBytes := []byte(content)
	signable := make([]byte, 0, len(contentBytes)+128+len(acknowledgements)*256)
	signable = append(signable, encodeInt32(1)...)
	signable = append(signable, c.uuid[:]...)
	signable = append(signable, c.chatSession.sessionID[:]...)
	signable = append(signable, encodeInt32(c.chatSession.messageIndex)...)
	c.chatSession.messageIndex++
	signable = append(signable, encodeInt64(salt)...)
	signable = append(signable, encodeInt64(timestampMillis/1000)...)
	signable = append(signable, encodeInt32(int32(len(contentBytes)))...)
	// 签名体这里是 Int32 length + 原始 UTF-8 内容，不是协议 String(VarInt length)。
	signable = append(signable, contentBytes...)
	signable = append(signable, encodeInt32(int32(len(acknowledgements)))...)
	for _, ack := range acknowledgements {
		signable = append(signable, ack...)
	}

	sum := sha256.Sum256(signable)
	sig, err := rsa.SignPKCS1v15(rand.Reader, c.chatSession.privateKey, crypto.SHA256, sum[:])
	if err != nil {
		return nil, err
	}
	return sig, nil
}

type commandNodeWire struct {
	NodeType byte
	Name     string
	ParserID int32
	Children []int32
}

func extractSignableCommandTargets(nodes []commandNodeWire, rootIndex int32) map[string]signableCommandTarget {
	targets := make(map[string]signableCommandTarget)
	if rootIndex < 0 || int(rootIndex) >= len(nodes) {
		return targets
	}

	root := nodes[rootIndex]
	for _, literalIdx := range root.Children {
		if literalIdx < 0 || int(literalIdx) >= len(nodes) {
			continue
		}
		literal := nodes[literalIdx]
		if literal.NodeType != 1 || literal.Name == "" {
			continue
		}

		pathSeen := map[int32]bool{}
		var walk func(idx int32, depth int)
		walk = func(idx int32, depth int) {
			if idx < 0 || int(idx) >= len(nodes) || pathSeen[idx] {
				return
			}
			pathSeen[idx] = true
			node := nodes[idx]

			if node.NodeType == 2 && (node.ParserID == 5 || node.ParserID == 20) {
				target := signableCommandTarget{
					ArgumentName: node.Name,
					SliceIndex:   depth,
				}
				old, exists := targets[literal.Name]
				if !exists || target.SliceIndex < old.SliceIndex {
					targets[literal.Name] = target
				}
			}

			for _, child := range node.Children {
				walk(child, depth+1)
			}
			delete(pathSeen, idx)
		}

		walk(literalIdx, 0)
	}

	return targets
}

func skipCommandParserProperties(r *bytes.Reader, parserID int32) error {
	switch parserID {
	case 1: // brigadier:float
		flags, err := readU8(r)
		if err != nil {
			return err
		}
		if flags&0x01 != 0 {
			if err := discardN(r, 4); err != nil {
				return err
			}
		}
		if flags&0x02 != 0 {
			if err := discardN(r, 4); err != nil {
				return err
			}
		}
	case 2: // brigadier:double
		flags, err := readU8(r)
		if err != nil {
			return err
		}
		if flags&0x01 != 0 {
			if err := discardN(r, 8); err != nil {
				return err
			}
		}
		if flags&0x02 != 0 {
			if err := discardN(r, 8); err != nil {
				return err
			}
		}
	case 3: // brigadier:integer
		flags, err := readU8(r)
		if err != nil {
			return err
		}
		if flags&0x01 != 0 {
			if err := discardN(r, 4); err != nil {
				return err
			}
		}
		if flags&0x02 != 0 {
			if err := discardN(r, 4); err != nil {
				return err
			}
		}
	case 4: // brigadier:long
		flags, err := readU8(r)
		if err != nil {
			return err
		}
		if flags&0x01 != 0 {
			if err := discardN(r, 8); err != nil {
				return err
			}
		}
		if flags&0x02 != 0 {
			if err := discardN(r, 8); err != nil {
				return err
			}
		}
	case 5: // brigadier:string
		if _, err := readVarInt(r); err != nil {
			return err
		}
	case 6: // minecraft:entity
		if _, err := readU8(r); err != nil {
			return err
		}
	case 31: // minecraft:score_holder
		if _, err := readU8(r); err != nil {
			return err
		}
	case 43: // minecraft:time
		if _, err := readInt32(r); err != nil {
			return err
		}
	case 44, 45, 46, 47, 48: // resource* parsers
		if _, err := readString(r, r); err != nil {
			return err
		}
	}
	return nil
}

func discardN(r io.Reader, n int64) error {
	_, err := io.CopyN(io.Discard, r, n)
	return err
}

func readAnonymousNBTJSON(r io.Reader) (string, error) {
	dec := nbt.NewDecoder(r)
	dec.NetworkFormat(true)
	var v any
	if err := dec.Decode(&v); err != nil {
		return "", err
	}

	raw, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func parsePublicKeyPEMToSPKIDER(publicKeyPEM string) ([]byte, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("无效 PEM 公钥")
	}

	if pubAny, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		return x509.MarshalPKIXPublicKey(pubAny)
	}
	if rsaPub, err := x509.ParsePKCS1PublicKey(block.Bytes); err == nil {
		return x509.MarshalPKIXPublicKey(rsaPub)
	}
	return nil, fmt.Errorf("不支持的 PEM 公钥类型: %s", block.Type)
}

func parsePrivateKeyPEM(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("无效 PEM 私钥")
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	if anyKey, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if key, ok := anyKey.(*rsa.PrivateKey); ok {
			return key, nil
		}
		return nil, fmt.Errorf("不支持的 PKCS8 私钥类型")
	}
	return nil, fmt.Errorf("不支持的 PEM 私钥类型: %s", block.Type)
}
