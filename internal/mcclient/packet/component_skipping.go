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
		0:   SkipNBT,                                                                   // custom_data
		1:   SkipVarInt,                                                                // max_stack_size
		2:   SkipVarInt,                                                                // max_damage
		3:   SkipVarInt,                                                                // damage
		4:   SkipNothing,                                                               // unbreakable
		5:   SkipUseEffects,                                                            // use_effects (1.21.11+)
		6:   SkipNBT,                                                                   // custom_name (text component)
		7:   SkipVarInt,                                                                // minimum_attack_charge
		8:   SkipRegistryValue,                                                         // damage_type
		9:   SkipNBT,                                                                   // item_name (text component)
		10:  SkipString,                                                                // item_model
		11:  func(r *bytes.Reader) error { return SkipPrefixedArray(r, SkipNBT) },      // lore
		12:  SkipVarInt,                                                                // rarity
		13:  SkipEnchantments,                                                          // enchantments
		14:  SkipBlockPredicates,                                                       // can_place_on
		15:  SkipBlockPredicates,                                                       // can_break
		16:  SkipAttributeModifiers,                                                    // attribute_modifiers
		17:  SkipCustomModelData,                                                       // custom_model_data
		18:  SkipTooltipDisplay,                                                        // tooltip_display
		19:  SkipVarInt,                                                                // repair_cost
		20:  SkipNothing,                                                               // creative_slot_lock
		21:  SkipBool,                                                                  // enchantment_glint_override
		22:  SkipNBT,                                                                   // intangible_projectile
		23:  SkipFood,                                                                  // food
		24:  SkipConsumable,                                                            // consumable
		25:  func(r *bytes.Reader) error { return SkipSlotData(r) },                    // use_remainder
		26:  SkipUseCooldown,                                                           // use_cooldown
		27:  SkipRegistryValue,                                                         // damage_resistant
		28:  SkipTool,                                                                  // tool
		29:  SkipWeapon,                                                                // weapon
		30:  SkipAttackRange,                                                           // attack_range (1.21.11+)
		31:  SkipVarInt,                                                                // enchantable
		32:  SkipEquippable,                                                            // equippable
		33:  SkipRepairable,                                                            // repairable
		34:  SkipNothing,                                                               // glider
		35:  SkipNBT,                                                                   // tooltip_style (text component)
		36:  SkipDeathProtection,                                                       // death_protection
		37:  SkipBlocksAttacks,                                                         // blocks_attacks
		38:  SkipPiercingWeapon,                                                        // piercing_weapon (1.21.11+)
		39:  SkipKineticWeapon,                                                         // kinetic_weapon (1.21.11+)
		40:  SkipSwingAnimation,                                                        // swing_animation (1.21.11+)
		41:  SkipEnchantments,                                                          // stored_enchantments
		42:  SkipInt32,                                                                 // dyed_color
		43:  SkipInt32,                                                                 // map_color
		44:  SkipVarInt,                                                                // map_id
		45:  SkipNBT,                                                                   // map_decorations
		46:  SkipVarInt,                                                                // map_post_processing
		47:  SkipFloat32,                                                               // potion_duration_scale
		48:  func(r *bytes.Reader) error { return SkipPrefixedArray(r, SkipSlotData) }, // charged_projectiles
		49:  func(r *bytes.Reader) error { return SkipPrefixedArray(r, SkipSlotData) }, // bundle_contents
		50:  SkipPotionContents,                                                        // potion_contents
		51:  SkipSuspiciousStewEffects,                                                 // suspicious_stew_effects
		52:  SkipWritableBookContent,                                                   // writable_book_content
		53:  SkipWrittenBookContent,                                                    // written_book_content
		54:  SkipTrim,                                                                  // trim
		55:  SkipNBT,                                                                   // debug_stick_state
		56:  SkipEntityData,                                                            // entity_data
		57:  SkipNBT,                                                                   // bucket_entity_data
		58:  SkipBlockEntityData,                                                       // block_entity_data
		59:  SkipInstrument,                                                            // instrument
		60:  SkipProvidesTrimMaterial,                                                  // provides_trim_material
		61:  SkipVarInt,                                                                // ominous_bottle_amplifier
		62:  SkipJukeboxPlayable,                                                       // jukebox_playable
		63:  SkipString,                                                                // provides_banner_patterns
		64:  SkipNBT,                                                                   // recipes
		65:  SkipLodestoneTracker,                                                      // lodestone_tracker
		66:  SkipFireworkExplosion,                                                     // firework_explosion
		67:  SkipFireworks,                                                             // fireworks
		68:  SkipProfile,                                                               // profile
		69:  SkipString,                                                                // note_block_sound
		70:  SkipBannerPatterns,                                                        // banner_patterns
		71:  SkipDyeColor,                                                              // base_color
		72:  SkipPotDecorations,                                                        // pot_decorations (1.21.11+)
		73:  SkipContainer,                                                             // container (1.21.11+)
		74:  SkipBlockState,                                                            // block_state
		75:  SkipBees,                                                                  // bees
		76:  SkipNBT,                                                                   // lock
		77:  SkipNBT,                                                                   // container_loot
		78:  SkipRegistryValue,                                                         // break_sound (1.21.11+)
		79:  SkipRegistryValue,                                                         // villager_variant
		80:  SkipRegistryValue,                                                         // wolf_variant
		81:  SkipRegistryValue,                                                         // cat_variant
		82:  SkipRegistryValue,                                                         // frog_variant
		83:  SkipRegistryValue,                                                         // axolotl_variant
		84:  SkipRegistryValue,                                                         // paintion_variant
		85:  SkipRegistryValue,                                                         // shulker_variant
		86:  SkipRegistryValue,                                                         // goat_variant
		87:  SkipRegistryValue,                                                         // sniffer_variant
		88:  SkipRegistryValue,                                                         // ghoul_variant
		89:  SkipRegistryValue,                                                         // breeze_variant
		90:  SkipRegistryValue,                                                         // bogged_variant
		91:  SkipVarInt,                                                                // bundle_remaining_space
		92:  SkipRegistryValue,                                                         // entity_color
		93:  SkipNBT,                                                                   // buckable
		94:  SkipRegistryValue,                                                         // armor_trim
		95:  SkipNBT,                                                                   // equippable_color
		96:  SkipNBT,                                                                   // trim_material
		97:  SkipNBT,                                                                   // trim_pattern
		98:  SkipNBT,                                                                   // compass_color
		99:  SkipNBT,                                                                   // map_color
		100: SkipVarInt,                                                                // frame_type
		101: SkipRegistryValue,                                                         // banner_pattern
		102: SkipVarInt,                                                                // base_color
		103: SkipVarInt,                                                                // color
	}
}

