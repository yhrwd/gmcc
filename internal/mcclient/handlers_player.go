package mcclient

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"gmcc/internal/logx"
	"gmcc/internal/nbt"
	"gmcc/internal/player"
)

func (c *Client) handleSetHealthPacket(data []byte) error {
	logx.Debugf("set_health raw: len=%d, hex=%x", len(data), data)

	r := bytes.NewReader(data)
	if r.Len() < 12 {
		logx.Warnf("set_health 数据太短: len=%d", len(data))
		return nil
	}

	var health float32
	binary.Read(r, binary.BigEndian, &health)

	food, _ := readVarIntFromReader(r)

	var saturation float32
	binary.Read(r, binary.BigEndian, &saturation)

	logx.Debugf("set_health: health=%.1f, food=%d, saturation=%.1f", health, food, saturation)

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
	logx.Debugf("set_held_slot: slot=%d, raw=%x", slot, data)
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

	entityID, err := readInt32FromReader(r)
	if err != nil {
		return fmt.Errorf("读取 entity_id 失败: %w", err)
	}
	c.Player.SetEntityID(entityID)

	isHardcore, err := readBoolFromReader(r)
	if err != nil {
		return fmt.Errorf("读取 is_hardcore 失败: %w", err)
	}
	_ = isHardcore

	numDimensions, err := readVarIntFromReader(r)
	if err != nil {
		return fmt.Errorf("读取 dimension count 失败: %w", err)
	}
	for i := int32(0); i < numDimensions; i++ {
		if _, err := readStringFromReader(r); err != nil {
			return fmt.Errorf("读取 dimension name 失败: %w", err)
		}
	}

	maxPlayers, err := readVarIntFromReader(r)
	if err != nil {
		return fmt.Errorf("读取 max_players 失败: %w", err)
	}
	_ = maxPlayers

	viewDistance, err := readVarIntFromReader(r)
	if err != nil {
		return fmt.Errorf("读取 view_distance 失败: %w", err)
	}
	_ = viewDistance

	simulationDistance, err := readVarIntFromReader(r)
	if err != nil {
		return fmt.Errorf("读取 simulation_distance 失败: %w", err)
	}
	_ = simulationDistance

	reducedDebugInfo, err := readBoolFromReader(r)
	if err != nil {
		return fmt.Errorf("读取 reduced_debug_info 失败: %w", err)
	}
	_ = reducedDebugInfo

	enableRespawnScreen, err := readBoolFromReader(r)
	if err != nil {
		return fmt.Errorf("读取 enable_respawn_screen 失败: %w", err)
	}
	_ = enableRespawnScreen

	doLimitedCrafting, err := readBoolFromReader(r)
	if err != nil {
		return fmt.Errorf("读取 do_limited_crafting 失败: %w", err)
	}
	_ = doLimitedCrafting

	dimensionType, err := readVarIntFromReader(r)
	if err != nil {
		return fmt.Errorf("读取 dimension_type 失败: %w", err)
	}
	_ = dimensionType

	dimensionName, err := readStringFromReader(r)
	if err != nil {
		return fmt.Errorf("读取 dimension_name 失败: %w", err)
	}
	c.Player.SetDimension(dimensionName)

	var hashedSeed int64
	binary.Read(r, binary.BigEndian, &hashedSeed)
	_ = hashedSeed

	gameMode, err := readU8(r)
	if err != nil {
		return fmt.Errorf("读取 game_mode 失败: %w", err)
	}
	c.Player.SetGameMode(player.GameMode(int(gameMode)))

	prevGameMode, _ := readU8(r)
	_ = prevGameMode

	isDebug, _ := readBoolFromReader(r)
	_ = isDebug
	isFlat, _ := readBoolFromReader(r)
	_ = isFlat

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
		_ = deathPos
	}

	portalCooldown, _ := readVarIntFromReader(r)
	_ = portalCooldown
	seaLevel, _ := readVarIntFromReader(r)
	_ = seaLevel
	secureChatEnforced, _ := readBoolFromReader(r)
	_ = secureChatEnforced

	logx.Infof("登录Play阶段: EntityID=%d, 维度=%s, 游戏模式=%s", entityID, dimensionName, c.Player.GameMode.String())

	return nil
}

