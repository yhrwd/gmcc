package packet

import (
	"bytes"
	"encoding/binary"

	"gmcc/internal/logx"
)

type componentSkipper func(*bytes.Reader) error

var componentSkippers map[int32]componentSkipper

func init() {
	componentSkippers = map[int32]componentSkipper{
		0:  SkipNBT,                                                                   // custom_data
		1:  SkipVarInt,                                                                // max_stack_size
		2:  SkipVarInt,                                                                // max_damage
		3:  SkipVarInt,                                                                // damage
		4:  SkipNothing,                                                               // unbreakable
		5:  SkipNBT,                                                                   // custom_name (text component)
		6:  SkipNBT,                                                                   // item_name (text component)
		7:  SkipString,                                                                // item_model
		8:  func(r *bytes.Reader) error { return SkipPrefixedArray(r, SkipNBT) },      // lore
		9:  SkipVarInt,                                                                // rarity
		10: SkipEnchantments,                                                          // enchantments
		11: SkipBlockPredicates,                                                       // can_place_on
		12: SkipBlockPredicates,                                                       // can_break
		13: SkipAttributeModifiers,                                                    // attribute_modifiers
		14: SkipCustomModelData,                                                       // custom_model_data
		15: SkipTooltipDisplay,                                                        // tooltip_display
		16: SkipVarInt,                                                                // repair_cost
		17: SkipNothing,                                                               // creative_slot_lock
		18: SkipBool,                                                                  // enchantment_glint_override
		19: SkipNBT,                                                                   // intangible_projectile
		20: SkipFood,                                                                  // food
		21: SkipConsumable,                                                            // consumable
		22: func(r *bytes.Reader) error { return SkipSlotData(r) },                    // use_remainder
		23: SkipUseCooldown,                                                           // use_cooldown
		24: SkipString,                                                                // damage_resistant
		25: SkipTool,                                                                  // tool
		26: SkipWeapon,                                                                // weapon
		27: SkipVarInt,                                                                // enchantable
		28: SkipEquippable,                                                            // equippable
		29: SkipRepairable,                                                            // repairable
		30: SkipNothing,                                                               // glider
		31: SkipString,                                                                // tooltip_style
		32: SkipDeathProtection,                                                       // death_protection
		33: SkipBlocksAttacks,                                                         // blocks_attacks
		34: SkipEnchantments,                                                          // stored_enchantments
		35: SkipInt32,                                                                 // dyed_color
		36: SkipInt32,                                                                 // map_color
		37: SkipVarInt,                                                                // map_id
		38: SkipNBT,                                                                   // map_decorations
		39: SkipVarInt,                                                                // map_post_processing
		40: SkipFloat32,                                                               // potion_duration_scale
		41: func(r *bytes.Reader) error { return SkipPrefixedArray(r, SkipSlotData) }, // charged_projectiles
		42: func(r *bytes.Reader) error { return SkipPrefixedArray(r, SkipSlotData) }, // bundle_contents
		43: SkipPotionContents,                                                        // potion_contents
		44: SkipSuspiciousStewEffects,                                                 // suspicious_stew_effects
		45: SkipWritableBookContent,                                                   // writable_book_content
		46: SkipWrittenBookContent,                                                    // written_book_content
		47: SkipTrim,                                                                  // trim
		48: SkipNBT,                                                                   // debug_stick_state
		49: SkipEntityData,                                                            // entity_data
		50: SkipNBT,                                                                   // bucket_entity_data
		51: SkipBlockEntityData,                                                       // block_entity_data
		52: SkipInstrument,                                                            // instrument
		53: SkipProvidesTrimMaterial,                                                  // provides_trim_material
		54: SkipVarInt,                                                                // ominous_bottle_amplifier
		55: SkipJukeboxPlayable,                                                       // jukebox_playable
		56: SkipString,                                                                // provides_banner_patterns
		57: SkipNBT,                                                                   // recipes
		58: SkipLodestoneTracker,                                                      // lodestone_tracker
		59: SkipFireworkExplosion,                                                     // firework_explosion
		60: SkipFireworks,                                                             // fireworks
		61: SkipProfile,                                                               // profile
		62: SkipString,                                                                // note_block_sound
		63: SkipBannerPatterns,                                                        // banner_patterns
		64: SkipDyeColor,                                                              // base_color
		65: func(r *bytes.Reader) error { return SkipPrefixedArray(r, SkipVarInt) },   // pot_decorations
		66: func(r *bytes.Reader) error { return SkipPrefixedArray(r, SkipSlotData) }, // container
		67: SkipBlockState,                                                            // block_state
		68: SkipBees,                                                                  // bees
		69: SkipNBT,                                                                   // lock
		70: SkipNBT,                                                                   // container_loot
		71: SkipSoundEvent,                                                            // break_sound
		// 72+ 是实体变种子组件，根据版本可能不存在
	}
}