func SkipComponentByType(r *bytes.Reader, componentType int32) error {
	if skipper, ok := componentSkippers[componentType]; ok {
		err := skipper(r)
		if err != nil && err.Error() == "unexpected EOF" {
			return nil
		}
		return err
	}
	// 未知组件：尝试跳过为NBT（常见情况）
	logx.Warnf("未知组件类型: %d, 尝试NBT跳过", componentType)
	if r.Len() == 0 {
		return nil
	}
	err := SkipNBT(r)
	if err != nil && err.Error() == "unexpected EOF" {
		return nil
	}
	return err
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
		if r.Len() == 0 {
			return nil
		}
		if err := fn(r); err != nil {
			if err.Error() == "unexpected EOF" {
				return nil
			}
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

func SkipUseEffects(r *bytes.Reader) error {
	if _, err := ReadBoolFromReader(r); err != nil {
		return err
	}
	if _, err := ReadBoolFromReader(r); err != nil {
		return err
	}
	_, err := ReadFloat32FromReader(r)
	return err
}

func SkipRegistryValue(r *bytes.Reader) error {
	_, err := ReadVarIntFromReader(r)
	return err
}

func SkipAttackRange(r *bytes.Reader) error {
	for i := 0; i < 6; i++ {
		if _, err := ReadFloat32FromReader(r); err != nil {
			return err
		}
	}
	return nil
}

func SkipPiercingWeapon(r *bytes.Reader) error {
	if _, err := ReadBoolFromReader(r); err != nil {
		return err
	}
	if _, err := ReadBoolFromReader(r); err != nil {
		return err
	}
	if err := SkipPrefixedOptional(r, SkipRegistryValue); err != nil {
		return err
	}
	if err := SkipPrefixedOptional(r, SkipRegistryValue); err != nil {
		return err
	}
	return nil
}

func SkipKineticWeapon(r *bytes.Reader) error {
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	if err := SkipPrefixedOptional(r, func(r *bytes.Reader) error {
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
		if _, err := ReadFloat32FromReader(r); err != nil {
			return err
		}
		_, err := ReadVarIntFromReader(r)
		return err
	}); err != nil {
		return err
	}
	return nil
}

func SkipPotDecorations(r *bytes.Reader) error {
	for i := 0; i < 4; i++ {
		if err := SkipPrefixedOptional(r, SkipRegistryValue); err != nil {
			return err
		}
	}
	return nil
}

func SkipContainer(r *bytes.Reader) error {
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	length, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < length; i++ {
		if err := SkipSlotData(r); err != nil {
			return err
		}
	}
	return nil
}
