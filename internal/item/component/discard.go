package component

import (
	"bytes"

	"gmcc/internal/logx"
	"gmcc/internal/mcclient/packet"
)

// makeDiscardHandler 创建丢弃处理器
func makeDiscardHandler(typeID int32) ComponentHandler {
	return func(id int32, r *bytes.Reader) (*ComponentResult, error) {
		// 根据组件类型选择跳过方式
		var err error
		switch id {
		// ===== NBT 组件 - 最常见的组件类型 =====
		// ID: 0-58 (大部分组件都是 NBT 格式)
		case CustomData, // 0 - 自定义数据
			MaxStackSize,             // 1 - 最大堆叠数
			MaxDamage,                // 2 - 最大耐久
			Damage,                   // 3 - 当前耐久
			Unbreakable,              // 4 - 不可破坏 (特殊: 无数据)
			UseEffects,               // 5 - 使用效果
			CustomName,               // 6 - 自定义名称
			MinimumAttackCharge,      // 7 - 最小攻击蓄力
			DamageType,               // 8 - 伤害类型
			ItemName,                 // 9 - 物品名称
			ItemModel,                // 10 - 物品模型 (特殊: String)
			Lore,                     // 11 - 物品描述
			Rarity,                   // 12 - 稀有度
			Enchantments,             // 13 - 附魔
			CanPlaceOn,               // 14 - 可放置于
			CanBreak,                 // 15 - 可破坏
			AttributeModifiers,       // 16 - 属性修饰符
			CustomModelData,          // 17 - 自定义模型数据
			TooltipDisplay,           // 18 - 提示框显示
			RepairCost,               // 19 - 修复成本
			CreativeSlotLock,         // 20 - 创造模式槽位锁定 (特殊: 无数据)
			EnchantmentGlintOverride, // 21 - 附魔光效覆盖 (特殊: Bool)
			IntangibleProjectile,     // 22 - 无形弹射物
			Food,                     // 23 - 食物属性
			Consumable,               // 24 - 可消耗
			UseRemainder,             // 25 - 使用剩余
			UseCooldown,              // 26 - 使用冷却
			DamageResistant,          // 27 - 伤害抗性
			Tool,                     // 28 - 工具属性
			Weapon,                   // 29 - 武器属性
			AttackRange,              // 30 - 攻击范围
			Enchantable,              // 31 - 可附魔
			Equippable,               // 32 - 可穿戴
			Repairable,               // 33 - 可修复
			Glider,                   // 34 - 滑翔
			TooltipStyle,             // 35 - 提示框样式
			DeathProtection,          // 36 - 死亡保护
			BlocksAttacks,            // 37 - 阻挡攻击
			PiercingWeapon,           // 38 - 穿透武器
			KineticWeapon,            // 39 - 动能武器
			SwingAnimation,           // 40 - 挥动动画
			StoredEnchantments,       // 41 - 存储的附魔
			DyedColor,                // 42 - 染色颜色 (特殊: Int32)
			MapColor,                 // 43 - 地图颜色 (特殊: Int32)
			MapID,                    // 44 - 地图ID
			MapDecorations,           // 45 - 地图标记
			MapPostProcessing,        // 46 - 地图后处理
			PotionDurationScale,      // 47 - 药水持续时间缩放 (特殊: Int32)
			ChargedProjectiles,       // 48 - 已充能的弹射物
			BundleContents,           // 49 - 收纳袋内容
			PotionContents,           // 50 - 药水内容
			SuspiciousStewEffects,    // 51 - 迷之炖菜效果
			WritableBookContent,      // 52 - 可写书内容
			WrittenBookContent,       // 53 - 成书内容
			Trim,                     // 54 - 盔甲纹饰
			DebugStickState,          // 55 - 调试棒状态
			EntityData,               // 56 - 实体数据
			BucketEntityData,         // 57 - 桶中实体数据
			BlockEntityData,          // 58 - 方块实体数据
			// 以下继续是 NBT 组件
			Instrument,             // 59
			ProvidesTrimMaterial,   // 60
			OminousBottleAmplifier, // 61
			JukeboxPlayable,        // 62
			ProvidesBannerPatterns, // 63 (注意：String类型，但在 NBT 组中)
			Recipes,                // 64
			LodestoneTracker,       // 65
			FireworkExplosion,      // 66
			Fireworks,              // 67
			Profile,                // 68
			NoteBlockSound,         // 69 (注意：String类型，但在 NBT 组中)
			BannerPatterns,         // 70
			BaseColor,              // 71
			PotDecorations,         // 72
			// 73 是 Container，需要特殊处理
			BlockState,           // 74
			Bees,                 // 75
			Lock,                 // 76
			ContainerLoot,        // 77
			BreakSound,           // 78 (注意：String类型，但在 NBT 组中)
			VillagerVariant,      // 79
			WolfVariant,          // 80
			CatVariant,           // 81
			FrogVariant,          // 82
			AxolotlVariant,       // 83
			PaintingVariant,      // 84
			ShulkerVariant,       // 85
			GoatVariant,          // 86
			SnifferVariant,       // 87
			GhoulVariant,         // 88
			BreezeVariant,        // 89
			BundleRemainingSpace, // 90 (注意：VarInt类型，但在 NBT 组中)
			EntityColor,          // 92
			Buckable,             // 93 (注意：Bool类型，但在 NBT 组中)
			ArmorTrim,            // 94
			EquippableColor,      // 95
			TrimMaterial,         // 96
			TrimPattern,          // 97
			CompassColor,         // 98
			MapDisplayColor,      // 99
			FrameType,            // 100
			BannerPattern,        // 101
			BaseColorComponent,   // 102
			ColorComponent,       // 103
			// NBT 组件
			err = packet.SkipNBT(r)

		// ===== 特殊处理 =====
		// Container 组件 (73) - 这个已经在 handlers.go 中特殊处理
		case Container:
			// 这里不应该被执行，因为 handlers.go 中 Container 有特殊处理
			// 但如果执行到这里，尝试跳过 NBT
			err = packet.SkipNBT(r)

		default:
			// 其他组件：尝试作为 NBT 跳过
			logx.Debugf("未知组件类型 %d，尝试作为 NBT 跳过", id)
			err = packet.SkipNBT(r)
		}

		// SkipNBT 已处理 EOF 错误，这里只需记录其他错误
		if err != nil {
			logx.Debugf("跳过组件失败: typeID=%d, err=%v", id, err)
			// 不返回错误，继续处理
		}

		return &ComponentResult{TypeID: id}, nil
	}
}
