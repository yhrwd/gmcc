package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"gmcc/internal/logx"
)

type componentSkipper func(*bytes.Reader) error

var componentSkippers map[int32]componentSkipper

func init() {
	componentSkippers = map[int32]componentSkipper{
		0:  SkipVarInt,
		1:  SkipBannerPatterns,
		2:  SkipVarInt,
		3:  SkipNBT,
		4:  SkipNBT,
		5:  SkipBlockState,
		6:  SkipSoundEvent,
		7:  SkipNBT,
		8:  func(r *bytes.Reader) error { return SkipPrefixedArray(r, SkipSlotData) },
		9:  func(r *bytes.Reader) error { return SkipPrefixedArray(r, SkipSlotData) },
		10: func(r *bytes.Reader) error { return SkipPrefixedArray(r, SkipSlotData) },
		11: SkipNothing,
		12: SkipNothing,
		13: SkipNBT,
		14: SkipNBT,
		15: SkipVarInt,
		16: SkipString,
		17: SkipString,
		18: SkipDeathProtection,
		19: SkipNothing,
		20: SkipVarInt,
		21: SkipVarInt,
		22: func(r *bytes.Reader) error { _, err := ReadBoolFromReader(r); return err },
		23: SkipNBT,
		24: SkipFireworks,
		25: SkipFireworkExplosion,
		26: SkipNothing,
		27: SkipNothing,
		28: SkipNothing,
		29: SkipNothing,
		30: SkipNothing,
		31: SkipString,
		32: SkipNBT,
		33: SkipString,
		34: SkipLodestoneTracker,
		35: func(r *bytes.Reader) error { return SkipPrefixedArray(r, SkipNBT) },
		36: SkipInt32,
		37: SkipNothing,
		38: SkipVarInt,
		39: SkipVarInt,
		40: SkipVarInt,
		41: SkipVarInt,
		42: SkipFloat32,
		43: SkipString,
		44: SkipVarInt,
		45: SkipFloat32,
		46: func(r *bytes.Reader) error { return SkipPrefixedArray(r, SkipVarInt) },
		47: SkipString,
		48: SkipProvidesTrimMaterial,
		49: SkipVarInt,
		50: func(r *bytes.Reader) error { return SkipPrefixedArray(r, SkipString) },
		51: SkipRepairable,
		52: SkipVarInt,
		53: SkipSuspiciousStewEffects,
		54: SkipSwingAnimation,
		55: SkipTooltipDisplay,
		56: SkipString,
		57: SkipUseCooldown,
		58: func(r *bytes.Reader) error { return SkipSlotData(r) },
		59: SkipWritableBookContent,
		60: SkipWrittenBookContent,
	}
}

func SkipComponentByType(r *bytes.Reader, componentType int32) error {
	if skipper, ok := componentSkippers[componentType]; ok {
		return skipper(r)
	}
	logx.Warnf("未知的组件类型: %d, 尝试跳过剩余数据", componentType)
	return fmt.Errorf("unknown component type %d", componentType)
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