func readInt64FromReader(r *bytes.Reader) (int64, error) {
	var v int64
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

func (c *Client) handleEntityDataPacket(data []byte) error {
	r := bytes.NewReader(data)

	entityID, err := readVarIntFromReader(r)
	if err != nil {
		return nil
	}

	if entityID != c.Player.EntityID {
		return nil
	}

	for {
		index, err := readU8(r)
		if err != nil {
			return nil
		}
		if index == 0xFF {
			break
		}

		typeID, err := readVarIntFromReader(r)
		if err != nil {
			return nil
		}

		switch index {
		case 1:
			airTicks, _ := readVarIntFromReader(r)
			c.Player.UpdateAir(int32(airTicks))
		case 9:
			var health float32
			binary.Read(r, binary.BigEndian, &health)
			c.Player.UpdateEntityHealth(health)
		default:
			_ = skipMetadataValue(r, typeID)
		}
	}

	return nil
}

func readFloat64FromReader(r *bytes.Reader) (float64, error) {
	var v float64
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

func readInt32FromReader(r *bytes.Reader) (int32, error) {
	var v int32
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

func readVarIntFromReader(r *bytes.Reader) (int32, error) {
	var result int32
	var shift uint8
	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		result |= int32(b&0x7F) << shift
		shift += 7
		if b&0x80 == 0 {
			break
		}
		if shift >= 35 {
			return 0, fmt.Errorf("VarInt too big")
		}
	}
	return result, nil
}

func readBoolFromReader(r *bytes.Reader) (bool, error) {
	b, err := r.ReadByte()
	if err != nil {
		return false, err
	}
	return b != 0, nil
}

func readStringFromReader(r *bytes.Reader) (string, error) {
	n, err := readVarIntFromReader(r)
	if err != nil {
		return "", err
	}
	if n < 0 || n > 32767*4 {
		return "", fmt.Errorf("invalid string length %d", n)
	}
	buf := make([]byte, n)
	if _, err := r.Read(buf); err != nil {
		return "", err
	}
	return nbt.CESU8ToUTF8(buf), nil
}

func readBytes(r *bytes.Reader, n int) []byte {
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return buf[:i]
		}
		buf[i] = b
	}
	return buf
}

func skipMetadataValue(r *bytes.Reader, typeID int32) error {
	switch typeID {
	case 0:
		_, _ = r.ReadByte()
	case 1:
		_, _ = readVarIntFromReader(r)
	case 2:
		_ = binary.Read(r, binary.BigEndian, new(float32))
	case 3, 6, 7:
		_ = skipNBT(r)
	case 4:
		_, _ = readStringFromReader(r)
	case 5:
		_, _ = readInt32FromReader(r)
		_, _ = readInt32FromReader(r)
		_, _ = readInt32FromReader(r)
	case 8, 9:
		_, _ = readBoolFromReader(r)
	case 10:
		_ = binary.Read(r, binary.BigEndian, new(float32))
	case 11:
		_, _ = readInt32FromReader(r)
		_, _ = readInt32FromReader(r)
		_, _ = readInt32FromReader(r)
	case 12:
		_, _ = readInt32FromReader(r)
		_, _ = readInt32FromReader(r)
		_, _ = readInt32FromReader(r)
		_, _ = readInt32FromReader(r)
		_, _ = readInt32FromReader(r)
		_, _ = readInt32FromReader(r)
	case 13:
		_, _ = readInt32FromReader(r)
		_, _ = readInt32FromReader(r)
		_, _ = readInt32FromReader(r)
		_ = skipNBT(r)
	case 14:
		_ = skipPrefixedArray(r, skipSlotData)
	case 15:
		_, _ = readInt32FromReader(r)
		_, _ = readInt32FromReader(r)
		_, _ = readInt32FromReader(r)
		_ = skipNBT(r)
	case 16:
		_, _ = readInt32FromReader(r)
		_, _ = readVarIntFromReader(r)
		_ = skipNBT(r)
	case 17:
		_ = binary.Read(r, binary.BigEndian, new(float32))
		_ = binary.Read(r, binary.BigEndian, new(float32))
		_ = binary.Read(r, binary.BigEndian, new(float32))
	case 18:
		_, _ = readInt32FromReader(r)
	case 19:
		_, _ = readInt32FromReader(r)
		_ = skipNBT(r)
	case 20:
		if has, _ := readBoolFromReader(r); has {
			if has2, _ := readBoolFromReader(r); has2 {
				_, _ = readStringFromReader(r)
				_, _ = readStringFromReader(r)
			}
		}
	case 21:
		_, _ = readInt32FromReader(r)
		_, _ = readVarIntFromReader(r)
		_ = skipNBT(r)
	case 22:
		_, _ = readInt32FromReader(r)
		_, _ = readInt32FromReader(r)
		_, _ = readInt32FromReader(r)
	}
	return nil
}

