package packet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"gmcc/internal/constants"
	"gmcc/internal/logx"
	"gmcc/internal/nbt"
)

// SlotData 兼容类型定义（避免循环依赖）
type SlotData struct {
	ID    int32
	Count int32
}

// MustReadBytes 读取字节，错误不终止但记录日志
// 用于解析非关键数据时忽略错误
func MustReadBytes(r io.Reader, n int, name string) []byte {
	b, err := ReadBytes(r, n)
	if err != nil {
		logx.PacketWarn(name, err)
	}
	return b
}

type byteReaderWrapper struct {
	r io.Reader
}

func newByteReaderWrapper(r io.Reader) *byteReaderWrapper {
	return &byteReaderWrapper{r: r}
}

func (b *byteReaderWrapper) Read(p []byte) (n int, err error) {
	return b.r.Read(p)
}

func (b *byteReaderWrapper) ReadByte() (byte, error) {
	var buf [1]byte
	if _, err := io.ReadFull(b.r, buf[:]); err != nil {
		return 0, err
	}
	return buf[0], nil
}

// MustReadVarInt 读取 VarInt，错误不终止但记录日志
func MustReadVarInt(r io.Reader, name string) int32 {
	v, err := ReadVarIntFromReader(r)
	if err != nil {
		logx.PacketWarn(name, err)
	}
	return v
}

// MustReadString 读取字符串，错误不终止但记录日志
func MustReadString(r io.Reader, name string) string {
	s, err := ReadStringFromReader(r)
	if err != nil {
		logx.PacketWarn(name, err)
	}
	return s
}

// MustReadBool 读取布尔值，错误不终止但记录日志
func MustReadBool(r io.Reader, name string) bool {
	v, err := ReadBool(r)
	if err != nil {
		logx.PacketWarn(name, err)
	}
	return v
}

// MustReadU8 读取 uint8，错误不终止但记录日志
func MustReadU8(r io.Reader, name string) byte {
	v, err := ReadU8(r)
	if err != nil {
		logx.PacketWarn(name, err)
	}
	return v
}

// MustReadInt32 读取 int32，错误不终止但记录日志
func MustReadInt32(r io.Reader, name string) int32 {
	v, err := ReadInt32(r)
	if err != nil {
		logx.PacketWarn(name, err)
	}
	return v
}

// MustReadFloat64 读取 float64，错误不终止但记录日志
func MustReadFloat64(r io.Reader, name string) float64 {
	v, err := ReadFloat64FromReader(r)
	if err != nil {
		logx.PacketWarn(name, err)
	}
	return v
}

func ReadVarIntFromReader(r io.Reader) (int32, error) {
	if br, ok := r.(io.ByteReader); ok {
		return ReadVarInt(br)
	}
	return ReadVarInt(newByteReaderWrapper(r))
}

func ReadStringFromReader(r io.Reader) (string, error) {
	if br, ok := r.(io.ByteReader); ok {
		return ReadString(br, r)
	}
	return ReadString(newByteReaderWrapper(r), r)
}

func ReadBoolFromReader(r io.Reader) (bool, error) {
	return ReadBool(r)
}

func ReadInt32FromReader(r io.Reader) (int32, error) {
	return ReadInt32(r)
}

func ReadFloat64FromReader(r io.Reader) (float64, error) {
	var v float64
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

func ReadFloat32FromReader(r io.Reader) (float32, error) {
	var v float32
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

func ReadU8(r io.Reader) (byte, error) {
	var b [1]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return 0, err
	}
	return b[0], nil
}

func ReadBytes(r io.Reader, n int) ([]byte, error) {
	if n < 0 {
		return nil, fmt.Errorf("negative read length: %d", n)
	}
	if n > constants.MaxPacketSize {
		return nil, fmt.Errorf("read length exceeds max allowed: %d > %d", n, constants.MaxPacketSize)
	}
	b := make([]byte, n)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, fmt.Errorf("read %d bytes: %w", n, err)
	}
	return b, nil
}

