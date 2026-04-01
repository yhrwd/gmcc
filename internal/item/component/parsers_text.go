package component

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gmcc/internal/mcclient/chat"
	"gmcc/internal/nbt"
)

// ParseTextComponent 解析文本组件 (NBT格式)
func ParseTextComponent(r *bytes.Reader) (*chat.TextComponent, error) {
	dec := nbt.NewDecoder(r)
	dec.NetworkFormat(true)

	var nbtData map[string]any
	if err := dec.Decode(&nbtData); err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal(nbtData)
	if err != nil {
		return nil, err
	}

	var tc chat.TextComponent
	if err := json.Unmarshal(jsonBytes, &tc); err != nil {
		return nil, err
	}

	return &tc, nil
}

// ParseTextComponentList 解析文本组件列表 (lore组件)
func ParseTextComponentList(r *bytes.Reader) ([]*chat.TextComponent, error) {
	dec := nbt.NewDecoder(r)
	dec.NetworkFormat(true)

	var list []any
	if err := dec.Decode(&list); err != nil {
		return nil, err
	}

	result := make([]*chat.TextComponent, 0, len(list))
	for _, item := range list {
		jsonBytes, err := json.Marshal(item)
		if err != nil {
			continue
		}

		var tc chat.TextComponent
		if err := json.Unmarshal(jsonBytes, &tc); err != nil {
			continue
		}

		result = append(result, &tc)
	}

	return result, nil
}

// ParseCustomName 解析 custom_name 组件 (ID: 6)
func ParseCustomName(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
	tc, err := ParseTextComponent(r)
	if err != nil {
		return nil, err
	}

	return &ComponentResult{
		TypeID: typeID,
		Data:   tc,
	}, nil
}

// ParseItemName 解析 item_name 组件 (ID: 9)
func ParseItemName(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
	tc, err := ParseTextComponent(r)
	if err != nil {
		return nil, err
	}

	return &ComponentResult{
		TypeID: typeID,
		Data:   tc,
	}, nil
}

// ParseLore 解析 lore 组件 (ID: 11)
func ParseLore(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
	list, err := ParseTextComponentList(r)
	if err != nil {
		return nil, err
	}

	return &ComponentResult{
		TypeID: typeID,
		Data:   list,
	}, nil
}

// ParseRarity 解析 rarity 组件 (ID: 12, VarInt枚举)
func ParseRarity(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
	value, err := readVarInt(r)
	if err != nil {
		return nil, err
	}

	rarityMap := map[int32]string{
		0: "common",
		1: "uncommon",
		2: "rare",
		3: "epic",
	}

	rarityName := rarityMap[value]
	if rarityName == "" {
		rarityName = "common"
	}

	return &ComponentResult{
		TypeID: typeID,
		Data:   rarityName,
	}, nil
}

// ParseVarInt 通用 VarInt 解析
func ParseVarInt(r *bytes.Reader) (*ComponentResult, error) {
	value, err := readVarInt(r)
	if err != nil {
		return nil, fmt.Errorf("读取 VarInt 失败: %w", err)
	}
	return &ComponentResult{Data: value}, nil
}

// EnchantmentEntry 附魔条目结构
type EnchantmentEntry struct {
	ID    string
	Level int32
}

// ParseEnchantments 解析 enchantments 组件 (ID: 13)
func ParseEnchantments(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
	dec := nbt.NewDecoder(r)
	dec.NetworkFormat(true)

	var enchantments map[string]int32
	if err := dec.Decode(&enchantments); err != nil {
		return nil, err
	}

	result := make([]EnchantmentEntry, 0, len(enchantments))
	for id, level := range enchantments {
		result = append(result, EnchantmentEntry{
			ID:    id,
			Level: level,
		})
	}

	return &ComponentResult{
		TypeID: typeID,
		Data:   result,
	}, nil
}

func readVarInt(r *bytes.Reader) (int32, error) {
	result := int32(0)
	shift := uint32(0)

	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		result |= int32(b&0x7F) << shift
		shift += 7

		if (b & 0x80) == 0 {
			break
		}
	}

	return result, nil
}
