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
		// 读取 component_type
		componentType, err := ReadVarIntFromReader(r)
		if err != nil {
			// 如果读取组件类型失败，可能是数据损坏，尝试跳过剩余部分
			logx.Debugf("Failed to read component type at index %d: %v, remaining bytes: %d", i, err, r.Len())
			// 尝试通过 NBT 跳过剩余部分
			if skipErr := SkipNBT(r); skipErr != nil {
				// 如果 NBT 也失败，跳过所有剩余字节
				logx.Debugf("SkipNBT also failed: %v, skipping all remaining bytes", skipErr)
				if r.Len() > 0 {
					_, _ = r.Seek(int64(r.Len()), 1)
				}
			}
			return fmt.Errorf("read component type at index %d: %w", i, err)
		}
		// 根据类型跳过数据
		if err := skipComponentByType(r, componentType); err != nil {
			logx.Debugf("Failed to skip component type %d at index %d: %v", componentType, i, err)
			// 如果跳过失败，尝试跳到下一个组件
			continue
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
	// VarInt 类型 (基础数字、稀有度、附魔能力等)
	varIntTypes := map[int32]bool{
		1:   true, // max_stack_size
		2:   true, // max_damage
		3:   true, // damage
		7:   true, // minimum_attack_charge
		12:  true, // rarity
		19:  true, // repair_cost
		20:  true, // creative_slot_lock
		31:  true, // enchantable
		44:  true, // map_id
		46:  true, // map_post_processing
		47:  true, // potion_duration_scale
		61:  true, // ominous_bottle_amplifier
		91:  true, // bundle_remaining_space
		102: true, // base_color_component
		103: true, // color_component
		54:  true, // trimming_material (VarInt for material ID)
		73:  true, // container (size VarInt)
	}

	if varIntTypes[componentType] {
		_, err := ReadVarIntFromReader(r)
		return err
	}

	// Int32 类型 (颜色、定制模型数据等)
	int32Types := map[int32]bool{
		17:  true, // custom_model_data
		42:  true, // dyed_color
		43:  true, // map_color
		71:  true, // base_color
		100: true, // frame_type
	}

	if int32Types[componentType] {
		_, err := ReadInt32FromReader(r)
		return err
	}

	// Bool 类型 (无数据、附魔光效等)
	boolTypes := map[int32]bool{
		4:  true, // unbreakable
		21: true, // enchantment_glint_override
		34: true, // glider
		36: true, // death_protection
	}

	if boolTypes[componentType] {
		_, err := ReadBoolFromReader(r)
		return err
	}

	// NBT 类型 (文本组件、复杂结构等)
	nbtTypes := map[int32]bool{
		0:   true, // custom_data
		6:   true, // custom_name
		9:   true, // item_name
		11:  true, // lore
		13:  true, // enchantments
		22:  true, // intangible_projectile
		23:  true, // food
		24:  true, // consumable
		25:  true, // use_remainder
		26:  true, // use_cooldown
		27:  true, // damage_resistant
		28:  true, // tool
		29:  true, // weapon
		30:  true, // attack_range
		32:  true, // equippable
		33:  true, // repairable
		35:  true, // tooltip_style
		37:  true, // blocks_attacks
		38:  true, // piercing_weapon
		39:  true, // kinetic_weapon
		40:  true, // swing_animation
		41:  true, // stored_enchantments
		45:  true, // map_decorations
		50:  true, // potion_contents
		51:  true, // suspicious_stew_effects
		52:  true, // writable_book_content
		53:  true, // written_book_content
		55:  true, // debug_stick_state
		56:  true, // entity_data
		57:  true, // bucket_entity_data
		58:  true, // block_entity_data
		59:  true, // instrument
		60:  true, // provides_trim_material
		62:  true, // jukebox_playable
		63:  true, // provides_banner_patterns
		64:  true, // recipes
		65:  true, // lodestone_tracker
		66:  true, // firework_explosion
		67:  true, // fireworks
		68:  true, // profile
		69:  true, // note_block_sound
		70:  true, // banner_patterns
		75:  true, // bees
		76:  true, // lock
		77:  true, // container_loot
		78:  true, // break_sound
		79:  true, // villager_variant
		80:  true, // wolf_variant
		81:  true, // cat_variant
		82:  true, // axolotl_variant
		83:  true, // frog_variant
		84:  true, // painting_variant
		85:  true, // shulker_variant
		86:  true, // goat_variant
		87:  true, // sniffer_variant
		88:  true, // ghoul_variant
		89:  true, // breeze_variant
		90:  true, // bogged_variant
		93:  true, // buckable
		94:  true, // armor_trim
		95:  true, // equippable_color
		96:  true, // trim_material
		97:  true, // trim_pattern
		98:  true, // compass_color
		99:  true, // map_display_color
		101: true, // banner_pattern
	}

	if nbtTypes[componentType] {
		return SkipNBT(r)
	}

	// 容器组件 (ID 73) 特殊处理
	if componentType == 73 {
		return skipContainerComponentData(r)
	}

	// 其他未知组件，尝试作为 NBT 跳过
	if err := SkipNBT(r); err != nil {
		logx.Debugf("SkipNBT failed for component type %d: %v", componentType, err)
		if r.Len() > 0 {
			logx.Debugf("Component type %d: attempting to skip remaining %d bytes", componentType, r.Len())
			// 作为最后手段，尝试跳过剩余字节
			_, err := r.Seek(int64(r.Len()), 1)
			return err
		}
		return err
	}
	return nil
}

// skipContainerComponentData 跳过容器组件数据
func skipContainerComponentData(r *bytes.Reader) error {
	// 读取容器大小(忽略，仅用于验证)
	_, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}

	// 读取内容数量
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}

	// 跳过所有槽位
	for i := int32(0); i < count; i++ {
		if err := skipSlotComponents(r); err != nil {
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