// ReadSlotData 读取物品槽数据 (兼容旧接口)
// 注意: 新代码应直接使用 internal/item.ReadSlotData 获取完整组件信息
func ReadSlotData(r *bytes.Reader) (*SlotData, error) {
	// 1. item_count (VarInt)
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil // 空物品
	}

	// 2. item_id (VarInt)
	itemID, err := ReadVarIntFromReader(r)
	if err != nil {
		return nil, err
	}

	// 3. 跳过组件 (使用新的组件系统)
	if err := skipSlotComponents(r); err != nil {
		logx.Warnf("Slot解析失败: itemID=%d, count=%d, err=%v", itemID, count, err)
		return nil, err
	}

	return &SlotData{ID: itemID, Count: count}, nil
}

func skipSlotData(r *bytes.Reader) error {
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	return skipSlotComponents(r)
}

// skipSlotComponents 跳过物品组件 (内部使用，简化版)
func skipSlotComponents(r *bytes.Reader) error {
	// 添加的组件数量
	numAdd, err := ReadVarIntFromReader(r)
	if err != nil {
		return fmt.Errorf("read num_add: %w", err)
	}

	// 移除的组件数量
	numRemove, err := ReadVarIntFromReader(r)
	if err != nil {
		return fmt.Errorf("read num_remove: %w", err)
	}

	// 跳过添加的组件
	for i := int32(0); i < numAdd; i++ {
		componentType, err := ReadVarIntFromReader(r)
		if err != nil {
			return fmt.Errorf("read component type at index %d: %w", i, err)
		}
		if err := skipComponentByType(r, componentType); err != nil {
			return fmt.Errorf("skip component type %d at index %d: %w", componentType, i, err)
		}
	}

	// 跳过移除的组件 (只有 component_type)
	for i := int32(0); i < numRemove; i++ {
		if _, err := ReadVarIntFromReader(r); err != nil {
			return fmt.Errorf("read removed component type at index %d: %w", i, err)
		}
	}
	return nil
}

// skipComponentByType 根据组件类型跳过数据 (参考 1.21.11 数据组件规范)
func skipComponentByType(r *bytes.Reader, componentType int32) error {
	switch componentType {
	case 0, 6, 9, 45, 55, 57, 64, 76, 77:
		return SkipNBT(r)
	case 1, 2, 3, 12, 19, 31, 44, 46, 61, 71, 79, 80, 81, 82, 83, 85, 86, 87, 88, 89, 90, 91, 92, 93, 95, 98, 99, 100, 101, 102, 103:
		_, err := ReadVarIntFromReader(r)
		return err
	case 4, 20, 22, 34:
		return nil
	case 5:
		return skipUseEffects(r)
	case 7:
		_, err := ReadFloat32FromReader(r)
		return err
	case 10, 27, 35, 63, 69:
		_, err := ReadStringFromReader(r)
		return err
	case 11:
		return skipNBTList(r)
	case 13, 41:
		return skipEnchantmentList(r)
	case 14, 15:
		return skipItemBlockPredicates(r)
	case 16:
		return skipAttributeModifiers(r)
	case 17:
		return skipCustomModelData(r)
	case 18:
		return skipTooltipDisplay(r)
	case 21:
		_, err := ReadBoolFromReader(r)
		return err
	case 23:
		return skipFoodComponent(r)
	case 24, 36:
		return skipConsumeEffectsComponent(r)
	case 25:
		return skipSlotData(r)
	case 26:
		return skipUseCooldown(r)
	case 28:
		return skipToolComponent(r)
	case 29:
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
		_, err := ReadFloat32FromReader(r)
		return err
	case 30:
		return skipFloat32Values(r, 6)
	case 32:
		return skipEquippableComponent(r)
	case 33:
		return skipIDSetApprox(r)
	case 37:
		return skipBlocksAttacksComponent(r)
	case 38:
		return skipPiercingWeaponComponent(r)
	case 39:
		return skipKineticWeaponComponent(r)
	case 40:
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
		_, err := ReadVarIntFromReader(r)
		return err
	case 42, 43:
		_, err := ReadInt32FromReader(r)
		return err
	case 47:
		_, err := ReadFloat32FromReader(r)
		return err
	case 48, 49, 73:
		return skipSlotList(r)
	case 50:
		return skipPotionContents(r)
	case 51:
		return skipSuspiciousStewEffects(r)
	case 52:
		return skipWritableBookContent(r)
	case 53:
		return skipWrittenBookContent(r)
	case 54:
		return skipTrimComponentApprox(r)
	case 56, 58:
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
		return SkipNBT(r)
	case 59:
		return skipHolderOrStringApprox(r)
	case 60:
		return skipHolderOrStringApprox(r)
	case 62:
		return skipHolderOrStringApprox(r)
	case 65:
		return skipLodestoneTracker(r)
	case 66:
		return skipFireworkExplosion(r)
	case 67:
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
		count, err := ReadVarIntFromReader(r)
		if err != nil {
			return err
		}
		for i := int32(0); i < count; i++ {
			if err := skipFireworkExplosion(r); err != nil {
				return err
			}
		}
		return nil
	case 68:
		return skipResolvableProfile(r)
	case 70:
		count, err := ReadVarIntFromReader(r)
		if err != nil {
			return err
		}
		for i := int32(0); i < count; i++ {
			if err := skipBannerPatternLayerApprox(r); err != nil {
				return err
			}
		}
		return nil
	case 75:
		return skipBeesComponent(r)
	case 78:
		return skipSoundHolderApprox(r)
	case 84:
		return skipRegistryEntryHolderApprox(r)
	case 94:
		return skipTrimComponentApprox(r)
	case 96, 97:
		return skipRegistryEntryHolderApprox(r)
	default:
		return SkipNBT(r)
	}
}

