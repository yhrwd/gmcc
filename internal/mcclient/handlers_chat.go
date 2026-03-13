package mcclient

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"gmcc/internal/logx"
	"gmcc/internal/mcclient/chat"
	"gmcc/internal/mcclient/packet"
)

func (c *Client) handleSystemChatPacket(data []byte) error {
	r := bytes.NewReader(data)
	rawJSON, err := c.readAnonymousNBTJSON(r)
	if err != nil {
		return fmt.Errorf("解析 system_chat 内容失败: %w", err)
	}
	isActionBar, err := packet.ReadBool(r)
	if err != nil {
		return fmt.Errorf("解析 system_chat actionBar 标记失败: %w", err)
	}

	msg := ChatMessage{
		Type:        "system",
		PlainText:   chat.ExtractPlainTextFromChatJSON(rawJSON),
		RawJSON:     rawJSON,
		IsActionBar: isActionBar,
		ReceivedAt:  time.Now(),
	}
	if isActionBar {
		msg.Type = "action_bar"
	}
	c.emitChat(msg)
	return nil
}

func (c *Client) handleActionBarPacket(data []byte) error {
	r := bytes.NewReader(data)
	rawJSON, err := c.readAnonymousNBTJSON(r)
	if err != nil {
		return fmt.Errorf("解析 action_bar 内容失败: %w", err)
	}
	c.emitChat(ChatMessage{
		Type:        "action_bar",
		PlainText:   chat.ExtractPlainTextFromChatJSON(rawJSON),
		RawJSON:     rawJSON,
		IsActionBar: true,
		ReceivedAt:  time.Now(),
	})
	return nil
}

func (c *Client) handleProfilelessChatPacket(data []byte) error {
	r := bytes.NewReader(data)
	rawJSON, err := c.readAnonymousNBTJSON(r)
	if err != nil {
		return fmt.Errorf("解析 profileless_chat 内容失败: %w", err)
	}
	c.emitChat(ChatMessage{
		Type:       "profileless_chat",
		PlainText:  chat.ExtractPlainTextFromChatJSON(rawJSON),
		RawJSON:    rawJSON,
		ReceivedAt: time.Now(),
	})
	return nil
}

func (c *Client) handlePlayerChatPacket(data []byte) error {
	r := bytes.NewReader(data)
	if _, err := packet.ReadVarInt(r); err != nil {
		return fmt.Errorf("解析 player_chat globalIndex 失败: %w", err)
	}
	sender, err := packet.ReadUUID(r)
	if err != nil {
		return fmt.Errorf("解析 player_chat senderUuid 失败: %w", err)
	}
	if _, err := packet.ReadVarInt(r); err != nil {
		return fmt.Errorf("解析 player_chat index 失败: %w", err)
	}

	hasSignature, err := packet.ReadBool(r)
	if err != nil {
		return fmt.Errorf("解析 player_chat signature 标记失败: %w", err)
	}
	if hasSignature {
		if err := packet.DiscardN(r, 256); err != nil {
			return fmt.Errorf("跳过 player_chat signature 失败: %w", err)
		}
	}

	plain, err := packet.ReadString(r, r)
	if err != nil {
		return fmt.Errorf("解析 player_chat plainMessage 失败: %w", err)
	}
	if _, err := packet.ReadInt64(r); err != nil {
		return fmt.Errorf("解析 player_chat timestamp 失败: %w", err)
	}
	if _, err := packet.ReadInt64(r); err != nil {
		return fmt.Errorf("解析 player_chat salt 失败: %w", err)
	}

	prevCount, err := packet.ReadVarInt(r)
	if err != nil {
		return fmt.Errorf("解析 player_chat previousMessages 数量失败: %w", err)
	}
	for i := int32(0); i < prevCount; i++ {
		id, err := packet.ReadVarInt(r)
		if err != nil {
			return fmt.Errorf("解析 player_chat previousMessage.id 失败: %w", err)
		}
		if id == 0 {
			if err := packet.DiscardN(r, 256); err != nil {
				return fmt.Errorf("跳过 player_chat previousMessage.signature 失败: %w", err)
			}
		}
	}

	var rawJSON string
	hasUnsigned, err := packet.ReadBool(r)
	if err != nil {
		return fmt.Errorf("解析 player_chat unsignedChatContent 标记失败: %w", err)
	}
	if hasUnsigned {
		rawJSON, err = c.readAnonymousNBTJSON(r)
		if err != nil {
			return fmt.Errorf("解析 player_chat unsignedChatContent 失败: %w", err)
		}
	}

	c.emitChat(ChatMessage{
		Type:       "player_chat",
		PlainText:  plain,
		RawJSON:    rawJSON,
		SenderUUID: packet.FormatUUID(sender),
		ReceivedAt: time.Now(),
	})
	return nil
}

