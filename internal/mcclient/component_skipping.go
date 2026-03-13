package mcclient

import (
	"bytes"
	"encoding/binary"
)

type componentSkipper func(*bytes.Reader) error

var componentSkippers map[int32]componentSkipper

func init() {
	componentSkippers = map[int32]componentSkipper{
		0:  skipNBT,
		1:  skipVarInt,
		2:  skipVarInt,
		3:  skipVarInt,
		4:  skipNothing,
		5:  skipNBT,
		6:  skipNBT,
		7:  skipString,
		8:  func(r *bytes.Reader) error { return skipPrefixedArray(r, skipNBT) },
		9:  skipVarInt,
		10: skipEnchantments,
		11: skipBlockPredicates,
		12: skipBlockPredicates,
		13: skipAttributeModifiers,
		14: skipCustomModelData,
		15: func(r *bytes.Reader) error {
			if _, err := readBoolFromReader(r); err != nil {
				return err
			}
			return skipPrefixedArray(r, skipVarInt)
		},
		16: skipVarInt,
		17: skipNothing,
		18: func(r *bytes.Reader) error { _, err := readBoolFromReader(r); return err },
		19: skipNBT,
		20: func(r *bytes.Reader) error {
			if _, err := readVarIntFromReader(r); err != nil {
				return err
			}
			if err := skipFloat32(r); err != nil {
				return err
			}
			_, err := readBoolFromReader(r)
			return err
		},
		21: skipConsumable,
		22: func(r *bytes.Reader) error { return skipSlotData(r) },
		23: func(r *bytes.Reader) error {
			if err := skipFloat32(r); err != nil {
				return err
			}
			return skipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := readStringFromReader(r); return err })
		},
		24: skipString,
		25: skipTool,
		26: func(r *bytes.Reader) error {
			if _, err := readVarIntFromReader(r); err != nil {
				return err
			}
			return skipFloat32(r)
		},
		27: skipVarInt,
		28: skipEquippable,
		29: skipIDSet,
		30: skipNothing,
		31: skipString,
		32: func(r *bytes.Reader) error { return skipPrefixedArray(r, skipConsumeEffect) },
		33: skipBlocksAttacks,
		34: skipEnchantments,
		35: func(r *bytes.Reader) error {
			if _, err := readBoolFromReader(r); err != nil {
				return err
			}
			return skipInt32(r)
		},
		36: func(r *bytes.Reader) error {
			if _, err := readBoolFromReader(r); err != nil {
				return err
			}
			return skipInt32(r)
		},
		37: skipVarInt,
		38: skipNBT,
		39: skipVarInt,
		40: func(r *bytes.Reader) error { return skipPrefixedArray(r, skipSlotData) },
		41: func(r *bytes.Reader) error { return skipPrefixedArray(r, skipSlotData) },
		42: skipPotionContents,
		43: skipFloat32,
		44: func(r *bytes.Reader) error { return skipPrefixedArray(r, skipPotionEffect) },
		45: skipWritableBookContent,
		46: skipWrittenBookContent,
		47: skipTrim,
		48: skipNBT,
		49: func(r *bytes.Reader) error {
			if _, err := readVarIntFromReader(r); err != nil {
				return err
			}
			return skipNBT(r)
		},
		50: func(r *bytes.Reader) error {
			if _, err := readVarIntFromReader(r); err != nil {
				return err
			}
			return skipNBT(r)
		},
		51: func(r *bytes.Reader) error { return skipIDOrX(r, skipInstrument) },
		52: func(r *bytes.Reader) error {
			if _, err := readStringFromReader(r); err != nil {
				return err
			}
			if _, err := readBoolFromReader(r); err != nil {
				return err
			}
			if err := skipFloat32(r); err != nil {
				return err
			}
			present, err := readBoolFromReader(r)
			if err != nil {
				return err
			}
			if present {
				return skipFloat32(r)
			}
			return nil
		},
		53: skipString,
		54: skipVarInt,
		55: skipJukeboxPlayable,
	}
}

func skipComponentByType(r *bytes.Reader, componentType int32) error {
	if skipper, ok := componentSkippers[componentType]; ok {
		return skipper(r)
	}
	return nil
}

// Helper skip functions

func skipNothing(r *bytes.Reader) error {
	return nil
}

func skipVarInt(r *bytes.Reader) error {
	_, err := readVarIntFromReader(r)
	return err
}

func skipString(r *bytes.Reader) error {
	_, err := readStringFromReader(r)
	return err
}

func skipInt32(r *bytes.Reader) error {
	_, err := readInt32FromReader(r)
	return err
}

func skipFloat32(r *bytes.Reader) error {
	return binary.Read(r, binary.BigEndian, new(float32))
}