// skipContainerComponentData 跳过容器组件数据
func skipContainerComponentData(r *bytes.Reader) error {
	return skipSlotList(r)
}

func skipUseEffects(r *bytes.Reader) error {
	if _, err := ReadBoolFromReader(r); err != nil {
		return err
	}
	if _, err := ReadBoolFromReader(r); err != nil {
		return err
	}
	_, err := ReadFloat32FromReader(r)
	return err
}

func skipNBTList(r *bytes.Reader) error {
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < count; i++ {
		if err := SkipNBT(r); err != nil {
			return err
		}
	}
	return nil
}

func skipEnchantmentList(r *bytes.Reader) error {
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < count; i++ {
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
	}
	return nil
}

func skipCustomModelData(r *bytes.Reader) error {
	if err := skipFloat32List(r); err != nil {
		return err
	}
	if err := skipBoolList(r); err != nil {
		return err
	}
	if err := skipStringList(r); err != nil {
		return err
	}
	return skipInt32List(r)
}

func skipTooltipDisplay(r *bytes.Reader) error {
	if _, err := ReadBoolFromReader(r); err != nil {
		return err
	}
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < count; i++ {
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
	}
	return nil
}

func skipFoodComponent(r *bytes.Reader) error {
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	if _, err := ReadFloat32FromReader(r); err != nil {
		return err
	}
	_, err := ReadBoolFromReader(r)
	return err
}

func skipConsumeEffectsComponent(r *bytes.Reader) error {
	if _, err := ReadFloat32FromReader(r); err != nil {
		return err
	}
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	if err := skipSoundHolderApprox(r); err != nil {
		return err
	}
	if _, err := ReadBoolFromReader(r); err != nil {
		return err
	}
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < count; i++ {
		if err := skipItemConsumeEffect(r); err != nil {
			return err
		}
	}
	return nil
}

func skipUseCooldown(r *bytes.Reader) error {
	if _, err := ReadFloat32FromReader(r); err != nil {
		return err
	}
	return skipOptionalString(r)
}

func skipToolComponent(r *bytes.Reader) error {
	rules, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < rules; i++ {
		if err := skipIDSetApprox(r); err != nil {
			return err
		}
		if err := skipOptionalFloat32(r); err != nil {
			return err
		}
		if err := skipOptionalBool(r); err != nil {
			return err
		}
	}
	if _, err := ReadFloat32FromReader(r); err != nil {
		return err
	}
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	_, err = ReadBoolFromReader(r)
	return err
}

