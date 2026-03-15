package mcclient

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"gmcc/internal/logx"
	"gmcc/internal/mcclient/packet"
	"gmcc/internal/player"
)

func (c *Client) handleSetHealthPacket(data []byte) error {
	r := bytes.NewReader(data)
	if r.Len() < 12 {
		return nil
	}

	var health float32
	binary.Read(r, binary.BigEndian, &health)

	food, _ := packet.ReadVarIntFromReader(r)

	var saturation float32
	binary.Read(r, binary.BigEndian, &saturation)

	c.Player.UpdateHealth(health, 0, int32(food), saturation)
	return nil
}

func (c *Client) handleSetExperiencePacket(data []byte) error {
	r := bytes.NewReader(data)
	if r.Len() < 9 {
		return nil
	}

	var expBar float32
	binary.Read(r, binary.BigEndian, &expBar)
	level, _ := packet.ReadVarIntFromReader(r)
	totalExp, _ := packet.ReadVarIntFromReader(r)

	c.Player.UpdateExperience(int32(level), expBar, float32(totalExp))
	return nil
}

func (c *Client) handleSetHeldSlotPacket(data []byte) error {
	r := bytes.NewReader(data)
	slot, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return nil
	}
	c.Player.SetHeldSlot(int8(slot))
	return nil
}

func (c *Client) handleGameEventPacket(data []byte) error {
	r := bytes.NewReader(data)
	if r.Len() < 5 {
		return nil
	}

	eventType, _ := packet.ReadU8(r)
	var value float32
	binary.Read(r, binary.BigEndian, &value)

	if eventType == 3 && value >= 0 && value <= 3 {
		mode := player.GameMode(int(value))
		c.Player.SetGameMode(mode)
		logx.Infof("游戏模式变更: %s", mode.String())
	}

	return nil
}

func (c *Client) handlePlayLoginPacket(data []byte) error {
	r := bytes.NewReader(data)

	if err := c.readLoginBasicInfo(r); err != nil {
		return err
	}
	if err := c.readLoginDimensions(r); err != nil {
		return err
	}
	if err := c.readLoginWorldSettings(r); err != nil {
		return err
	}
	if err := c.readLoginPlayerState(r); err != nil {
		return err
	}
	if err := c.readLoginDeathLocation(r); err != nil {
		return err
	}
	if err := c.readLoginMisc(r); err != nil {
		return err
	}

	logx.Infof("登录Play阶段: EntityID=%d, 维度=%s, 游戏模式=%s", c.Player.EntityID, c.Player.Dimension, c.Player.GameMode.String())
	return nil
}

func (c *Client) readLoginBasicInfo(r *bytes.Reader) error {
	entityID, err := packet.ReadInt32FromReader(r)
	if err != nil {
		return fmt.Errorf("读取 entity_id 失败: %w", err)
	}
	c.Player.SetEntityID(entityID)

	if _, err := packet.ReadBoolFromReader(r); err != nil {
		return fmt.Errorf("读取 is_hardcore 失败: %w", err)
	}
	return nil
}

func (c *Client) readLoginDimensions(r *bytes.Reader) error {
	numDimensions, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return fmt.Errorf("读取 dimension count 失败: %w", err)
	}
	for i := int32(0); i < numDimensions; i++ {
		if _, err := packet.ReadStringFromReader(r); err != nil {
			return fmt.Errorf("读取 dimension name 失败: %w", err)
		}
	}
	return nil
}

func (c *Client) readLoginWorldSettings(r *bytes.Reader) error {
	fields := []string{"max_players", "view_distance", "simulation_distance"}
	for _, field := range fields {
		if _, err := packet.ReadVarIntFromReader(r); err != nil {
			return fmt.Errorf("读取 %s 失败: %w", field, err)
		}
	}

	boolFields := []string{"reduced_debug_info", "enable_respawn_screen", "do_limited_crafting"}
	for _, field := range boolFields {
		if _, err := packet.ReadBoolFromReader(r); err != nil {
			return fmt.Errorf("读取 %s 失败: %w", field, err)
		}
	}
	return nil
}

func (c *Client) readLoginPlayerState(r *bytes.Reader) error {
	if _, err := packet.ReadVarIntFromReader(r); err != nil {
		return fmt.Errorf("读取 dimension_type 失败: %w", err)
	}

	dimensionName, err := packet.ReadStringFromReader(r)
	if err != nil {
		return fmt.Errorf("读取 dimension_name 失败: %w", err)
	}
	c.Player.SetDimension(dimensionName)

	var hashedSeed int64
	binary.Read(r, binary.BigEndian, &hashedSeed)

	gameMode, err := packet.ReadU8(r)
	if err != nil {
		return fmt.Errorf("读取 game_mode 失败: %w", err)
	}
	c.Player.SetGameMode(player.GameMode(int(gameMode)))

	_, _ = packet.ReadU8(r)             // prevGameMode
	_, _ = packet.ReadBoolFromReader(r) // isDebug
	_, _ = packet.ReadBoolFromReader(r) // isFlat
	return nil
}