func SkipComponentByType(r *bytes.Reader, componentType int32) error {
	if skipper, ok := componentSkippers[componentType]; ok {
		return skipper(r)
	}
	// 未知组件：尝试跳过为NBT（常见情况）
	logx.Warnf("未知组件类型: %d, 尝试NBT跳过", componentType)
	return SkipNBT(r)
}

func SkipNothing(r *bytes.Reader) error {
	return nil
}

func SkipVarInt(r *bytes.Reader) error {
	_, err := ReadVarIntFromReader(r)
	return err
}

func SkipString(r *bytes.Reader) error {
	_, err := ReadStringFromReader(r)
	return err
}

func SkipInt32(r *bytes.Reader) error {
	_, err := ReadInt32FromReader(r)
	return err
}

func SkipFloat32(r *bytes.Reader) error {
	return binary.Read(r, binary.BigEndian, new(float32))
}

func SkipPrefixedArray(r *bytes.Reader, fn func(*bytes.Reader) error) error {
	length, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < length; i++ {
		if err := fn(r); err != nil {
			return err
		}
	}
	return nil
}

func SkipPrefixedOptional(r *bytes.Reader, fn func(*bytes.Reader) error) error {
	present, err := ReadBoolFromReader(r)
	if err != nil {
		return err
	}
	if present {
		return fn(r)
	}
	return nil
}

func SkipSoundEvent(r *bytes.Reader) error {
	if _, err := ReadStringFromReader(r); err != nil {
		return err
	}
	present, err := ReadBoolFromReader(r)
	if err != nil {
		return err
	}
	if present {
		return SkipFloat32(r)
	}
	return nil
}

func SkipBannerPatterns(r *bytes.Reader) error {
	return SkipPrefixedArray(r, func(r *bytes.Reader) error {
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
		_, err := ReadVarIntFromReader(r)
		return err
	})
}

func SkipBlockState(r *bytes.Reader) error {
	return SkipPrefixedArray(r, func(r *bytes.Reader) error {
		if _, err := ReadStringFromReader(r); err != nil {
			return err
		}
		_, err := ReadStringFromReader(r)
		return err
	})
}

func SkipFireworks(r *bytes.Reader) error {
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	return SkipPrefixedArray(r, SkipFireworkExplosion)
}

func SkipFireworkExplosion(r *bytes.Reader) error {
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	if err := SkipPrefixedArray(r, SkipInt32); err != nil {
		return err
	}
	if err := SkipPrefixedArray(r, SkipInt32); err != nil {
		return err
	}
	if _, err := ReadBoolFromReader(r); err != nil {
		return err
	}
	_, err := ReadBoolFromReader(r)
	return err
}

func SkipLodestoneTracker(r *bytes.Reader) error {
	if err := SkipPrefixedOptional(r, func(r *bytes.Reader) error {
		if _, err := ReadStringFromReader(r); err != nil {
			return err
		}
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
		_, err := ReadVarIntFromReader(r)
		return err
	}); err != nil {
		return err
	}
	_, err := ReadBoolFromReader(r)
	return err
}

func SkipProvidesTrimMaterial(r *bytes.Reader) error {
	id, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	if id == 0 {
		if _, err := ReadStringFromReader(r); err != nil {
			return err
		}
		if _, err := ReadStringFromReader(r); err != nil {
			return err
		}
		return SkipPrefixedOptional(r, SkipString)
	}
	return nil
}