func skipEquippableComponent(r *bytes.Reader) error {
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	if err := skipSoundHolderApprox(r); err != nil {
		return err
	}
	if err := skipOptionalString(r); err != nil {
		return err
	}
	if err := skipOptionalString(r); err != nil {
		return err
	}
	if err := skipOptionalIDSetApprox(r); err != nil {
		return err
	}
	for i := 0; i < 5; i++ {
		if _, err := ReadBoolFromReader(r); err != nil {
			return err
		}
	}
	return skipSoundHolderApprox(r)
}

func skipBlocksAttacksComponent(r *bytes.Reader) error {
	if _, err := ReadFloat32FromReader(r); err != nil {
		return err
	}
	if _, err := ReadFloat32FromReader(r); err != nil {
		return err
	}
	damageReductions, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < damageReductions; i++ {
		if _, err := ReadFloat32FromReader(r); err != nil {
			return err
		}
		if err := skipOptionalIDSetApprox(r); err != nil {
			return err
		}
		if err := skipFloat32Values(r, 2); err != nil {
			return err
		}
	}
	if err := skipFloat32Values(r, 3); err != nil {
		return err
	}
	if err := skipOptionalString(r); err != nil {
		return err
	}
	if err := skipOptionalSoundHolderApprox(r); err != nil {
		return err
	}
	return skipOptionalSoundHolderApprox(r)
}

func skipPiercingWeaponComponent(r *bytes.Reader) error {
	if _, err := ReadBoolFromReader(r); err != nil {
		return err
	}
	if _, err := ReadBoolFromReader(r); err != nil {
		return err
	}
	if err := skipOptionalSoundHolderApprox(r); err != nil {
		return err
	}
	return skipOptionalSoundHolderApprox(r)
}

func skipKineticWeaponComponent(r *bytes.Reader) error {
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	for i := 0; i < 3; i++ {
		if err := skipOptionalKineticWeaponCondition(r); err != nil {
			return err
		}
	}
	if err := skipFloat32Values(r, 2); err != nil {
		return err
	}
	if err := skipOptionalSoundHolderApprox(r); err != nil {
		return err
	}
	return skipOptionalSoundHolderApprox(r)
}

func skipPotionContents(r *bytes.Reader) error {
	if err := skipOptionalVarInt(r); err != nil {
		return err
	}
	if err := skipOptionalInt32(r); err != nil {
		return err
	}
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < count; i++ {
		if err := skipItemPotionEffect(r); err != nil {
			return err
		}
	}
	return skipOptionalString(r)
}

func skipSuspiciousStewEffects(r *bytes.Reader) error {
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < count; i++ {
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
	}
	return nil
}

func skipWritableBookContent(r *bytes.Reader) error {
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < count; i++ {
		if _, err := ReadStringFromReader(r); err != nil {
			return err
		}
		if err := skipOptionalString(r); err != nil {
			return err
		}
	}
	return nil
}

func skipWrittenBookContent(r *bytes.Reader) error {
	if _, err := ReadStringFromReader(r); err != nil {
		return err
	}
	if err := skipOptionalString(r); err != nil {
		return err
	}
	if _, err := ReadStringFromReader(r); err != nil {
		return err
	}
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < count; i++ {
		if err := SkipNBT(r); err != nil {
			return err
		}
		if err := skipOptionalNBT(r); err != nil {
			return err
		}
	}
	_, err = ReadBoolFromReader(r)
	return err
}

func skipTrimComponentApprox(r *bytes.Reader) error {
	if err := skipRegistryEntryHolderApprox(r); err != nil {
		return err
	}
	return skipRegistryEntryHolderApprox(r)
}

func skipLodestoneTracker(r *bytes.Reader) error {
	if err := skipOptionalGlobalPos(r); err != nil {
		return err
	}
	_, err := ReadBoolFromReader(r)
	return err
}