func (c *Client) handleDeclareCommandsPacket(data []byte) error {
	r := bytes.NewReader(data)
	nodeCount, err := packet.ReadVarInt(r)
	if err != nil {
		return fmt.Errorf("解析 declare_commands 节点数失败: %w", err)
	}
	if nodeCount < 0 {
		return fmt.Errorf("declare_commands 节点数无效: %d", nodeCount)
	}

	nodes := make([]commandNodeWire, nodeCount)
	for i := int32(0); i < nodeCount; i++ {
		flags, err := packet.ReadU8(r)
		if err != nil {
			return fmt.Errorf("解析 declare_commands[%d] flags 失败: %w", i, err)
		}
		nodeType := flags & 0x03
		hasRedirect := (flags & 0x08) != 0
		hasSuggestions := (flags & 0x10) != 0

		childCount, err := packet.ReadVarInt(r)
		if err != nil {
			return fmt.Errorf("解析 declare_commands[%d] children count 失败: %w", i, err)
		}
		if childCount < 0 {
			return fmt.Errorf("declare_commands[%d] children count 无效: %d", i, childCount)
		}
		children := make([]int32, 0, childCount)
		for j := int32(0); j < childCount; j++ {
			child, err := packet.ReadVarInt(r)
			if err != nil {
				return fmt.Errorf("解析 declare_commands[%d] child[%d] 失败: %w", i, j, err)
			}
			children = append(children, child)
		}

		if hasRedirect {
			if _, err := packet.ReadVarInt(r); err != nil {
				return fmt.Errorf("解析 declare_commands[%d] redirect 失败: %w", i, err)
			}
		}

		node := commandNodeWire{
			NodeType: nodeType,
			Children: children,
		}

		switch nodeType {
		case 1: // literal
			name, err := packet.ReadString(r, r)
			if err != nil {
				return fmt.Errorf("解析 declare_commands[%d] literal name 失败: %w", i, err)
			}
			node.Name = strings.ToLower(strings.TrimSpace(name))
		case 2: // argument
			name, err := packet.ReadString(r, r)
			if err != nil {
				return fmt.Errorf("解析 declare_commands[%d] argument name 失败: %w", i, err)
			}
			parserID, err := packet.ReadVarInt(r)
			if err != nil {
				return fmt.Errorf("解析 declare_commands[%d] parser 失败: %w", i, err)
			}
			if err := skipCommandParserProperties(r, parserID); err != nil {
				return fmt.Errorf("跳过 declare_commands[%d] parser=%d properties 失败: %w", i, parserID, err)
			}
			if hasSuggestions {
				if _, err := packet.ReadString(r, r); err != nil {
					return fmt.Errorf("解析 declare_commands[%d] suggestions 失败: %w", i, err)
				}
			}
			node.Name = name
			node.ParserID = parserID
		}

		nodes[i] = node
	}

	rootIndex, err := packet.ReadVarInt(r)
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
		flags, err := packet.ReadU8(r)
		if err != nil {
			return err
		}
		if flags&0x01 != 0 {
			if err := packet.DiscardN(r, 4); err != nil {
				return err
			}
		}
		if flags&0x02 != 0 {
			if err := packet.DiscardN(r, 4); err != nil {
				return err
			}
		}
	case 2: // brigadier:double
		flags, err := packet.ReadU8(r)
		if err != nil {
			return err
		}
		if flags&0x01 != 0 {
			if err := packet.DiscardN(r, 8); err != nil {
				return err
			}
		}
		if flags&0x02 != 0 {
			if err := packet.DiscardN(r, 8); err != nil {
				return err
			}
		}
	case 3: // brigadier:integer
		flags, err := packet.ReadU8(r)
		if err != nil {
			return err
		}
		if flags&0x01 != 0 {
			if err := packet.DiscardN(r, 4); err != nil {
				return err
			}
		}
		if flags&0x02 != 0 {
			if err := packet.DiscardN(r, 4); err != nil {
				return err
			}
		}
	case 4: // brigadier:long
		flags, err := packet.ReadU8(r)
		if err != nil {
			return err
		}
		if flags&0x01 != 0 {
			if err := packet.DiscardN(r, 8); err != nil {
				return err
			}
		}
		if flags&0x02 != 0 {
			if err := packet.DiscardN(r, 8); err != nil {
				return err
			}
		}
	case 5: // brigadier:string
		if _, err := packet.ReadVarInt(r); err != nil {
			return err
		}
	case 6: // minecraft:entity
		if _, err := packet.ReadU8(r); err != nil {
			return err
		}
	case 31: // minecraft:score_holder
		if _, err := packet.ReadU8(r); err != nil {
			return err
		}
	case 43: // minecraft:time
		if _, err := packet.ReadInt32(r); err != nil {
			return err
		}
	case 44, 45, 46, 47, 48: // resource* parsers
		if _, err := packet.ReadString(r, r); err != nil {
			return err
		}
	}
	return nil
}
