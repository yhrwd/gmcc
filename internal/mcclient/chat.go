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

	"github.com/Tnze/go-mc/nbt"

	mcauth "gmcc/internal/auth/minecraft"
	"gmcc/internal/logx"
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

func (c *Client) GetCharacterAnalyzer() *CharacterAnalyzer {
	return c.charAnalyzer
}

func (c *Client) AddCharacterMapping(unicodeStr, description, replaceWith string) error {
	if c.charAnalyzer == nil {
		return fmt.Errorf("字符映射分析器未初始化")
	}
	return c.charAnalyzer.AddMapping(unicodeStr, description, replaceWith)
}

func (c *Client) RemoveCharacterMapping(unicodeStr string) error {
	if c.charAnalyzer == nil {
		return fmt.Errorf("字符映射分析器未初始化")
	}
	return c.charAnalyzer.RemoveMapping(unicodeStr)
}

func (c *Client) SetCharacterReplaceEnabled(enabled bool) error {
	if c.charAnalyzer == nil {
		return fmt.Errorf("字符映射分析器未初始化")
	}
	return c.charAnalyzer.SetEnableReplace(enabled)
}

func (c *Client) SetShowUnicodeInfo(show bool) error {
	if c.charAnalyzer == nil {
		return fmt.Errorf("字符映射分析器未初始化")
	}
	return c.charAnalyzer.SetShowUnicodeInfo(show)
}