func SkipRepairable(r *bytes.Reader) error {
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err := ReadStringFromReader(r)
		return err
	}
	for i := int32(0); i < count; i++ {
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
	}
	return nil
}

func SkipSuspiciousStewEffects(r *bytes.Reader) error {
	return SkipPrefixedArray(r, func(r *bytes.Reader) error {
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
		_, err := ReadVarIntFromReader(r)
		return err
	})
}

func SkipSwingAnimation(r *bytes.Reader) error {
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	_, err := ReadVarIntFromReader(r)
	return err
}

func SkipTooltipDisplay(r *bytes.Reader) error {
	if _, err := ReadBoolFromReader(r); err != nil {
		return err
	}
	return SkipPrefixedArray(r, SkipVarInt)
}

func SkipUseCooldown(r *bytes.Reader) error {
	if err := SkipFloat32(r); err != nil {
		return err
	}
	return SkipPrefixedOptional(r, SkipString)
}

func SkipWritableBookContent(r *bytes.Reader) error {
	num, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < num && i < 100; i++ {
		if _, err := ReadStringFromReader(r); err != nil {
			return err
		}
		if err := SkipPrefixedOptional(r, SkipString); err != nil {
			return err
		}
	}
	return nil
}

func SkipWrittenBookContent(r *bytes.Reader) error {
	if _, err := ReadStringFromReader(r); err != nil {
		return err
	}
	if err := SkipPrefixedOptional(r, SkipString); err != nil {
		return err
	}
	if _, err := ReadStringFromReader(r); err != nil {
		return err
	}
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	num, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < num && i < 100; i++ {
		if err := SkipNBT(r); err != nil {
			return err
		}
		if err := SkipPrefixedOptional(r, SkipNBT); err != nil {
			return err
		}
	}
	_, err = ReadBoolFromReader(r)
	return err
}

func SkipDeathProtection(r *bytes.Reader) error {
	return SkipPrefixedArray(r, func(r *bytes.Reader) error {
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
		effectType, err := ReadVarIntFromReader(r)
		if err != nil {
			return err
		}
		switch effectType {
		case 0:
			if _, err := ReadVarIntFromReader(r); err != nil {
				return err
			}
			if err := SkipFloat32(r); err != nil {
				return err
			}
			_, err := ReadVarIntFromReader(r)
			return err
		case 1:
			return SkipDeathProtection(r)
		case 2:
			_, err := ReadVarIntFromReader(r)
			return err
		}
		return nil
	})
}

func SkipSlotData(r *bytes.Reader) error {
	_, err := ReadSlotData(r)
	return err
}

func SkipIDOrX(r *bytes.Reader, skipInline func(*bytes.Reader) error) error {
	id, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	if id == 0 {
		return skipInline(r)
	}
	return nil
}

func SkipBool(r *bytes.Reader) error {
	_, err := ReadBoolFromReader(r)
	return err
}

func SkipBooleanAND(r *bytes.Reader) error {
	length, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < length; i++ {
		if err := SkipNBT(r); err != nil {
			return err
		}
	}
	return nil
}

func SkipEnchantments(r *bytes.Reader) error {
	return SkipPrefixedArray(r, func(r *bytes.Reader) error {
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
		_, err := ReadVarIntFromReader(r)
		return err
	})
}

func SkipBlockPredicates(r *bytes.Reader) error {
	length, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < length; i++ {
		if err := SkipBlockPredicate(r); err != nil {
			return err
		}
	}
	return nil
}

func SkipBlockPredicate(r *bytes.Reader) error {
	present, err := ReadBoolFromReader(r)
	if err != nil {
		return err
	}
	if present {
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
	}
	length, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < length; i++ {
		if _, err := ReadStringFromReader(r); err != nil {
			return err
		}
		if _, err := ReadStringFromReader(r); err != nil {
			return err
		}
	}
	return nil
}

func SkipAttributeModifiers(r *bytes.Reader) error {
	return SkipPrefixedArray(r, func(r *bytes.Reader) error {
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
		if _, err := ReadStringFromReader(r); err != nil {
			return err
		}
		if err := binary.Read(r, binary.BigEndian, new(float64)); err != nil {
			return err
		}
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
		_, err := ReadVarIntFromReader(r)
		return err
	})
}

