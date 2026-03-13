package packet

import (
	"bytes"
	"encoding/binary"
)

type componentSkipper func(*bytes.Reader) error

var componentSkippers map[int32]componentSkipper

func init() {
	componentSkippers = map[int32]componentSkipper{
		0:  SkipNBT,
		1:  SkipVarInt,
		2:  SkipVarInt,
		3:  SkipVarInt,
		4:  SkipNothing,
		5:  SkipNBT,
		6:  SkipNBT,
		7:  SkipString,
		8:  func(r *bytes.Reader) error { return SkipPrefixedArray(r, SkipNBT) },
		9:  SkipVarInt,
		10: SkipEnchantments,
		11: SkipBlockPredicates,
		12: SkipBlockPredicates,
		13: SkipAttributeModifiers,
		14: SkipCustomModelData,
		15: func(r *bytes.Reader) error {
			if _, err := ReadBoolFromReader(r); err != nil {
				return err
			}
			return SkipPrefixedArray(r, SkipVarInt)
		},
		16: SkipVarInt,
		17: SkipNothing,
		18: func(r *bytes.Reader) error { _, err := ReadBoolFromReader(r); return err },
		19: SkipNBT,
		20: func(r *bytes.Reader) error {
			if _, err := ReadVarIntFromReader(r); err != nil {
				return err
			}
			if err := SkipFloat32(r); err != nil {
				return err
			}
			_, err := ReadBoolFromReader(r)
			return err
		},
		21: SkipConsumable,
		22: func(r *bytes.Reader) error { return SkipSlotData(r) },
		23: func(r *bytes.Reader) error {
			if err := SkipFloat32(r); err != nil {
				return err
			}
			return SkipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := ReadStringFromReader(r); return err })
		},
		24: SkipString,
		25: SkipTool,
		26: func(r *bytes.Reader) error {
			if _, err := ReadVarIntFromReader(r); err != nil {
				return err
			}
			return SkipFloat32(r)
		},
		27: SkipVarInt,
		28: SkipEquippable,
		29: SkipIDSet,
		30: SkipNothing,
		31: SkipString,
		32: func(r *bytes.Reader) error { return SkipPrefixedArray(r, SkipConsumeEffect) },
		33: SkipBlocksAttacks,
		34: SkipEnchantments,
		35: func(r *bytes.Reader) error {
			if _, err := ReadBoolFromReader(r); err != nil {
				return err
			}
			return SkipInt32(r)
		},
		36: func(r *bytes.Reader) error {
			if _, err := ReadBoolFromReader(r); err != nil {
				return err
			}
			return SkipInt32(r)
		},
		37: SkipVarInt,
		38: SkipNBT,
		39: SkipVarInt,
		40: func(r *bytes.Reader) error { return SkipPrefixedArray(r, SkipSlotData) },
		41: func(r *bytes.Reader) error { return SkipPrefixedArray(r, SkipSlotData) },
		42: SkipPotionContents,
		43: SkipFloat32,
		44: func(r *bytes.Reader) error { return SkipPrefixedArray(r, SkipPotionEffect) },
		45: SkipWritableBookContent,
		46: SkipWrittenBookContent,
		47: SkipTrim,
		48: SkipNBT,
		49: func(r *bytes.Reader) error {
			if _, err := ReadVarIntFromReader(r); err != nil {
				return err
			}
			return SkipNBT(r)
		},
		50: func(r *bytes.Reader) error {
			if _, err := ReadVarIntFromReader(r); err != nil {
				return err
			}
			return SkipNBT(r)
		},
		51: func(r *bytes.Reader) error { return SkipIDOrX(r, SkipInstrument) },
		52: func(r *bytes.Reader) error {
			if _, err := ReadStringFromReader(r); err != nil {
				return err
			}
			if _, err := ReadBoolFromReader(r); err != nil {
				return err
			}
			if err := SkipFloat32(r); err != nil {
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
		},
		53: SkipString,
		54: SkipVarInt,
		55: SkipJukeboxPlayable,
	}
}

func SkipComponentByType(r *bytes.Reader, componentType int32) error {
	if skipper, ok := componentSkippers[componentType]; ok {
		return skipper(r)
	}
	return nil
}

// Helper skip functions

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
	return SkipPrefixedArray(r, SkipBlockPredicate)
}

