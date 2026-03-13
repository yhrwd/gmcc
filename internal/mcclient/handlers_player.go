package mcclient

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"gmcc/internal/logx"
	"gmcc/internal/player"
)

func (c *Client) handleSetHealthPacket(data []byte) error {
	r := bytes.NewReader(data)
	if r.Len() < 12 {
		return nil
	}

	var health float32
	binary.Read(r, binary.BigEndian, &health)

	food, _ := readVarIntFromReader(r)

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
	level, _ := readVarIntFromReader(r)
	totalExp, _ := readVarIntFromReader(r)

	c.Player.UpdateExperience(int32(level), expBar, float32(totalExp))
	return nil
}

func (c *Client) handleSetHeldSlotPacket(data []byte) error {
	r := bytes.NewReader(data)
	slot, err := readVarIntFromReader(r)
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

	eventType, _ := readU8(r)
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
	entityID, err := readInt32FromReader(r)
	if err != nil {
		return fmt.Errorf("读取 entity_id 失败: %w", err)
	}
	c.Player.SetEntityID(entityID)

	if _, err := readBoolFromReader(r); err != nil {
		return fmt.Errorf("读取 is_hardcore 失败: %w", err)
	}
	return nil
}

func (c *Client) readLoginDimensions(r *bytes.Reader) error {
	numDimensions, err := readVarIntFromReader(r)
	if err != nil {
		return fmt.Errorf("读取 dimension count 失败: %w", err)
	}
	for i := int32(0); i < numDimensions; i++ {
		if _, err := readStringFromReader(r); err != nil {
			return fmt.Errorf("读取 dimension name 失败: %w", err)
		}
	}
	return nil
}

func (c *Client) readLoginWorldSettings(r *bytes.Reader) error {
	fields := []string{"max_players", "view_distance", "simulation_distance"}
	for _, field := range fields {
		if _, err := readVarIntFromReader(r); err != nil {
			return fmt.Errorf("读取 %s 失败: %w", field, err)
		}
	}

	boolFields := []string{"reduced_debug_info", "enable_respawn_screen", "do_limited_crafting"}
	for _, field := range boolFields {
		if _, err := readBoolFromReader(r); err != nil {
			return fmt.Errorf("读取 %s 失败: %w", field, err)
		}
	}
	return nil
}

func (c *Client) readLoginPlayerState(r *bytes.Reader) error {
	if _, err := readVarIntFromReader(r); err != nil {
		return fmt.Errorf("读取 dimension_type 失败: %w", err)
	}

	dimensionName, err := readStringFromReader(r)
	if err != nil {
		return fmt.Errorf("读取 dimension_name 失败: %w", err)
	}
	c.Player.SetDimension(dimensionName)

	var hashedSeed int64
	binary.Read(r, binary.BigEndian, &hashedSeed)

	gameMode, err := readU8(r)
	if err != nil {
		return fmt.Errorf("读取 game_mode 失败: %w", err)
	}
	c.Player.SetGameMode(player.GameMode(int(gameMode)))

	_, _ = readU8(r)             // prevGameMode
	_, _ = readBoolFromReader(r) // isDebug
	_, _ = readBoolFromReader(r) // isFlat
	return nil
}

func (c *Client) readLoginDeathLocation(r *bytes.Reader) error {
	hasDeathLocation, err := readBoolFromReader(r)
	if err != nil {
		return fmt.Errorf("读取 has_death_location 失败: %w", err)
	}
	if hasDeathLocation {
		if _, err := readStringFromReader(r); err != nil {
			return fmt.Errorf("读取 death_dimension 失败: %w", err)
		}
		var deathPos int64
		binary.Read(r, binary.BigEndian, &deathPos)
	}
	return nil
}

func (c *Client) readLoginMisc(r *bytes.Reader) error {
	_, _ = readVarIntFromReader(r) // portalCooldown
	_, _ = readVarIntFromReader(r) // seaLevel
	_, _ = readBoolFromReader(r)   // secureChatEnforced
	return nil
}

func (c *Client) handleEntityDataPacket(data []byte) error {
	r := bytes.NewReader(data)
	entityID, err := readVarIntFromReader(r)
	if err != nil || entityID != c.Player.EntityID {
		return nil
	}

	for {
		index, err := readU8(r)
		if err != nil || index == 0xFF {
			break
		}
		typeID, err := readVarIntFromReader(r)
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
	windowID, err := readVarIntFromReader(r)
	if err != nil {
		return nil
	}
	stateID, _ := readVarIntFromReader(r)
	numItems, _ := readVarIntFromReader(r)
	if numItems > 1000 { // Safety limit
		numItems = 1000
	}

	logx.Debugf("container_content: windowID=%d, stateID=%d, numItems=%d", windowID, stateID, numItems)

	items := make([]*player.SlotData, numItems)
	for i := int32(0); i < numItems; i++ {
		slot, err := readSlotData(r)
		if err == nil && slot != nil {
			items[i] = &player.SlotData{ID: slot.ID, Count: slot.Count}
		}
	}

	carriedItem, _ := readSlotData(r)
	var carried *player.SlotData
	if carriedItem != nil {
		carried = &player.SlotData{ID: carriedItem.ID, Count: carriedItem.Count}
	}

	c.Player.UpdateInventory(windowID, items, carried)
	return nil
}

func (c *Client) handleContainerSlotPacket(data []byte) error {
	r := bytes.NewReader(data)
	windowID, _ := readVarIntFromReader(r)
	stateID, _ := readVarIntFromReader(r)

	var slot int16
	binary.Read(r, binary.BigEndian, &slot)

	item, _ := readSlotData(r)
	var slotItem *player.SlotData
	if item != nil {
		slotItem = &player.SlotData{ID: item.ID, Count: item.Count}
	}

	logx.Debugf("container_slot: windowID=%d, stateID=%d, slot=%d, item=%+v", windowID, stateID, slot, item)
	c.Player.UpdateSlot(windowID, int32(slot), slotItem)
	return nil
}
