package mcclient

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"gmcc/internal/logx"
	"gmcc/internal/player"
)

func readFloat32(r *bytes.Reader) (float32, error) {
	var v float32
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

func (c *Client) handleSetHealthPacket(data []byte) error {
	r := bytes.NewReader(data)
	health, err := readFloat32(r)
	if err != nil {
		return fmt.Errorf("读取 health 失败: %w", err)
	}
	food, err := readVarInt(r)
	if err != nil {
		return fmt.Errorf("读取 food 失败: %w", err)
	}
	saturation, err := readFloat32(r)
	if err != nil {
		return fmt.Errorf("读取 saturation 失败: %w", err)
	}

	c.Player.UpdateHealth(health, 0, int32(food), saturation)
	logx.Debugf("生命值更新: health=%.1f food=%d saturation=%.1f", health, food, saturation)
	return nil
}

func (c *Client) handleSetExperiencePacket(data []byte) error {
	r := bytes.NewReader(data)
	experienceBar, err := readFloat32(r)
	if err != nil {
		return fmt.Errorf("读取 experience_bar 失败: %w", err)
	}
	level, err := readVarInt(r)
	if err != nil {
		return fmt.Errorf("读取 level 失败: %w", err)
	}
	totalExperience, err := readVarInt(r)
	if err != nil {
		return fmt.Errorf("读取 total_experience 失败: %w", err)
	}

	c.Player.UpdateExperience(int32(level), experienceBar, float32(totalExperience))
	logx.Debugf("经验值更新: level=%d exp=%.1f total=%d", level, experienceBar, totalExperience)
	return nil
}

func (c *Client) handleSetHeldSlotPacket(data []byte) error {
	r := bytes.NewReader(data)
	slot, err := readU8(r)
	if err != nil {
		return fmt.Errorf("读取 held slot 失败: %w", err)
	}

	c.Player.SetHeldSlot(int8(slot))
	logx.Debugf("手持槽位更新: slot=%d", slot)
	return nil
}

func (c *Client) handleContainerContentPacket(data []byte) error {
	r := bytes.NewReader(data)
	windowID, err := readU8(r)
	if err != nil {
		return fmt.Errorf("读取 window_id 失败: %w", err)
	}
	stateID, err := readVarInt(r)
	if err != nil {
		return fmt.Errorf("读取 state_id 失败: %w", err)
	}

	count, err := readVarInt(r)
	if err != nil {
		return fmt.Errorf("读取 slot count 失败: %w", err)
	}

	for i := int32(0); i < count; i++ {
		item, err := readSlot(r)
		if err != nil {
			return fmt.Errorf("读取 slot[%d] 失败: %w", i, err)
		}
		var pItem *player.Item
		if item != nil {
			pItem = &player.Item{ID: item.ID, Count: item.Count}
		}
		c.Player.UpdateInventorySlot(int8(windowID), stateID, int8(i), pItem)
	}

	_, err = readSlot(r)
	if err != nil {
		return fmt.Errorf("读取 cursor item 失败: %w", err)
	}

	logx.Debugf("容器内容更新: window=%d state=%d slots=%d", windowID, stateID, count)
	return nil
}

func (c *Client) handleContainerSlotPacket(data []byte) error {
	r := bytes.NewReader(data)
	windowID, err := readU8(r)
	if err != nil {
		return fmt.Errorf("读取 window_id 失败: %w", err)
	}
	stateID, err := readVarInt(r)
	if err != nil {
		return fmt.Errorf("读取 state_id 失败: %w", err)
	}
	slot, err := readU8(r)
	if err != nil {
		return fmt.Errorf("读取 slot 失败: %w", err)
	}

	item, err := readSlot(r)
	if err != nil {
		return fmt.Errorf("读取 item 失败: %w", err)
	}

	var pItem *player.Item
	if item != nil {
		pItem = &player.Item{ID: item.ID, Count: item.Count}
	}
	c.Player.UpdateInventorySlot(int8(windowID), stateID, int8(slot), pItem)
	logx.Debugf("容器槽位更新: window=%d slot=%d", windowID, slot)
	return nil
}

func (c *Client) handleGameEventPacket(data []byte) error {
	r := bytes.NewReader(data)
	eventType, err := readU8(r)
	if err != nil {
		return fmt.Errorf("读取 game_event type 失败: %w", err)
	}
	value, err := readFloat32(r)
	if err != nil {
		return fmt.Errorf("读取 game_event value 失败: %w", err)
	}

	switch eventType {
	case 3:
		var mode player.GameMode
		if value >= 0 && value <= 3 {
			mode = player.GameMode(int(value))
		}
		c.Player.SetGameMode(mode)
		logx.Infof("游戏模式变更: %s", mode.String())
	case 7, 8:
		logx.Debugf("降雨强度变更: type=%d value=%.1f", eventType, value)
	}

	return nil
}

func readSlot(r *bytes.Reader) (*Slot, error) {
	hasItem, err := readBool(r)
	if err != nil {
		return nil, err
	}
	if !hasItem {
		return nil, nil
	}

	item := &Slot{}

	itemID, err := readVarInt(r)
	if err != nil {
		return nil, err
	}
	item.ID = fmt.Sprintf("minecraft:%d", itemID)

	count, err := readU8(r)
	if err != nil {
		return nil, err
	}
	item.Count = int32(count)

	mask, err := readU8(r)
	if err != nil {
		return nil, err
	}

	hasComponents := mask&0x01 != 0
	hasRemoved := mask&0x02 != 0

	if hasComponents {
		compCount, err := readVarInt(r)
		if err != nil {
			return nil, err
		}
		for i := int32(0); i < compCount; i++ {
			compType, err := readVarInt(r)
			if err != nil {
				return nil, err
			}
			if err := skipComponent(r, compType); err != nil {
				return nil, err
			}
		}
	}

	if hasRemoved {
		removedCount, err := readVarInt(r)
		if err != nil {
			return nil, err
		}
		for i := int32(0); i < removedCount; i++ {
			if _, err := readVarInt(r); err != nil {
				return nil, err
			}
		}
	}

	return item, nil
}

func skipComponent(r *bytes.Reader, compType int32) error {
	switch compType {
	case 1, 2, 3:
		for i := 0; i < 2; i++ {
			flags, err := readU8(r)
			if err != nil {
				return err
			}
			if flags&0x01 != 0 {
				if err := discardBytes(r, 4); err != nil {
					return err
				}
			}
			if flags&0x02 != 0 {
				var size int
				if compType == 1 || compType == 3 {
					size = 4
				} else {
					size = 8
				}
				if err := discardBytes(r, size); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func discardBytes(r *bytes.Reader, n int) error {
	for i := 0; i < n; i++ {
		if _, err := r.ReadByte(); err != nil {
			return err
		}
	}
	return nil
}

type Slot struct {
	ID    string
	Count int32
}
