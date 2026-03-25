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

	var health float32
	if err := binary.Read(r, binary.BigEndian, &health); err != nil {
		logx.PacketError("set_health", data, err)
		return nil
	}

	food, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		logx.PacketError("set_health", data, fmt.Errorf("读取 food 失败: %w", err))
		return nil
	}

	var saturation float32
	if err := binary.Read(r, binary.BigEndian, &saturation); err != nil {
		logx.PacketError("set_health", data, fmt.Errorf("读取 saturation 失败: %w", err))
		return nil
	}

	c.Player.UpdateHealth(health, 0, int32(food), saturation)
	return nil
}

func (c *Client) handleSetExperiencePacket(data []byte) error {
	r := bytes.NewReader(data)
	if r.Len() < 6 {
		logx.PacketError("set_experience", data, fmt.Errorf("数据过短: %d bytes (需要至少6)", r.Len()))
		return nil
	}

	var expBar float32
	if err := binary.Read(r, binary.BigEndian, &expBar); err != nil {
		logx.PacketError("set_experience", data, err)
		return nil
	}
	level, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		logx.PacketError("set_experience", data, fmt.Errorf("读取 level 失败: %w", err))
		return nil
	}
	totalExp, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		logx.PacketError("set_experience", data, fmt.Errorf("读取 totalExp 失败: %w", err))
		return nil
	}

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

	eventType := packet.MustReadU8(r, "player_info.eventType")
	var value float32
	_ = binary.Read(r, binary.BigEndian, &value)

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

	_ = packet.MustReadU8(r, "login.prevGameMode") // prevGameMode
	_ = packet.MustReadBool(r, "login.isDebug")    // isDebug
	_ = packet.MustReadBool(r, "login.isFlat")     // isFlat
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
	_ = packet.MustReadVarInt(r, "login.portalCooldown")   // portalCooldown
	_ = packet.MustReadVarInt(r, "login.seaLevel")         // seaLevel
	_ = packet.MustReadBool(r, "login.secureChatEnforced") // secureChatEnforced
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
		break
	}
	return nil
}

func (c *Client) handlePlayerAbilitiesPacket(data []byte) error {
	r := bytes.NewReader(data)

	flags, err := packet.ReadU8(r)
	if err != nil {
		logx.Warnf("player_abilities: 读取 flags 失败: %v", err)
		return nil
	}

	flyingSpeed, err := packet.ReadFloat32FromReader(r)
	if err != nil {
		logx.Warnf("player_abilities: 读取 flyingSpeed 失败: %v", err)
		return nil
	}

	walkingSpeed, err := packet.ReadFloat32FromReader(r)
	if err != nil {
		logx.Warnf("player_abilities: 读取 walkingSpeed 失败: %v", err)
		return nil
	}

	invulnerable := (flags & 0x01) != 0
	flying := (flags & 0x02) != 0
	canFly := (flags & 0x04) != 0
	creativeMode := (flags & 0x08) != 0

	c.Player.UpdateAbilities(int8(flags), flyingSpeed, walkingSpeed)

	logx.Infof("player_abilities: invulnerable=%v, flying=%v, canFly=%v, creativeMode=%v, flyingSpeed=%.2f, walkingSpeed=%.2f",
		invulnerable, flying, canFly, creativeMode, flyingSpeed, walkingSpeed)
	return nil
}
