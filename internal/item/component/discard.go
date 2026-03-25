package component

import (
	"bytes"

	"gmcc/internal/logx"
	"gmcc/internal/mcclient/packet"
)

func makeDiscardHandler(typeID int32) ComponentHandler {
	return func(id int32, r *bytes.Reader) (*ComponentResult, error) {
		err := skipComponentByID(id, r)
		if err != nil {
			logx.Debugf("跳过组件失败: typeID=%d, err=%v", id, err)
		}
		return &ComponentResult{TypeID: id}, nil
	}
}

func skipComponentByID(id int32, r *bytes.Reader) error {
	switch id {
	case Unbreakable, CreativeSlotLock, Glider:
		return nil
	case DyedColor, MapColor, CustomModelData, EnchantmentGlintOverride,
		PotionDurationScale, UseEffects:
		_, err := packet.ReadInt32FromReader(r)
		return err
	case MaxStackSize, MaxDamage, Damage, MinimumAttackCharge, Rarity,
		RepairCost, Enchantable, MapID, MapPostProcessing,
		OminousBottleAmplifier, BundleRemainingSpace, FrameType,
		BaseColorComponent, ColorComponent:
		_, err := packet.ReadVarIntFromReader(r)
		return err
	default:
		return packet.SkipNBT(r)
	}
}