func skipNBT(r *bytes.Reader) error {
	dec := nbt.NewDecoder(r).NetworkFormat(true)
	return dec.Skip()
}

func skipPrefixedArray(r *bytes.Reader, fn func(*bytes.Reader) error) error {
	count, err := readVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < count; i++ {
		if err := fn(r); err != nil {
			return err
		}
	}
	return nil
}

func skipSlotData(r *bytes.Reader) error {
	count, err := readVarIntFromReader(r)
	if err != nil {
		return err
	}
	if count <= 0 {
		return nil
	}
	_, _ = readVarIntFromReader(r)
	numComponentsToAdd, _ := readVarIntFromReader(r)
	for i := int32(0); i < numComponentsToAdd; i++ {
		if err := skipComponentData(r); err != nil {
			return nil
		}
	}
	numComponentsToRemove, _ := readVarIntFromReader(r)
	for i := int32(0); i < numComponentsToRemove; i++ {
		_, _ = readVarIntFromReader(r)
	}
	return nil
}

func skipComponentData(r *bytes.Reader) error {
	componentType, err := readVarIntFromReader(r)
	if err != nil {
		return err
	}
	return skipComponentByType(r, componentType)
}

func skipComponentByType(r *bytes.Reader, componentType int32) error {
	switch componentType {
	case 0:
		return skipNBT(r)
	case 1, 2, 3:
		_, _ = readVarIntFromReader(r)
	case 4:
	case 5, 6:
		return skipNBT(r)
	case 7:
		_, _ = readStringFromReader(r)
	case 8:
		return skipPrefixedArray(r, skipNBT)
	case 9:
		_, _ = readVarIntFromReader(r)
	case 10:
		return skipEnchantments(r)
	case 11, 12:
		return skipBlockPredicates(r)
	case 13:
		return skipAttributeModifiers(r)
	case 14:
		return skipCustomModelData(r)
	case 15:
		_, _ = readBoolFromReader(r)
		return skipPrefixedArray(r, func(r *bytes.Reader) error { _, _ = readVarIntFromReader(r); return nil })
	case 16:
		_, _ = readVarIntFromReader(r)
	case 17:
	case 18:
		_, _ = readBoolFromReader(r)
	case 19:
		return skipNBT(r)
	case 20:
		_, _ = readVarIntFromReader(r)
		_ = binary.Read(r, binary.BigEndian, new(float32))
		_, _ = readBoolFromReader(r)
	case 21:
		return skipConsumable(r)
	case 22:
		_ = skipSlotData(r)
	case 23:
		_ = binary.Read(r, binary.BigEndian, new(float32))
		return skipPrefixedOptional(r, func(r *bytes.Reader) error { _, _ = readStringFromReader(r); return nil })
	case 24:
		_, _ = readStringFromReader(r)
	case 25:
		return skipTool(r)
	case 26:
		_, _ = readVarIntFromReader(r)
		_ = binary.Read(r, binary.BigEndian, new(float32))
	case 27:
		_, _ = readVarIntFromReader(r)
	case 28:
		return skipEquippable(r)
	case 29:
		return skipIDSet(r)
	case 30:
	case 31:
		_, _ = readStringFromReader(r)
	case 32:
		return skipPrefixedArray(r, skipConsumeEffect)
	case 33:
		return skipBlocksAttacks(r)
	case 34:
		return skipEnchantments(r)
	case 35:
		_, _ = readBoolFromReader(r)
		_, _ = readInt32FromReader(r)
	case 36:
		_, _ = readBoolFromReader(r)
		_, _ = readInt32FromReader(r)
	case 37:
		_, _ = readVarIntFromReader(r)
	case 38:
		return skipNBT(r)
	case 39:
		_, _ = readVarIntFromReader(r)
	case 40, 41:
		return skipPrefixedArray(r, skipSlotData)
	case 42:
		return skipPotionContents(r)
	case 43:
		_ = binary.Read(r, binary.BigEndian, new(float32))
	case 44:
		return skipPrefixedArray(r, skipPotionEffect)
	case 45:
		return skipWritableBookContent(r)
	case 46:
		return skipWrittenBookContent(r)
	case 47:
		return skipTrim(r)
	case 48:
		return skipNBT(r)
	case 49, 50:
		_, _ = readVarIntFromReader(r)
		return skipNBT(r)
	case 51:
		return skipIDOrX(r, skipInstrument)
	case 52:
		_, _ = readStringFromReader(r)
		_, _ = readBoolFromReader(r)
		_ = binary.Read(r, binary.BigEndian, new(float32))
		if v, _ := readBoolFromReader(r); v {
			_ = binary.Read(r, binary.BigEndian, new(float32))
		}
	case 53:
		_, _ = readStringFromReader(r)
	case 54:
		_, _ = readVarIntFromReader(r)
	case 55:
		return skipJukeboxPlayable(r)
	}
	return nil
}