func skipPrefixedArray(r *bytes.Reader, fn func(*bytes.Reader) error) error {
	length, err := readVarIntFromReader(r)
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
		if _, err := readVarIntFromReader(r); err != nil {
			return err
		}
		_, err := readVarIntFromReader(r)
		return err
	})
}

func skipBlockPredicates(r *bytes.Reader) error {
	return skipPrefixedArray(r, skipBlockPredicate)
}

func skipBlockPredicate(r *bytes.Reader) error {
	holder, err := readStringFromReader(r)
	if err != nil {
		return err
	}
	if holder != "" {
		return nil
	}
	has, err := readBoolFromReader(r)
	if err != nil {
		return err
	}
	if has {
		if _, err := readStringFromReader(r); err != nil {
			return err
		}
	}
	num, err := readVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < num; i++ {
		if _, err := readStringFromReader(r); err != nil {
			return err
		}
	}
	num2, err := readVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < num2; i++ {
		if _, err := readStringFromReader(r); err != nil {
			return err
		}
		if _, err := readStringFromReader(r); err != nil {
			return err
		}
	}
	return nil
}

func skipAttributeModifiers(r *bytes.Reader) error {
	return skipPrefixedArray(r, func(r *bytes.Reader) error {
		if _, err := readVarIntFromReader(r); err != nil {
			return err
		}
		if _, err := readStringFromReader(r); err != nil {
			return err
		}
		if _, err := readFloat64FromReader(r); err != nil {
			return err
		}
		if _, err := readVarIntFromReader(r); err != nil {
			return err
		}
		_, err := readVarIntFromReader(r)
		return err
	})
}

func skipCustomModelData(r *bytes.Reader) error {
	if err := skipPrefixedArray(r, func(r *bytes.Reader) error { return skipFloat32(r) }); err != nil {
		return err
	}
	if err := skipPrefixedArray(r, func(r *bytes.Reader) error { _, err := readBoolFromReader(r); return err }); err != nil {
		return err
	}
	if err := skipPrefixedArray(r, func(r *bytes.Reader) error { _, err := readStringFromReader(r); return err }); err != nil {
		return err
	}
	return skipPrefixedArray(r, func(r *bytes.Reader) error { _, err := readInt32FromReader(r); return err })
}

func skipConsumable(r *bytes.Reader) error {
	if err := skipFloat32(r); err != nil {
		return err
	}
	if _, err := readVarIntFromReader(r); err != nil {
		return err
	}
	if err := skipIDOrX(r, skipSoundEvent); err != nil {
		return err
	}
	if _, err := readBoolFromReader(r); err != nil {
		return err
	}
	return skipPrefixedArray(r, skipConsumeEffect)
}

func skipTool(r *bytes.Reader) error {
	if err := skipPrefixedArray(r, func(r *bytes.Reader) error {
		if err := skipIDSet(r); err != nil {
			return err
		}
		if err := skipPrefixedOptional(r, func(r *bytes.Reader) error { return skipFloat32(r) }); err != nil {
			return err
		}
		return skipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := readBoolFromReader(r); return err })
	}); err != nil {
		return err
	}
	if err := skipFloat32(r); err != nil {
		return err
	}
	if _, err := readVarIntFromReader(r); err != nil {
		return err
	}
	_, err := readBoolFromReader(r)
	return err
}

func skipEquippable(r *bytes.Reader) error {
	if _, err := readVarIntFromReader(r); err != nil {
		return err
	}
	if err := skipIDOrX(r, skipSoundEvent); err != nil {
		return err
	}
	if err := skipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := readStringFromReader(r); return err }); err != nil {
		return err
	}
	if err := skipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := readStringFromReader(r); return err }); err != nil {
		return err
	}
	if err := skipPrefixedOptional(r, skipIDSet); err != nil {
		return err
	}
	if _, err := readBoolFromReader(r); err != nil {
		return err
	}
	if _, err := readBoolFromReader(r); err != nil {
		return err
	}
	_, err := readBoolFromReader(r)
	return err
}

func skipIDSet(r *bytes.Reader) error {
	setType, err := readVarIntFromReader(r)
	if err != nil {
		return err
	}
	if setType == 0 {
		_, err := readStringFromReader(r)
		return err
	}
	for i := int32(0); i < setType-1; i++ {
		if _, err := readVarIntFromReader(r); err != nil {
			return err
		}
	}
	return nil
}

func skipConsumeEffect(r *bytes.Reader) error {
	effectType, err := readVarIntFromReader(r)
	if err != nil {
		return err
	}
	switch effectType {
	case 0:
		if _, err := readVarIntFromReader(r); err != nil {
			return err
		}
		if err := skipFloat32(r); err != nil {
			return err
		}
		_, err := readVarIntFromReader(r)
		return err
	case 1:
		return skipPrefixedArray(r, skipConsumeEffect)
	case 2:
		_, err := readVarIntFromReader(r)
		return err
	}
	return nil
}