func skipFireworkExplosion(r *bytes.Reader) error {
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	if err := skipInt32List(r); err != nil {
		return err
	}
	if err := skipInt32List(r); err != nil {
		return err
	}
	if _, err := ReadBoolFromReader(r); err != nil {
		return err
	}
	_, err := ReadBoolFromReader(r)
	return err
}

func skipResolvableProfile(r *bytes.Reader) error {
	profileType, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	switch profileType {
	case 0:
		if err := skipOptionalString(r); err != nil {
			return err
		}
		if err := skipOptionalUUID(r); err != nil {
			return err
		}
	case 1:
		if err := skipUUID(r); err != nil {
			return err
		}
		if _, err := ReadStringFromReader(r); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown resolvable profile type %d", profileType)
	}

	properties, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < properties; i++ {
		if _, err := ReadStringFromReader(r); err != nil {
			return err
		}
		if _, err := ReadStringFromReader(r); err != nil {
			return err
		}
		if err := skipOptionalString(r); err != nil {
			return err
		}
	}

	for i := 0; i < 3; i++ {
		if err := skipOptionalString(r); err != nil {
			return err
		}
	}
	return skipOptionalVarInt(r)
}

func skipBeesComponent(r *bytes.Reader) error {
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < count; i++ {
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

func skipItemBlockPredicates(r *bytes.Reader) error {
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < count; i++ {
		if err := skipOptionalIDSetApprox(r); err != nil {
			return err
		}
		hasProperties, err := ReadBoolFromReader(r)
		if err != nil {
			return err
		}
		if hasProperties {
			propertyCount, err := ReadVarIntFromReader(r)
			if err != nil {
				return err
			}
			for j := int32(0); j < propertyCount; j++ {
				if _, err := ReadStringFromReader(r); err != nil {
					return err
				}
				exact, err := ReadBoolFromReader(r)
				if err != nil {
					return err
				}
				if _, err := ReadStringFromReader(r); err != nil {
					return err
				}
				if !exact {
					if _, err := ReadStringFromReader(r); err != nil {
						return err
					}
				}
			}
		}
		if err := skipOptionalNBT(r); err != nil {
			return err
		}
		exactMatchers, err := ReadVarIntFromReader(r)
		if err != nil {
			return err
		}
		for j := int32(0); j < exactMatchers; j++ {
			componentType, err := ReadVarIntFromReader(r)
			if err != nil {
				return err
			}
			if err := skipComponentByType(r, componentType); err != nil {
				return err
			}
		}
		partialMatchers, err := ReadVarIntFromReader(r)
		if err != nil {
			return err
		}
		for j := int32(0); j < partialMatchers; j++ {
			if _, err := ReadVarIntFromReader(r); err != nil {
				return err
			}
		}
	}
	return nil
}

func skipAttributeModifiers(r *bytes.Reader) error {
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < count; i++ {
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
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
	}
	displayType, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	if displayType == 2 {
		return SkipNBT(r)
	}
	return nil
}

func skipItemPotionEffect(r *bytes.Reader) error {
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	return skipItemEffectDetail(r)
}

func skipItemEffectDetail(r *bytes.Reader) error {
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	for i := 0; i < 3; i++ {
		if _, err := ReadBoolFromReader(r); err != nil {
			return err
		}
	}
	hasHidden, err := ReadBoolFromReader(r)
	if err != nil {
		return err
	}
	if hasHidden {
		return skipItemEffectDetail(r)
	}
	return nil
}

func skipItemConsumeEffect(r *bytes.Reader) error {
	effectType, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	switch effectType {
	case 0:
		count, err := ReadVarIntFromReader(r)
		if err != nil {
			return err
		}
		for i := int32(0); i < count; i++ {
			if err := skipItemPotionEffect(r); err != nil {
				return err
			}
		}
		_, err = ReadFloat32FromReader(r)
		return err
	case 1:
		return skipIDSetApprox(r)
	case 2:
		return nil
	case 3:
		_, err := ReadFloat32FromReader(r)
		return err
	case 4:
		return skipSoundHolderApprox(r)
	default:
		return fmt.Errorf("unknown consume effect type %d", effectType)
	}
}

func skipSoundHolderApprox(r *bytes.Reader) error {
	_, err := ReadVarIntFromReader(r)
	return err
}

func skipOptionalSoundHolderApprox(r *bytes.Reader) error {
	hasValue, err := ReadBoolFromReader(r)
	if err != nil {
		return err
	}
	if !hasValue {
		return nil
	}
	return skipSoundHolderApprox(r)
}

func skipHolderOrStringApprox(r *bytes.Reader) error {
	hasHolder, err := ReadBoolFromReader(r)
	if err != nil {
		return err
	}
	if hasHolder {
		return skipRegistryEntryHolderApprox(r)
	}
	_, err = ReadStringFromReader(r)
	return err
}

func skipRegistryEntryHolderApprox(r *bytes.Reader) error {
	_, err := ReadVarIntFromReader(r)
	return err
}

func skipIDSetApprox(r *bytes.Reader) error {
	lengthOrCount, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	if lengthOrCount < 0 {
		return fmt.Errorf("invalid IDSet marker: %d", lengthOrCount)
	}

	if lengthOrCount == 0 {
		looksLikeName, strLen, err := peekIdentifierStringLength(r)
		if err != nil {
			return err
		}
		if looksLikeName {
			return DiscardN(r, strLen)
		}
		return nil
	}

	looksLikeName, strLen, err := peekFixedLengthIdentifier(r, int(lengthOrCount))
	if err != nil {
		return err
	}
	if looksLikeName {
		return DiscardN(r, strLen)
	}

	for i := int32(0); i < lengthOrCount; i++ {
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
	}
	return nil
}

func skipOptionalIDSetApprox(r *bytes.Reader) error {
	hasValue, err := ReadBoolFromReader(r)
	if err != nil {
		return err
	}
	if !hasValue {
		return nil
	}
	return skipIDSetApprox(r)
}

func skipBannerPatternLayerApprox(r *bytes.Reader) error {
	if err := skipRegistryEntryHolderApprox(r); err != nil {
		return err
	}
	_, err := ReadVarIntFromReader(r)
	return err
}

func skipOptionalGlobalPos(r *bytes.Reader) error {
	hasValue, err := ReadBoolFromReader(r)
	if err != nil {
		return err
	}
	if !hasValue {
		return nil
	}
	if _, err := ReadStringFromReader(r); err != nil {
		return err
	}
	return DiscardN(r, 8)
}

func skipOptionalUUID(r *bytes.Reader) error {
	hasValue, err := ReadBoolFromReader(r)
	if err != nil {
		return err
	}
	if !hasValue {
		return nil
	}
	return skipUUID(r)
}

func skipUUID(r io.Reader) error {
	_, err := ReadBytes(r, 16)
	return err
}

func skipOptionalString(r *bytes.Reader) error {
	hasValue, err := ReadBoolFromReader(r)
	if err != nil {
		return err
	}
	if !hasValue {
		return nil
	}
	_, err = ReadStringFromReader(r)
	return err
}

func skipOptionalNBT(r *bytes.Reader) error {
	hasValue, err := ReadBoolFromReader(r)
	if err != nil {
		return err
	}
	if !hasValue {
		return nil
	}
	return SkipNBT(r)
}

func skipOptionalVarInt(r *bytes.Reader) error {
	hasValue, err := ReadBoolFromReader(r)
	if err != nil {
		return err
	}
	if !hasValue {
		return nil
	}
	_, err = ReadVarIntFromReader(r)
	return err
}

func skipOptionalInt32(r *bytes.Reader) error {
	hasValue, err := ReadBoolFromReader(r)
	if err != nil {
		return err
	}
	if !hasValue {
		return nil
	}
	_, err = ReadInt32FromReader(r)
	return err
}

func skipOptionalFloat32(r *bytes.Reader) error {
	hasValue, err := ReadBoolFromReader(r)
	if err != nil {
		return err
	}
	if !hasValue {
		return nil
	}
	_, err = ReadFloat32FromReader(r)
	return err
}

func skipOptionalBool(r *bytes.Reader) error {
	hasValue, err := ReadBoolFromReader(r)
	if err != nil {
		return err
	}
	if !hasValue {
		return nil
	}
	_, err = ReadBoolFromReader(r)
	return err
}

func skipOptionalKineticWeaponCondition(r *bytes.Reader) error {
	hasValue, err := ReadBoolFromReader(r)
	if err != nil {
		return err
	}
	if !hasValue {
		return nil
	}
	if _, err := ReadVarIntFromReader(r); err != nil {
		return err
	}
	if _, err := ReadFloat32FromReader(r); err != nil {
		return err
	}
	_, err = ReadFloat32FromReader(r)
	return err
}

func peekIdentifierStringLength(r *bytes.Reader) (bool, int, error) {
	start, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return false, 0, err
	}
	defer func() {
		_, _ = r.Seek(start, io.SeekStart)
	}()

	length, err := ReadVarIntFromReader(r)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return false, 0, nil
		}
		return false, 0, err
	}
	if length <= 0 {
		return false, 0, nil
	}
	return peekFixedLengthIdentifier(r, int(length))
}