type SlotData struct {
	ID    int32
	Count int32
}

func readSlotData(r *bytes.Reader) (*SlotData, error) {
	count, err := readVarIntFromReader(r)
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		return nil, nil
	}

	itemID, err := readVarIntFromReader(r)
	if err != nil {
		return nil, nil
	}

	numComponentsToAdd, err := readVarIntFromReader(r)
	if err != nil {
		return nil, nil
	}

	for i := int32(0); i < numComponentsToAdd; i++ {
		if err := skipComponentData(r); err != nil {
			return nil, nil
		}
	}

	numComponentsToRemove, err := readVarIntFromReader(r)
	if err != nil {
		return nil, nil
	}

	for i := int32(0); i < numComponentsToRemove; i++ {
		if _, err := readVarIntFromReader(r); err != nil {
			return nil, nil
		}
	}

	return &SlotData{ID: itemID, Count: count}, nil
}

func skipPrefixedOptional(r *bytes.Reader, fn func(*bytes.Reader) error) error {
	present, err := readBoolFromReader(r)
	if err != nil {
		return err
	}
	if present {
		return fn(r)
	}
	return nil
}

func skipEnchantments(r *bytes.Reader) error {
	return skipPrefixedArray(r, func(r *bytes.Reader) error {
		_, _ = readVarIntFromReader(r)
		_, _ = readVarIntFromReader(r)
		return nil
	})
}

func skipBlockPredicates(r *bytes.Reader) error {
	return skipPrefixedArray(r, skipBlockPredicate)
}

func skipBlockPredicate(r *bytes.Reader) error {
	holder, _ := readStringFromReader(r)
	if holder != "" {
		return nil
	}
	has, _ := readBoolFromReader(r)
	if has {
		_, _ = readStringFromReader(r)
	}
	num, _ := readVarIntFromReader(r)
	for i := int32(0); i < num; i++ {
		_, _ = readStringFromReader(r)
		_, _ = readVarIntFromReader(r)
	}
	num2, _ := readVarIntFromReader(r)
	for i := int32(0); i < num2; i++ {
		_, _ = readStringFromReader(r)
		_, _ = readStringFromReader(r)
	}
	return nil
}