func SkipBlockPredicate(r *bytes.Reader) error {
	holder, err := ReadStringFromReader(r)
	if err != nil {
		return err
	}
	if holder != "" {
		return nil
	}
	has, err := ReadBoolFromReader(r)
	if err != nil {
		return err
	}
	if has {
		if _, err := ReadStringFromReader(r); err != nil {
			return err
		}
	}
	num, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < num; i++ {
		if _, err := ReadStringFromReader(r); err != nil {
			return err
		}
	}
	num2, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < num2; i++ {
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
		if _, err := ReadFloat64FromReader(r); err != nil {
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
	if err := SkipPrefixedArray(r, func(r *bytes.Reader) error { return SkipFloat32(r) }); err != nil {
		return err
	}
	if err := SkipPrefixedArray(r, func(r *bytes.Reader) error { _, err := ReadBoolFromReader(r); return err }); err != nil {
		return err
	}
	if err := SkipPrefixedArray(r, func(r *bytes.Reader) error { _, err := ReadStringFromReader(r); return err }); err != nil {
		return err
	}
	return SkipPrefixedArray(r, func(r *bytes.Reader) error { _, err := ReadInt32FromReader(r); return err })
}

func SkipConsumable(r *bytes.Reader) error {
	if err := SkipFloat32(r); err != nil {
		return err
	}
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	if err := SkipIDOrX(r, SkipSoundEvent); err != nil {
		return err
	}
	if _, err := ReadBoolFromReader(r); err != nil {
		return err
	}
	return SkipPrefixedArray(r, SkipConsumeEffect)
}

func SkipTool(r *bytes.Reader) error {
	if err := SkipPrefixedArray(r, func(r *bytes.Reader) error {
		if err := SkipIDSet(r); err != nil {
			return err
		}
		if err := SkipPrefixedOptional(r, func(r *bytes.Reader) error { return SkipFloat32(r) }); err != nil {
			return err
		}
		return SkipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := ReadBoolFromReader(r); return err })
	}); err != nil {
		return err
	}
	if err := SkipFloat32(r); err != nil {
		return err
	}
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	_, err := ReadBoolFromReader(r)
	return err
}

func SkipEquippable(r *bytes.Reader) error {
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	if err := SkipIDOrX(r, SkipSoundEvent); err != nil {
		return err
	}
	if err := SkipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := ReadStringFromReader(r); return err }); err != nil {
		return err
	}
	if err := SkipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := ReadStringFromReader(r); return err }); err != nil {
		return err
	}
	if err := SkipPrefixedOptional(r, SkipIDSet); err != nil {
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

func SkipIDSet(r *bytes.Reader) error {
	setType, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	if setType == 0 {
		_, err := ReadStringFromReader(r)
		return err
	}
	for i := int32(0); i < setType-1; i++ {
		if _, err := ReadVarIntFromReader(r); err != nil {
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
		return SkipPrefixedArray(r, SkipConsumeEffect)
	case 2:
		_, err := ReadVarIntFromReader(r)
		return err
	}
	return nil
}

func SkipPotionContents(r *bytes.Reader) error {
	if err := SkipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := ReadVarIntFromReader(r); return err }); err != nil {
		return err
	}
	if err := SkipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := ReadInt32FromReader(r); return err }); err != nil {
		return err
	}
	if err := SkipPrefixedArray(r, SkipPotionEffect); err != nil {
		return err
	}
	_, err := ReadStringFromReader(r)
	return err
}

func SkipPotionEffect(r *bytes.Reader) error {
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	_, err := ReadVarIntFromReader(r)
	return err
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
		if err := SkipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := ReadStringFromReader(r); return err }); err != nil {
			return err
		}
	}
	return nil
}

func SkipWrittenBookContent(r *bytes.Reader) error {
	if _, err := ReadStringFromReader(r); err != nil {
		return err
	}
	if err := SkipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := ReadStringFromReader(r); return err }); err != nil {
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

func SkipTrim(r *bytes.Reader) error {
	if err := SkipIDOrX(r, SkipTrimMaterial); err != nil {
		return err
	}
	return SkipIDOrX(r, SkipTrimPattern)
}

func SkipTrimMaterial(r *bytes.Reader) error {
	if _, err := ReadStringFromReader(r); err != nil {
		return err
	}
	if _, err := ReadStringFromReader(r); err != nil {
		return err
	}
	return SkipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := ReadStringFromReader(r); return err })
}

func SkipTrimPattern(r *bytes.Reader) error {
	if err := SkipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := ReadStringFromReader(r); return err }); err != nil {
		return err
	}
	_, err := ReadStringFromReader(r)
	return err
}

func SkipInstrument(r *bytes.Reader) error {
	if _, err := ReadStringFromReader(r); err != nil {
		return err
	}
	if _, err := ReadBoolFromReader(r); err != nil {
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

func SkipJukeboxPlayable(r *bytes.Reader) error {
	mode, err := ReadU8(r)
	if err != nil {
		return err
	}
	if mode == 0 {
		_, err = ReadStringFromReader(r)
		return err
	}
	_, err = ReadVarIntFromReader(r)
	return err
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

func SkipSlotData(r *bytes.Reader) error {
	_, err := ReadSlotData(r)
	return err
}

func SkipBlocksAttacks(r *bytes.Reader) error {
	if err := SkipFloat32(r); err != nil {
		return err
	}
	if err := SkipFloat32(r); err != nil {
		return err
	}
	if err := SkipPrefixedArray(r, func(r *bytes.Reader) error {
		if err := SkipFloat32(r); err != nil {
			return err
		}
		if err := SkipPrefixedOptional(r, SkipIDSet); err != nil {
			return err
		}
		if err := SkipFloat32(r); err != nil {
			return err
		}
		if err := SkipFloat32(r); err != nil {
			return err
		}
		if err := SkipFloat32(r); err != nil {
			return err
		}
		return SkipFloat32(r)
	}); err != nil {
		return err
	}
	if err := SkipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := ReadStringFromReader(r); return err }); err != nil {
		return err
	}
	if err := SkipPrefixedOptional(r, func(r *bytes.Reader) error { return SkipIDOrX(r, SkipSoundEvent) }); err != nil {
		return err
	}
	return SkipPrefixedOptional(r, func(r *bytes.Reader) error { return SkipIDOrX(r, SkipSoundEvent) })
}