func peekFixedLengthIdentifier(r *bytes.Reader, length int) (bool, int, error) {
	if length <= 0 || length > r.Len() {
		return false, 0, nil
	}

	start, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return false, 0, err
	}
	defer func() {
		_, _ = r.Seek(start, io.SeekStart)
	}()

	data, err := ReadBytes(r, length)
	if err != nil {
		return false, 0, err
	}
	if !looksLikeIdentifierBytes(data) {
		return false, 0, nil
	}
	return true, length, nil
}

func looksLikeIdentifierBytes(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	for _, b := range data {
		if b >= 'a' && b <= 'z' {
			continue
		}
		if b >= 'A' && b <= 'Z' {
			continue
		}
		if b >= '0' && b <= '9' {
			continue
		}
		switch b {
		case ':', '/', '_', '-', '.', '#':
			continue
		default:
			return false
		}
	}
	return true
}

func skipSlotList(r *bytes.Reader) error {
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < count; i++ {
		if err := skipSlotData(r); err != nil {
			return err
		}
	}
	return nil
}

func skipFloat32Values(r *bytes.Reader, count int) error {
	for i := 0; i < count; i++ {
		if _, err := ReadFloat32FromReader(r); err != nil {
			return err
		}
	}
	return nil
}

func skipFloat32List(r *bytes.Reader) error {
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < count; i++ {
		if _, err := ReadFloat32FromReader(r); err != nil {
			return err
		}
	}
	return nil
}