func skipAttributeModifiers(r *bytes.Reader) error {
	return skipPrefixedArray(r, func(r *bytes.Reader) error {
		_, _ = readVarIntFromReader(r)
		_, _ = readStringFromReader(r)
		_, _ = readFloat64FromReader(r)
		_, _ = readVarIntFromReader(r)
		_, _ = readVarIntFromReader(r)
		return nil
	})
}

func skipCustomModelData(r *bytes.Reader) error {
	_ = skipPrefixedArray(r, func(r *bytes.Reader) error { _ = binary.Read(r, binary.BigEndian, new(float32)); return nil })
	_ = skipPrefixedArray(r, func(r *bytes.Reader) error { _, _ = readBoolFromReader(r); return nil })
	_ = skipPrefixedArray(r, func(r *bytes.Reader) error { _, _ = readStringFromReader(r); return nil })
	return skipPrefixedArray(r, func(r *bytes.Reader) error { _, _ = readInt32FromReader(r); return nil })
}

func skipConsumable(r *bytes.Reader) error {
	_ = binary.Read(r, binary.BigEndian, new(float32))
	_, _ = readVarIntFromReader(r)
	_ = skipIDOrX(r, skipSoundEvent)
	_, _ = readBoolFromReader(r)
	return skipPrefixedArray(r, skipConsumeEffect)
}

func skipTool(r *bytes.Reader) error {
	_ = skipPrefixedArray(r, func(r *bytes.Reader) error {
		_ = skipIDSet(r)
		_ = skipPrefixedOptional(r, func(r *bytes.Reader) error { _ = binary.Read(r, binary.BigEndian, new(float32)); return nil })
		_ = skipPrefixedOptional(r, func(r *bytes.Reader) error { _, _ = readBoolFromReader(r); return nil })
		return nil
	})
	_ = binary.Read(r, binary.BigEndian, new(float32))
	_, _ = readVarIntFromReader(r)
	_, _ = readBoolFromReader(r)
	return nil
}

func skipEquippable(r *bytes.Reader) error {
	_, _ = readVarIntFromReader(r)
	_ = skipIDOrX(r, skipSoundEvent)
	_ = skipPrefixedOptional(r, func(r *bytes.Reader) error { _, _ = readStringFromReader(r); return nil })
	_ = skipPrefixedOptional(r, func(r *bytes.Reader) error { _, _ = readStringFromReader(r); return nil })
	_ = skipPrefixedOptional(r, skipIDSet)
	_, _ = readBoolFromReader(r)
	_, _ = readBoolFromReader(r)
	_, _ = readBoolFromReader(r)
	return nil
}

func skipIDSet(r *bytes.Reader) error {
	setType, err := readVarIntFromReader(r)
	if err != nil {
		return err
	}
	if setType == 0 {
		_, _ = readStringFromReader(r)
	} else {
		for i := int32(0); i < setType-1; i++ {
			_, _ = readVarIntFromReader(r)
		}
	}
	return nil
}

func skipConsumeEffect(r *bytes.Reader) error {
	effectType, _ := readVarIntFromReader(r)
	switch effectType {
	case 0:
		_, _ = readVarIntFromReader(r)
		_ = binary.Read(r, binary.BigEndian, new(float32))
		_, _ = readVarIntFromReader(r)
	case 1:
		_ = skipPrefixedArray(r, skipConsumeEffect)
	case 2:
		_, _ = readVarIntFromReader(r)
	}
	return nil
}

func skipPotionContents(r *bytes.Reader) error {
	_ = skipPrefixedOptional(r, func(r *bytes.Reader) error { _, _ = readVarIntFromReader(r); return nil })
	_ = skipPrefixedOptional(r, func(r *bytes.Reader) error { _, _ = readInt32FromReader(r); return nil })
	_ = skipPrefixedArray(r, skipPotionEffect)
	_, _ = readStringFromReader(r)
	return nil
}