func skipPotionContents(r *bytes.Reader) error {
	if err := skipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := readVarIntFromReader(r); return err }); err != nil {
		return err
	}
	if err := skipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := readInt32FromReader(r); return err }); err != nil {
		return err
	}
	if err := skipPrefixedArray(r, skipPotionEffect); err != nil {
		return err
	}
	_, err := readStringFromReader(r)
	return err
}

func skipPotionEffect(r *bytes.Reader) error {
	if _, err := readVarIntFromReader(r); err != nil {
		return err
	}
	if _, err := readVarIntFromReader(r); err != nil {
		return err
	}
	_, err := readVarIntFromReader(r)
	return err
}

func skipWritableBookContent(r *bytes.Reader) error {
	num, err := readVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < num && i < 100; i++ {
		if _, err := readStringFromReader(r); err != nil {
			return err
		}
		if err := skipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := readStringFromReader(r); return err }); err != nil {
			return err
		}
	}
	return nil
}

func skipWrittenBookContent(r *bytes.Reader) error {
	if _, err := readStringFromReader(r); err != nil {
		return err
	}
	if err := skipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := readStringFromReader(r); return err }); err != nil {
		return err
	}
	if _, err := readStringFromReader(r); err != nil {
		return err
	}
	if _, err := readVarIntFromReader(r); err != nil {
		return err
	}
	num, err := readVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < num && i < 100; i++ {
		if err := skipNBT(r); err != nil {
			return err
		}
		if err := skipPrefixedOptional(r, skipNBT); err != nil {
			return err
		}
	}
	_, err = readBoolFromReader(r)
	return err
}

func skipTrim(r *bytes.Reader) error {
	if err := skipIDOrX(r, skipTrimMaterial); err != nil {
		return err
	}
	return skipIDOrX(r, skipTrimPattern)
}

func skipTrimMaterial(r *bytes.Reader) error {
	if _, err := readStringFromReader(r); err != nil {
		return err
	}
	if _, err := readStringFromReader(r); err != nil {
		return err
	}
	return skipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := readStringFromReader(r); return err })
}

func skipTrimPattern(r *bytes.Reader) error {
	if err := skipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := readStringFromReader(r); return err }); err != nil {
		return err
	}
	_, err := readStringFromReader(r)
	return err
}

func skipInstrument(r *bytes.Reader) error {
	if _, err := readStringFromReader(r); err != nil {
		return err
	}
	if _, err := readBoolFromReader(r); err != nil {
		return err
	}
	present, err := readBoolFromReader(r)
	if err != nil {
		return err
	}
	if present {
		return skipFloat32(r)
	}
	return nil
}

func skipJukeboxPlayable(r *bytes.Reader) error {
	mode, err := readU8(r)
	if err != nil {
		return err
	}
	if mode == 0 {
		_, err = readStringFromReader(r)
		return err
	}
	_, err = readVarIntFromReader(r)
	return err
}

func skipSoundEvent(r *bytes.Reader) error {
	if _, err := readStringFromReader(r); err != nil {
		return err
	}
	present, err := readBoolFromReader(r)
	if err != nil {
		return err
	}
	if present {
		return skipFloat32(r)
	}
	return nil
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

func skipSlotData(r *bytes.Reader) error {
	_, err := readSlotData(r)
	return err
}

func skipBlocksAttacks(r *bytes.Reader) error {
	if err := skipFloat32(r); err != nil {
		return err
	}
	if err := skipFloat32(r); err != nil {
		return err
	}
	if err := skipPrefixedArray(r, func(r *bytes.Reader) error {
		if err := skipFloat32(r); err != nil {
			return err
		}
		if err := skipPrefixedOptional(r, skipIDSet); err != nil {
			return err
		}
		if err := skipFloat32(r); err != nil {
			return err
		}
		if err := skipFloat32(r); err != nil {
			return err
		}
		if err := skipFloat32(r); err != nil {
			return err
		}
		return skipFloat32(r)
	}); err != nil {
		return err
	}
	if err := skipPrefixedOptional(r, func(r *bytes.Reader) error { _, err := readStringFromReader(r); return err }); err != nil {
		return err
	}
	if err := skipPrefixedOptional(r, func(r *bytes.Reader) error { return skipIDOrX(r, skipSoundEvent) }); err != nil {
		return err
	}
	return skipPrefixedOptional(r, func(r *bytes.Reader) error { return skipIDOrX(r, skipSoundEvent) })
}
