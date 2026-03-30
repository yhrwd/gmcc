package mcclient

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"gmcc/internal/entity"
	"gmcc/internal/logx"
	"gmcc/internal/mcclient/packet"
	"gmcc/internal/mcclient/protocol"
)

// handleAddEntity 处理实体生成包 (0x01)
func (c *Client) handleAddEntity(data []byte) error {
	r := bytes.NewReader(data)

	// 读取实体ID
	entityID, err := packet.ReadVarInt(r)
	if err != nil {
		return fmt.Errorf("读取实体ID失败: %w", err)
	}

	// 读取UUID
	uuid, err := packet.ReadUUID(r)
	if err != nil {
		return fmt.Errorf("读取实体UUID失败: %w", err)
	}

	// 读取实体类型 (注册表ID)
	entityTypeID, err := packet.ReadVarInt(r)
	if err != nil {
		return fmt.Errorf("读取实体类型失败: %w", err)
	}

	// 读取位置
	x, err := readFloat64(r)
	if err != nil {
		return fmt.Errorf("读取X坐标失败: %w", err)
	}
	y, err := readFloat64(r)
	if err != nil {
		return fmt.Errorf("读取Y坐标失败: %w", err)
	}
	z, err := readFloat64(r)
	if err != nil {
		return fmt.Errorf("读取Z坐标失败: %w", err)
	}

	// 读取低精度速度向量
	velocity, err := readLpVec3(r)
	if err != nil {
		return fmt.Errorf("读取实体速度失败: %w", err)
	}

	// 读取角度
	if _, err := packet.ReadU8(r); err != nil {
		return fmt.Errorf("读取pitch失败: %w", err)
	}
	if _, err := packet.ReadU8(r); err != nil {
		return fmt.Errorf("读取yaw失败: %w", err)
	}
	if _, err := packet.ReadU8(r); err != nil {
		return fmt.Errorf("读取head_yaw失败: %w", err)
	}

	// 读取实体数据
	if _, err := packet.ReadVarInt(r); err != nil {
		return fmt.Errorf("读取实体数据失败: %w", err)
	}

	// 转换为实体类型字符串 (简化处理，实际应该查注册表)
	entityType := fmt.Sprintf("minecraft:entity_%d", entityTypeID)

	// 检查是否为玩家 (玩家类型的注册表ID通常是特定的)
	// 这里简化处理，实际应该在注册表中查找
	if entityTypeID == 0 { // 假设0是玩家类型，实际需要查注册表
		entityType = "minecraft:player"
	}

	pos := entity.Position{X: x, Y: y, Z: z}

	if c.entityTracker != nil {
		c.entityTracker.SpawnEntity(int32(entityID), entityType, uuid, pos, velocity)
	}

	logx.Debugf("实体生成: ID=%d, Type=%s, Pos=(%.2f, %.2f, %.2f)", entityID, entityType, x, y, z)

	return nil
}

// handleTeleportEntity 处理实体传送包 (0x48)
func (c *Client) handleTeleportEntity(data []byte) error {
	r := bytes.NewReader(data)

	// 读取实体ID
	entityID, err := packet.ReadVarInt(r)
	if err != nil {
		return fmt.Errorf("读取实体ID失败: %w", err)
	}

	// 读取新位置
	x, err := readFloat64(r)
	if err != nil {
		return fmt.Errorf("读取X坐标失败: %w", err)
	}
	y, err := readFloat64(r)
	if err != nil {
		return fmt.Errorf("读取Y坐标失败: %w", err)
	}
	z, err := readFloat64(r)
	if err != nil {
		return fmt.Errorf("读取Z坐标失败: %w", err)
	}

	// 读取速度
	if _, err := readFloat64(r); err != nil {
		return fmt.Errorf("读取速度X失败: %w", err)
	}
	if _, err := readFloat64(r); err != nil {
		return fmt.Errorf("读取速度Y失败: %w", err)
	}
	if _, err := readFloat64(r); err != nil {
		return fmt.Errorf("读取速度Z失败: %w", err)
	}

	// 读取旋转
	if _, err := readFloat32(r); err != nil {
		return fmt.Errorf("读取yaw失败: %w", err)
	}
	if _, err := readFloat32(r); err != nil {
		return fmt.Errorf("读取pitch失败: %w", err)
	}

	// 读取onGround
	if _, err := packet.ReadBool(r); err != nil {
		return fmt.Errorf("读取onGround失败: %w", err)
	}

	if c.entityTracker != nil {
		newPos := entity.Position{X: x, Y: y, Z: z}
		c.entityTracker.UpdatePosition(int32(entityID), newPos)
	}

	return nil
}

