package mcclient

import (
	"bytes"
	"fmt"

	"gmcc/internal/logx"
	"gmcc/internal/mcclient/packet"
)

const (
	playerInfoActionAddPlayer         byte = 0x01
	playerInfoActionInitializeChat    byte = 0x02
	playerInfoActionUpdateGameMode    byte = 0x04
	playerInfoActionUpdateListed      byte = 0x08
	playerInfoActionUpdateLatency     byte = 0x10
	playerInfoActionUpdateDisplayName byte = 0x20
	playerInfoActionUpdateListOrder   byte = 0x40
	playerInfoActionUpdateShowHat     byte = 0x80
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
		if action&playerInfoActionAddPlayer != 0 {
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

			logx.Debugf("列表添加玩家 %s (%s)", playerName, packet.FormatUUIDShort(uuid))
		} else {
			playerName = c.findPlayerNameByUUIDLocked(uuid)
			if playerName == "" {
				logx.Warnf("玩家信息更新: 收到更新但未知玩家 UUID %s (action=0x%02x)", packet.FormatUUIDShort(uuid), action)
			}
		}

		if action&playerInfoActionInitializeChat != 0 {
			hasSession := packet.MustReadBool(r, "player_info_update.has_chat_session")
			if hasSession {
				if _, err := packet.ReadUUID(r); err != nil {
					return fmt.Errorf("读取 player_info_update.chat_session_id 失败: %w", err)
				}
				if _, err := packet.ReadInt64(r); err != nil {
					return fmt.Errorf("读取 player_info_update.public_key_expiry 失败: %w", err)
				}
				if _, err := packet.ReadByteArray(r, r); err != nil {
					return fmt.Errorf("读取 player_info_update.public_key 失败: %w", err)
				}
				if _, err := packet.ReadByteArray(r, r); err != nil {
					return fmt.Errorf("读取 player_info_update.public_key_signature 失败: %w", err)
				}
			}
		}

		if action&playerInfoActionUpdateGameMode != 0 {
			_ = packet.MustReadVarInt(r, "player_info_update.gamemode")
		}

		if action&playerInfoActionUpdateListed != 0 {
			_ = packet.MustReadBool(r, "player_info_update.listed")
		}

		if action&playerInfoActionUpdateLatency != 0 {
			_ = packet.MustReadVarInt(r, "player_info_update.latency")
		}

		if action&playerInfoActionUpdateDisplayName != 0 {
			hasDisplayName := packet.MustReadBool(r, "player_info_update.has_display_name")
			if hasDisplayName {
				if _, err := c.readAnonymousNBTJSON(r); err != nil {
					return fmt.Errorf("读取 player_info_update.display_name 失败: %w", err)
				}
			}
		}

		if action&playerInfoActionUpdateListOrder != 0 {
			_ = packet.MustReadVarInt(r, "player_info_update.list_order")
		}

		if action&playerInfoActionUpdateShowHat != 0 {
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

func (c *Client) findPlayerNameByUUIDLocked(uuid [16]byte) string {
	for name, info := range c.players {
		if info.uuid == uuid {
			return name
		}
	}
	return ""
}

func (c *Client) findPlayerNameByUUID(uuid [16]byte) string {
	c.playersMu.RLock()
	defer c.playersMu.RUnlock()

	return c.findPlayerNameByUUIDLocked(uuid)
}