func (c *Client) GenerateCharacterMappingTemplate(outputPath string) error {
	if c.charAnalyzer == nil {
		return fmt.Errorf("字符映射分析器未初始化")
	}
	return c.charAnalyzer.GenerateMappingTemplate(outputPath)
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

	// 1.21.11 使用 chat_command_signed。
	// 当前实现会对 /say 的 message 参数做真实 RSA 签名，其余命令暂走基础签名字段。
	if err := c.sendSignedCommand(cmd); err != nil {
		logx.Warnf("发送签名命令失败，回退无签名命令: %v", err)
		if err2 := c.sendUnsignedCommand(cmd); err2 != nil {
			return fmt.Errorf("发送命令失败: signedErr=%v, unsignedErr=%w", err, err2)
		}
		logx.Infof("已发送命令(unsigned fallback): /%s", cmd)
		return nil
	}
	logx.Infof("已发送命令(signed): /%s", cmd)
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

	for _, cmd := range c.cfg.Actions.OnJoinCommands {
		if err := c.SendCommand(cmd); err != nil {
			logx.Warnf("执行 on_join_commands 失败: %v", err)
		}
	}
	for _, msg := range c.cfg.Actions.OnJoinMessages {
		if err := c.SendMessage(msg); err != nil {
			logx.Warnf("执行 on_join_messages 失败: %v", err)
		}
	}
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

func (c *Client) handleSystemChatPacket(data []byte) error {
	r := bytes.NewReader(data)
	rawJSON, err := readAnonymousNBTJSON(r)
	if err != nil {
		return fmt.Errorf("解析 system_chat 内容失败: %w", err)
	}
	isActionBar, err := readBool(r)
	if err != nil {
		return fmt.Errorf("解析 system_chat actionBar 标记失败: %w", err)
	}

	chat := ChatMessage{
		Type:        "system",
		PlainText:   extractPlainTextFromChatJSON(rawJSON),
		RawJSON:     rawJSON,
		IsActionBar: isActionBar,
		ReceivedAt:  time.Now(),
	}
	if isActionBar {
		chat.Type = "action_bar"
	}
	c.emitChat(chat)
	return nil
}

func (c *Client) handleActionBarPacket(data []byte) error {
	r := bytes.NewReader(data)
	rawJSON, err := readAnonymousNBTJSON(r)
	if err != nil {
		return fmt.Errorf("解析 action_bar 内容失败: %w", err)
	}
	c.emitChat(ChatMessage{
		Type:        "action_bar",
		PlainText:   extractPlainTextFromChatJSON(rawJSON),
		RawJSON:     rawJSON,
		IsActionBar: true,
		ReceivedAt:  time.Now(),
	})
	return nil
}

func (c *Client) handleProfilelessChatPacket(data []byte) error {
	r := bytes.NewReader(data)
	rawJSON, err := readAnonymousNBTJSON(r)
	if err != nil {
		return fmt.Errorf("解析 profileless_chat 内容失败: %w", err)
	}
	c.emitChat(ChatMessage{
		Type:       "profileless_chat",
		PlainText:  extractPlainTextFromChatJSON(rawJSON),
		RawJSON:    rawJSON,
		ReceivedAt: time.Now(),
	})
	return nil
}

func (c *Client) handlePlayerChatPacket(data []byte) error {
	r := bytes.NewReader(data)
	if _, err := readVarInt(r); err != nil {
		return fmt.Errorf("解析 player_chat globalIndex 失败: %w", err)
	}
	sender, err := readUUID(r)
	if err != nil {
		return fmt.Errorf("解析 player_chat senderUuid 失败: %w", err)
	}
	if _, err := readVarInt(r); err != nil {
		return fmt.Errorf("解析 player_chat index 失败: %w", err)
	}

	hasSignature, err := readBool(r)
	if err != nil {
		return fmt.Errorf("解析 player_chat signature 标记失败: %w", err)
	}
	if hasSignature {
		if err := discardN(r, 256); err != nil {
			return fmt.Errorf("跳过 player_chat signature 失败: %w", err)
		}
	}

	plain, err := readString(r, r)
	if err != nil {
		return fmt.Errorf("解析 player_chat plainMessage 失败: %w", err)
	}
	if _, err := readInt64(r); err != nil {
		return fmt.Errorf("解析 player_chat timestamp 失败: %w", err)
	}
	if _, err := readInt64(r); err != nil {
		return fmt.Errorf("解析 player_chat salt 失败: %w", err)
	}

	prevCount, err := readVarInt(r)
	if err != nil {
		return fmt.Errorf("解析 player_chat previousMessages 数量失败: %w", err)
	}
	for i := int32(0); i < prevCount; i++ {
		id, err := readVarInt(r)
		if err != nil {
			return fmt.Errorf("解析 player_chat previousMessage.id 失败: %w", err)
		}
		if id == 0 {
			if err := discardN(r, 256); err != nil {
				return fmt.Errorf("跳过 player_chat previousMessage.signature 失败: %w", err)
			}
		}
	}

	var rawJSON string
	hasUnsigned, err := readBool(r)
	if err != nil {
		return fmt.Errorf("解析 player_chat unsignedChatContent 标记失败: %w", err)
	}
	if hasUnsigned {
		rawJSON, err = readAnonymousNBTJSON(r)
		if err != nil {
			return fmt.Errorf("解析 player_chat unsignedChatContent 失败: %w", err)
		}
	}

	c.emitChat(ChatMessage{
		Type:       "player_chat",
		PlainText:  plain,
		RawJSON:    rawJSON,
		SenderUUID: formatUUID(sender),
		ReceivedAt: time.Now(),
	})
	return nil
}

func (c *Client) emitChat(chat ChatMessage) {
	originalText := chat.PlainText

	// 应用字符替换
	if c.charAnalyzer != nil {
		chat.PlainText = c.charAnalyzer.ReplaceText(chat.PlainText)
	}

	if strings.TrimSpace(chat.PlainText) == "" && strings.TrimSpace(chat.RawJSON) != "" {
		logx.Infof("收到聊天消息: type=%s sender=%s plain=<empty> raw=%s", chat.Type, chat.SenderUUID, shortString(chat.RawJSON, 180))
	} else {
		logx.Infof("收到聊天消息: type=%s sender=%s plain=%q", chat.Type, chat.SenderUUID, chat.PlainText)
	}
	if strings.TrimSpace(chat.RawJSON) != "" {
		logx.Debugf("聊天 JSON: %s", chat.RawJSON)
	}

	// 显示Unicode字符信息
	if c.charAnalyzer != nil {
		if unicodeInfo := c.charAnalyzer.AnalyzeText(originalText); unicodeInfo != "" {
			logx.Infof("%s", unicodeInfo)
		}
	}

	// 检测未知字符（替换字符 �）
	if strings.Contains(originalText, "�") {
		logx.Warnf("检测到未知字符，原始JSON: %s", chat.RawJSON)
		logx.Debugf("PlainText with unknown chars: %q", originalText)
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
	return c.conn.WritePacket(playServerChatCommand, encodeString(cmd))
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

func (c *Client) handleDeclareCommandsPacket(data []byte) error {
	r := bytes.NewReader(data)
	nodeCount, err := readVarInt(r)
	if err != nil {
		return fmt.Errorf("解析 declare_commands 节点数失败: %w", err)
	}
	if nodeCount < 0 {
		return fmt.Errorf("declare_commands 节点数无效: %d", nodeCount)
	}

	nodes := make([]commandNodeWire, nodeCount)
	for i := int32(0); i < nodeCount; i++ {
		flags, err := readU8(r)
		if err != nil {
			return fmt.Errorf("解析 declare_commands[%d] flags 失败: %w", i, err)
		}
		nodeType := flags & 0x03
		hasRedirect := (flags & 0x08) != 0
		hasSuggestions := (flags & 0x10) != 0

		childCount, err := readVarInt(r)
		if err != nil {
			return fmt.Errorf("解析 declare_commands[%d] children count 失败: %w", i, err)
		}
		if childCount < 0 {
			return fmt.Errorf("declare_commands[%d] children count 无效: %d", i, childCount)
		}
		children := make([]int32, 0, childCount)
		for j := int32(0); j < childCount; j++ {
			child, err := readVarInt(r)
			if err != nil {
				return fmt.Errorf("解析 declare_commands[%d] child[%d] 失败: %w", i, j, err)
			}
			children = append(children, child)
		}

		if hasRedirect {
			if _, err := readVarInt(r); err != nil {
				return fmt.Errorf("解析 declare_commands[%d] redirect 失败: %w", i, err)
			}
		}

		node := commandNodeWire{
			NodeType: nodeType,
			Children: children,
		}

		switch nodeType {
		case 1: // literal
			name, err := readString(r, r)
			if err != nil {
				return fmt.Errorf("解析 declare_commands[%d] literal name 失败: %w", i, err)
			}
			node.Name = strings.ToLower(strings.TrimSpace(name))
		case 2: // argument
			name, err := readString(r, r)
			if err != nil {
				return fmt.Errorf("解析 declare_commands[%d] argument name 失败: %w", i, err)
			}
			parserID, err := readVarInt(r)
			if err != nil {
				return fmt.Errorf("解析 declare_commands[%d] parser 失败: %w", i, err)
			}
			if err := skipCommandParserProperties(r, parserID); err != nil {
				return fmt.Errorf("跳过 declare_commands[%d] parser=%d properties 失败: %w", i, parserID, err)
			}
			if hasSuggestions {
				if _, err := readString(r, r); err != nil {
					return fmt.Errorf("解析 declare_commands[%d] suggestions 失败: %w", i, err)
				}
			}
			node.Name = name
			node.ParserID = parserID
		}

		nodes[i] = node
	}

	rootIndex, err := readVarInt(r)
	if err != nil {
		return fmt.Errorf("解析 declare_commands rootIndex 失败: %w", err)
	}

	targets := extractSignableCommandTargets(nodes, rootIndex)
	c.chatSignMu.Lock()
	c.commandSign = targets
	c.chatSignMu.Unlock()

	logx.Debugf("已解析 declare_commands: signableCommands=%d", len(targets))
	if say, ok := targets["say"]; ok {
		logx.Debugf("命令签名规则: /say arg=%s sliceIndex=%d", say.ArgumentName, say.SliceIndex)
	}
	return nil
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

			if node.NodeType == 2 && node.ParserID == 20 {
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
	var v any
	dec := nbt.NewDecoder(r)
	dec.NetworkFormat(true)
	if _, err := dec.Decode(&v); err != nil {
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

func readU8(r io.Reader) (byte, error) {
	var b [1]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return 0, err
	}
	return b[0], nil
}

func shortString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