// handleMoveEntityPos 处理实体位置增量更新包 (0x09)
func (c *Client) handleMoveEntityPos(data []byte) error {
	r := bytes.NewReader(data)

	// 读取实体ID
	entityID, err := packet.ReadVarInt(r)
	if err != nil {
		return fmt.Errorf("读取实体ID失败: %w", err)
	}

	// 读取增量 (short类型，需要除以4096)
	var deltaX, deltaY, deltaZ int16
	if err := binary.Read(r, binary.BigEndian, &deltaX); err != nil {
		return fmt.Errorf("读取deltaX失败: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &deltaY); err != nil {
		return fmt.Errorf("读取deltaY失败: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &deltaZ); err != nil {
		return fmt.Errorf("读取deltaZ失败: %w", err)
	}

	// 读取onGround
	if _, err := packet.ReadBool(r); err != nil {
		return fmt.Errorf("读取onGround失败: %w", err)
	}

	if c.entityTracker != nil {
		c.entityTracker.UpdatePositionDelta(int32(entityID), deltaX, deltaY, deltaZ)
	}

	return nil
}

// handleRemoveEntities 处理实体移除包 (0x4B)
func (c *Client) handleRemoveEntities(data []byte) error {
	ids, err := decodeRemoveEntityIDs(data)
	if err != nil {
		return err
	}

	if c.entityTracker != nil {
		c.entityTracker.RemoveEntities(ids)
	}

	logx.Debugf("移除 %d 个实体", len(ids))

	return nil
}

func decodeRemoveEntityIDs(data []byte) ([]int32, error) {
	if ids, err := decodeRemoveEntityIDsByteArray(data); err == nil {
		return ids, nil
	}
	return decodeRemoveEntityIDsCounted(data)
}

func decodeRemoveEntityIDsByteArray(data []byte) ([]int32, error) {
	r := bytes.NewReader(data)
	payload, err := packet.ReadByteArray(r, r)
	if err != nil {
		return nil, fmt.Errorf("读取实体ID字节数组失败: %w", err)
	}
	if r.Len() != 0 {
		return nil, fmt.Errorf("实体ID字节数组后仍有 %d 字节未读取", r.Len())
	}
	return decodeRemoveEntityIDsPayload(payload)
}

func decodeRemoveEntityIDsPayload(data []byte) ([]int32, error) {
	r := bytes.NewReader(data)
	ids := make([]int32, 0)
	for r.Len() > 0 {
		id, err := packet.ReadVarInt(r)
		if err != nil {
			return nil, fmt.Errorf("读取实体ID失败: %w", err)
		}
		ids = append(ids, int32(id))
	}
	return ids, nil
}

func decodeRemoveEntityIDsCounted(data []byte) ([]int32, error) {
	r := bytes.NewReader(data)
	count, err := packet.ReadVarInt(r)
	if err != nil {
		return nil, fmt.Errorf("读取实体数量失败: %w", err)
	}

	ids := make([]int32, 0, count)
	for i := int32(0); i < count; i++ {
		id, err := packet.ReadVarInt(r)
		if err != nil {
			return nil, fmt.Errorf("读取实体ID失败: %w", err)
		}
		ids = append(ids, int32(id))
	}

	if r.Len() != 0 {
		return nil, fmt.Errorf("实体ID列表后仍有 %d 字节未读取", r.Len())
	}
	return ids, nil
}

const (
	lpDataBits    = 15
	lpDataMask    = (1 << lpDataBits) - 1
	lpMaxQuantize = 32766.0
	lpScaleMask   = 3
	lpContFlag    = 4
)

// readLpVec3 读取实体速度使用的低精度三维向量。
// 协议格式与 Minecraft 的 net.minecraft.network.LpVec3 一致。
func readLpVec3(r io.Reader) (entity.Vector3, error) {
	lowest, err := packet.ReadU8(r)
	if err != nil {
		return entity.Vector3{}, fmt.Errorf("读取最低字节失败: %w", err)
	}

	if lowest == 0 {
		return entity.Vector3{}, nil
	}

	middle, err := packet.ReadU8(r)
	if err != nil {
		return entity.Vector3{}, fmt.Errorf("读取中间字节失败: %w", err)
	}

	rest, err := packet.ReadBytes(r, 4)
	if err != nil {
		return entity.Vector3{}, fmt.Errorf("读取高位字节失败: %w", err)
	}

	highest := uint64(binary.BigEndian.Uint32(rest))
	buffer := highest<<16 | uint64(middle)<<8 | uint64(lowest)

	scale := uint64(lowest & lpScaleMask)
	if lowest&lpContFlag != 0 {
		extra, err := packet.ReadVarIntFromReader(r)
		if err != nil {
			return entity.Vector3{}, fmt.Errorf("读取速度缩放扩展失败: %w", err)
		}
		scale |= uint64(uint32(extra)) << 2
	}

	return entity.Vector3{
		X: lpUnpack(buffer>>3) * float64(scale),
		Y: lpUnpack(buffer>>18) * float64(scale),
		Z: lpUnpack(buffer>>33) * float64(scale),
	}, nil
}

func lpUnpack(value uint64) float64 {
	raw := float64(value & lpDataMask)
	return math.Min(raw, lpMaxQuantize)*2.0/lpMaxQuantize - 1.0
}

// readFloat64 从reader读取float64
func readFloat64(r *bytes.Reader) (float64, error) {
	var v float64
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

// readFloat32 从reader读取float32
func readFloat32(r *bytes.Reader) (float32, error) {
	var v float32
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

// init 注册包处理器
func init() {
	// 在Client初始化时调用 registerEntityHandlers
}

// registerEntityHandlers 注册实体相关的包处理器
func (c *Client) registerEntityHandlers() {
	// 这个函数将在初始化Trackers时调用
	// 这里不做任何事情，因为 handler 已经在各自的函数中实现了
	_ = protocol.PlayClientAddEntity
}