func skipPotionEffect(r *bytes.Reader) error {
	_, _ = readVarIntFromReader(r)
	_, _ = readVarIntFromReader(r)
	_, _ = readVarIntFromReader(r)
	return nil
}

func skipWritableBookContent(r *bytes.Reader) error {
	num, _ := readVarIntFromReader(r)
	for i := int32(0); i < num && i < 100; i++ {
		_, _ = readStringFromReader(r)
		_ = skipPrefixedOptional(r, func(r *bytes.Reader) error { _, _ = readStringFromReader(r); return nil })
	}
	return nil
}

func skipWrittenBookContent(r *bytes.Reader) error {
	_, _ = readStringFromReader(r)
	_ = skipPrefixedOptional(r, func(r *bytes.Reader) error { _, _ = readStringFromReader(r); return nil })
	_, _ = readStringFromReader(r)
	_, _ = readVarIntFromReader(r)
	num, _ := readVarIntFromReader(r)
	for i := int32(0); i < num && i < 100; i++ {
		_ = skipNBT(r)
		_ = skipPrefixedOptional(r, skipNBT)
	}
	_, _ = readBoolFromReader(r)
	return nil
}

func skipTrim(r *bytes.Reader) error {
	_ = skipIDOrX(r, skipTrimMaterial)
	return skipIDOrX(r, skipTrimPattern)
}

func skipIDOrX(r *bytes.Reader, skipInline func(*bytes.Reader) error) error {
	id, err := readVarIntFromReader(r)
	if err != nil {
		return err
	}
	if id == 0 {
		return skipInline(r)
	}
	return nil
}

func skipTrimMaterial(r *bytes.Reader) error {
	_, _ = readStringFromReader(r)
	_, _ = readStringFromReader(r)
	_ = skipPrefixedOptional(r, func(r *bytes.Reader) error { _, _ = readStringFromReader(r); return nil })
	return nil
}

func skipTrimPattern(r *bytes.Reader) error {
	_ = skipPrefixedOptional(r, func(r *bytes.Reader) error { _, _ = readStringFromReader(r); return nil })
	_, _ = readStringFromReader(r)
	return nil
}

func skipInstrument(r *bytes.Reader) error {
	_, _ = readStringFromReader(r)
	_, _ = readBoolFromReader(r)
	if v, _ := readBoolFromReader(r); v {
		_ = binary.Read(r, binary.BigEndian, new(float32))
	}
	return nil
}

func skipJukeboxPlayable(r *bytes.Reader) error {
	mode, _ := readU8(r)
	if mode == 0 {
		_, _ = readStringFromReader(r)
	} else {
		_, _ = readVarIntFromReader(r)
	}
	return nil
}

func skipSoundEvent(r *bytes.Reader) error {
	_, _ = readStringFromReader(r)
	if v, _ := readBoolFromReader(r); v {
		_ = binary.Read(r, binary.BigEndian, new(float32))
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
	if numItems > 100 {
		numItems = 100
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

func skipBlocksAttacks(r *bytes.Reader) error {
	_ = binary.Read(r, binary.BigEndian, new(float32))
	_ = binary.Read(r, binary.BigEndian, new(float32))
	_ = skipPrefixedArray(r, func(r *bytes.Reader) error {
		_ = binary.Read(r, binary.BigEndian, new(float32))
		_ = skipPrefixedOptional(r, skipIDSet)
		_ = binary.Read(r, binary.BigEndian, new(float32))
		_ = binary.Read(r, binary.BigEndian, new(float32))
		_ = binary.Read(r, binary.BigEndian, new(float32))
		_ = binary.Read(r, binary.BigEndian, new(float32))
		return nil
	})
	_ = skipPrefixedOptional(r, func(r *bytes.Reader) error { _, _ = readStringFromReader(r); return nil })
	_ = skipPrefixedOptional(r, skipIDOrXSoundEvent)
	_ = skipPrefixedOptional(r, skipIDOrXSoundEvent)
	return nil
}

func skipIDOrXSoundEvent(r *bytes.Reader) error {
	return skipIDOrX(r, skipSoundEvent)
}
