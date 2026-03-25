package mcclient

import (
	"bytes"
	"fmt"

	"gmcc/internal/logx"
	"gmcc/internal/mcclient/packet"
)

func (c *Client) handlePlayerInfoUpdate(data []byte) error {
	r := bytes.NewReader(data)

	action, err := packet.ReadVarInt(r)
	if err != nil {
		return fmt.Errorf("读取 player_info_update action 失败: %w", err)
	}

	count, err := packet.ReadVarInt(r)
	if err != nil {
		return fmt.Errorf("读取 player_info_update count 失败: %w", err)
	}

	c.playersMu.Lock()
	defer c.playersMu.Unlock()

	for i := int32(0); i < count; i++ {
		uuid, err := packet.ReadUUID(r)
		if err != nil {
			return fmt.Errorf("读取玩家 UUID 失败: %w", err)
		}

		var playerName string
		if action&1 != 0 {
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

			logx.Debugf("玩家信息更新: 添加玩家 %s (%s)", playerName, formatUUIDShort(uuid))
		} else {
			playerName = c.findPlayerNameByUUID(uuid)
		}

		if action&2 != 0 {
			_ = packet.MustReadVarInt(r, "player_info_update.chat_session")
		}

		if action&4 != 0 {
			gamemode := packet.MustReadVarInt(r, "player_info_update.gamemode")
			if playerName != "" {
				logx.Debugf("玩家 %s 游戏模式更新: %d", playerName, gamemode)
			}
		}

		if action&8 != 0 {
			listed := packet.MustReadBool(r, "player_info_update.listed")
			if playerName != "" {
				logx.Debugf("玩家 %s 列表状态: %v", playerName, listed)
			}
		}

		if action&16 != 0 {
			latency := packet.MustReadVarInt(r, "player_info_update.latency")
			if playerName != "" {
				logx.Debugf("玩家 %s 延迟: %dms", playerName, latency)
			}
		}

		if action&32 != 0 {
			hasDisplayName := packet.MustReadBool(r, "player_info_update.has_display_name")
			if hasDisplayName {
				_ = packet.MustReadString(r, "player_info_update.display_name")
			}
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