func (c *Client) readLoginDeathLocation(r *bytes.Reader) error {
	hasDeathLocation, err := packet.ReadBoolFromReader(r)
	if err != nil {
		return fmt.Errorf("读取 has_death_location 失败: %w", err)
	}
	if hasDeathLocation {
		if _, err := packet.ReadStringFromReader(r); err != nil {
			return fmt.Errorf("读取 death_dimension 失败: %w", err)
		}
		var deathPos int64
		binary.Read(r, binary.BigEndian, &deathPos)
	}
	return nil
}

func (c *Client) readLoginMisc(r *bytes.Reader) error {
	_, _ = packet.ReadVarIntFromReader(r) // portalCooldown
	_, _ = packet.ReadVarIntFromReader(r) // seaLevel
	_, _ = packet.ReadBoolFromReader(r)   // secureChatEnforced
	return nil
}

func (c *Client) handleEntityDataPacket(data []byte) error {
	r := bytes.NewReader(data)
	entityID, err := packet.ReadVarIntFromReader(r)
	if err != nil || entityID != c.Player.EntityID {
		return nil
	}

	for {
		index, err := packet.ReadU8(r)
		if err != nil || index == 0xFF {
			break
		}
		typeID, err := packet.ReadVarIntFromReader(r)
		if err != nil {
			break
		}
		_ = typeID
		// In a real client we'd parse the metadata here
		// For now we just skip the packet if it's too complex or just ignore it
		break
	}
	return nil
}

func (c *Client) handleContainerContentPacket(data []byte) error {
	r := bytes.NewReader(data)
	windowID, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		logx.Warnf("container_content: 读取 windowID 失败: %v", err)
		return nil
	}
	stateID, _ := packet.ReadVarIntFromReader(r)
	numItems, _ := packet.ReadVarIntFromReader(r)

	logx.Infof("container_content: windowID=%d, stateID=%d, numItems=%d, remaining=%d bytes", windowID, stateID, numItems, r.Len())

	if numItems > 1000 {
		logx.Warnf("container_content: numItems 过大 (%d), 限制为 1000", numItems)
		numItems = 1000
	}

	items := make([]*player.SlotData, numItems)
	for i := int32(0); i < numItems; i++ {
		slot, err := packet.ReadSlotData(r)
		if err != nil {
			logx.Warnf("container_content: 读取 slot %d 失败: %v, 剩余 %d 字节", i, err, r.Len())
			break
		}
		if slot != nil {
			items[i] = &player.SlotData{ID: slot.ID, Count: slot.Count}
			if i < 10 || slot.ID != 0 {
				logx.Debugf("  slot[%d]: id=%d, count=%d", i, slot.ID, slot.Count)
			}
		}
	}

	carriedItem, _ := packet.ReadSlotData(r)
	var carried *player.SlotData
	if carriedItem != nil {
		carried = &player.SlotData{ID: carriedItem.ID, Count: carriedItem.Count}
		logx.Infof("container_content: carried item: id=%d, count=%d", carriedItem.ID, carriedItem.Count)
	}

	c.Player.UpdateInventory(windowID, items, carried)
	return nil
}

func (c *Client) handleContainerSlotPacket(data []byte) error {
	r := bytes.NewReader(data)
	windowID, _ := packet.ReadVarIntFromReader(r)
	stateID, _ := packet.ReadVarIntFromReader(r)

	var slot int16
	binary.Read(r, binary.BigEndian, &slot)

	item, err := packet.ReadSlotData(r)
	if err != nil {
		logx.Warnf("container_slot: 读取物品数据失败: windowID=%d, slot=%d, err=%v", windowID, slot, err)
		return nil
	}

	var slotItem *player.SlotData
	if item != nil {
		slotItem = &player.SlotData{ID: item.ID, Count: item.Count}
		logx.Infof("container_slot: windowID=%d, stateID=%d, slot=%d, item_id=%d, count=%d", windowID, stateID, slot, item.ID, item.Count)
	} else {
		logx.Debugf("container_slot: windowID=%d, stateID=%d, slot=%d, item=empty", windowID, stateID, slot)
	}

	c.Player.UpdateSlot(windowID, int32(slot), slotItem)
	return nil
}
