package mcclient

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"gmcc/internal/logx"
)

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

	logx.Debugf("已解析 declare_commands: nodes=%d signableCommands=%d", nodeCount, len(targets))
	for cmd, target := range targets {
		logx.Debugf("命令签名规则: /%s arg=%s sliceIndex=%d", cmd, target.ArgumentName, target.SliceIndex)
	}
	return nil
}