func skipBoolList(r *bytes.Reader) error {
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < count; i++ {
		if _, err := ReadBoolFromReader(r); err != nil {
			return err
		}
	}
	return nil
}

func skipStringList(r *bytes.Reader) error {
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < count; i++ {
		if _, err := ReadStringFromReader(r); err != nil {
			return err
		}
	}
	return nil
}

func skipInt32List(r *bytes.Reader) error {
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < count; i++ {
		if _, err := ReadInt32FromReader(r); err != nil {
			return err
		}
	}
	return nil
}

// SkipNBT 跳过 Network NBT 格式 (无 name 字段)
func SkipNBT(r *bytes.Reader) error {
	if r.Len() == 0 {
		return nil
	}
	dec := nbt.NewDecoder(r).NetworkFormat(true)
	err := dec.Skip()
	if err != nil {
		errMsg := err.Error()
		if errMsg == "unexpected EOF" || strings.HasPrefix(errMsg, "unknown tag type: ") {
			logx.Warnf("SkipNBT 警告: %v, 剩余 %d 字节", err, r.Len())
			return nil
		}
		return err
	}
	return nil
}

// ReadAnonymousNBTJSON 解析 Network NBT 并返回 JSON 字符串
func ReadAnonymousNBTJSON(r io.Reader) (string, error) {
	dec := nbt.NewDecoder(r)
	dec.NetworkFormat(true)
	var v any
	if err := dec.Decode(&v); err != nil {
		return "", err
	}

	raw, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