func SkipCustomModelData(r *bytes.Reader) error {
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	return nil
}

func SkipFood(r *bytes.Reader) error {
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	if err := SkipFloat32(r); err != nil {
		return err
	}
	_, err := ReadBoolFromReader(r)
	return err
}

func SkipConsumable(r *bytes.Reader) error {
	if err := SkipFloat32(r); err != nil {
		return err
	}
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	if err := SkipSoundEvent(r); err != nil {
		return err
	}
	if _, err := ReadBoolFromReader(r); err != nil {
		return err
	}
	length, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < length; i++ {
		if err := SkipConsumeEffect(r); err != nil {
			return err
		}
	}
	return nil
}

func SkipConsumeEffect(r *bytes.Reader) error {
	effectType, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	switch effectType {
	case 0:
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
		if err := SkipFloat32(r); err != nil {
			return err
		}
		_, err := ReadVarIntFromReader(r)
		return err
	case 1:
		return SkipConsumeEffect(r)
	case 2:
		_, err := ReadVarIntFromReader(r)
		return err
	}
	return nil
}

func SkipTool(r *bytes.Reader) error {
	length, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < length; i++ {
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
	}
	return nil
}

func SkipWeapon(r *bytes.Reader) error {
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	return SkipFloat32(r)
}

func SkipEquippable(r *bytes.Reader) error {
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	if err := SkipSoundEvent(r); err != nil {
		return err
	}
	if err := SkipPrefixedOptional(r, SkipString); err != nil {
		return err
	}
	if err := SkipPrefixedOptional(r, SkipString); err != nil {
		return err
	}
	if err := SkipPrefixedOptional(r, SkipVarInt); err != nil {
		return err
	}
	if _, err := ReadBoolFromReader(r); err != nil {
		return err
	}
	if _, err := ReadBoolFromReader(r); err != nil {
		return err
	}
	if _, err := ReadBoolFromReader(r); err != nil {
		return err
	}
	_, err := ReadBoolFromReader(r)
	return err
}

func SkipBlocksAttacks(r *bytes.Reader) error {
	if err := SkipFloat32(r); err != nil {
		return err
	}
	if err := SkipFloat32(r); err != nil {
		return err
	}
	length, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < length; i++ {
		if err := SkipFloat32(r); err != nil {
			return err
		}
	}
	return nil
}

func SkipPotionContents(r *bytes.Reader) error {
	if err := SkipPrefixedOptional(r, SkipVarInt); err != nil {
		return err
	}
	if err := SkipPrefixedOptional(r, SkipInt32); err != nil {
		return err
	}
	length, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < length; i++ {
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
	}
	_, err = ReadStringFromReader(r)
	return err
}

func SkipTrim(r *bytes.Reader) error {
	if err := SkipIDOrX(r, func(r *bytes.Reader) error {
		return SkipNBT(r)
	}); err != nil {
		return err
	}
	return SkipIDOrX(r, func(r *bytes.Reader) error {
		return SkipNBT(r)
	})
}

func SkipEntityData(r *bytes.Reader) error {
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	return SkipNBT(r)
}

func SkipBlockEntityData(r *bytes.Reader) error {
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	return SkipNBT(r)
}

func SkipInstrument(r *bytes.Reader) error {
	return SkipIDOrX(r, SkipString)
}

func SkipJukeboxPlayable(r *bytes.Reader) error {
	if _, err := ReadU8(r); err != nil {
		return err
	}
	return SkipIDOrX(r, SkipString)
}

func SkipProfile(r *bytes.Reader) error {
	present, err := ReadBoolFromReader(r)
	if err != nil {
		return err
	}
	if present {
		if _, err := ReadStringFromReader(r); err != nil {
			return err
		}
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
	}
	return nil
}

func SkipDyeColor(r *bytes.Reader) error {
	_, err := ReadVarIntFromReader(r)
	return err
}

func SkipBees(r *bytes.Reader) error {
	length, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < length; i++ {
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
		if err := SkipNBT(r); err != nil {
			return err
		}
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
	}
	return nil
}
