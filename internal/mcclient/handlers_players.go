package mcclient

import (
	"bytes"
	"fmt"

	"gmcc/internal/logx"
	"gmcc/internal/mcclient/packet"
)

func (c *Client) handlePlayerInfoUpdate(data []byte) error {
	r := bytes.NewReader(data)

	// 协议 774: action 是单字节位域，不是 varint
	action, err := packet.ReadU8(r)
	if err != nil {
		return fmt.Errorf("读取 player_info_update action 失败: %w", err)
	}

	count, err := packet.ReadVarInt(r)
	if err != nil {
		return fmt.Errorf("读取 player_info_update count 失败: %w", err)
	}

	logx.Debugf("player_info_update: action=0x%02x, count=%d", action, count)

	c.playersMu.Lock()
	defer c.playersMu.Unlock()

	for i := int32(0); i < count; i++ {
		uuid, err := packet.ReadUUID(r)
		if err != nil {
			return fmt.Errorf("读取玩家 UUID 失败: %w", err)
		}

		var playerName string
		// 协议 774 位定义: bit7=add_player, bit6=init_chat, bit5=game_mode, bit4=listed, bit3=latency, bit2=display_name, bit1=list_order, bit0=show_hat
		if action&0x80 != 0 {
			playerName = packet.MustReadString(r, "player_info_update.name")
			propertiesCount := packet.MustReadVarInt(r, "player_info_update.properties_count")

			for j := int32(0); j < propertiesCount; j++ {
				_ = packet.MustReadString(r, "player_info_update.property_name")
				_ = packet.MustReadString(r, "player_info_update.property_value")
				isSigned := packet.MustReadBool(r, "player_info_update.property_is_signed")
				if isSigned {
					_ = packet.MustReadString(r, "player_info_update.property_signature")
				}
			}

			c.players[playerName] = playerInfo{uuid: uuid}

			logx.Infof("玩家信息更新: 添加玩家 %s (%s)", playerName, formatUUIDShort(uuid))
		} else {
			playerName = c.findPlayerNameByUUID(uuid)
			if playerName == "" {
				logx.Warnf("玩家信息更新: 收到更新但未知玩家 UUID %s (action=0x%02x)", formatUUIDShort(uuid), action)
			}
		}

		if action&0x40 != 0 {
			// initialize_chat: optional chat session
			hasSession := packet.MustReadBool(r, "player_info_update.has_chat_session")
			if hasSession {
				_, _ = packet.ReadUUID(r) // session UUID
				_, _ = packet.ReadVarInt(r)
				keyLen := packet.MustReadVarInt(r, "player_info_update.key_len")
				_, _ = packet.ReadBytes(r, int(keyLen)) // public key
				sigLen := packet.MustReadVarInt(r, "player_info_update.sig_len")
				_, _ = packet.ReadBytes(r, int(sigLen)) // signature
			}
		}

		if action&0x20 != 0 {
			gamemode := packet.MustReadVarInt(r, "player_info_update.gamemode")
			if playerName != "" {
				logx.Debugf("玩家 %s 游戏模式更新: %d", playerName, gamemode)
			}
		}

		if action&0x10 != 0 {
			listed := packet.MustReadBool(r, "player_info_update.listed")
			if playerName != "" {
				logx.Debugf("玩家 %s 列表状态: %v", playerName, listed)
			}
		}

		if action&0x08 != 0 {
			latency := packet.MustReadVarInt(r, "player_info_update.latency")
			if playerName != "" {
				logx.Debugf("玩家 %s 延迟: %dms", playerName, latency)
			}
		}

		if action&0x04 != 0 {
			hasDisplayName := packet.MustReadBool(r, "player_info_update.has_display_name")
			if hasDisplayName {
				_ = packet.MustReadString(r, "player_info_update.display_name")
			}
		}

		if action&0x02 != 0 {
			_ = packet.MustReadVarInt(r, "player_info_update.list_order")
		}

		if action&0x01 != 0 {
			_ = packet.MustReadBool(r, "player_info_update.show_hat")
		}
	}

	return nil
}

func (c *Client) handlePlayerInfoRemove(data []byte) error {
	r := bytes.NewReader(data)

	count, err := packet.ReadVarInt(r)
	if err != nil {
		return fmt.Errorf("读取 player_info_remove count 失败: %w", err)
	}

	c.playersMu.Lock()
	defer c.playersMu.Unlock()

	for i := int32(0); i < count; i++ {
		uuid, err := packet.ReadUUID(r)
		if err != nil {
			return fmt.Errorf("读取玩家 UUID 失败: %w", err)
		}

		for name, info := range c.players {
			if info.uuid == uuid {
				delete(c.players, name)
				logx.Debugf("玩家移除: %s", name)
				break
			}
		}
	}

	return nil
}

func (c *Client) GetOnlinePlayers() []string {
	c.playersMu.RLock()
	defer c.playersMu.RUnlock()

	players := make([]string, 0, len(c.players))
	for name := range c.players {
		players = append(players, name)
	}

	return players
}

func (c *Client) logOnlinePlayers() {
	players := c.GetOnlinePlayers()
	if len(players) == 0 {
		logx.Infof("在线玩家: 无")
		return
	}

	playerList := players[0]
	if len(players) > 1 {
		playerList = fmt.Sprintf("%s 等 %d 人", playerList, len(players))
	}
	logx.Infof("在线玩家: %s", playerList)
}

func formatUUIDShort(uuid [16]byte) string {
	const hexChars = "0123456789abcdef"
	hex := make([]byte, 8)
	for i := 0; i < 4; i++ {
		hex[i*2] = hexChars[uuid[i]>>4]
		hex[i*2+1] = hexChars[uuid[i]&0x0f]
	}
	return string(hex)
}

func (c *Client) findPlayerNameByUUID(uuid [16]byte) string {
	for name, info := range c.players {
		if info.uuid == uuid {
			return name
		}
	}
	return ""
}
