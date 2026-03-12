package mcclient

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"gmcc/internal/logx"
	"gmcc/internal/player"
)

func (c *Client) handleSetHealthPacket(data []byte) error {
	r := bytes.NewReader(data)
	if r.Len() < 9 {
		return nil
	}

	health, _ := readFloat32(r)
	food, _ := readVarInt(r)
	saturation, _ := readFloat32(r)

	c.Player.UpdateHealth(health, 0, int32(food), saturation)
	return nil
}

func (c *Client) handleSetExperiencePacket(data []byte) error {
	r := bytes.NewReader(data)
	if r.Len() < 9 {
		return nil
	}

	experienceBar, _ := readFloat32(r)
	level, _ := readVarInt(r)
	totalExperience, _ := readVarInt(r)

	c.Player.UpdateExperience(int32(level), experienceBar, float32(totalExperience))
	return nil
}

func (c *Client) handleSetHeldSlotPacket(data []byte) error {
	if len(data) < 1 {
		return nil
	}
	slot := data[0]
	c.Player.SetHeldSlot(int8(slot))
	return nil
}

func (c *Client) handleContainerContentPacket(data []byte) error {
	r := bytes.NewReader(data)

	if r.Len() < 3 {
		return nil
	}

	windowID, _ := readU8(r)
	stateID, _ := readVarInt(r)
	count, _ := readVarInt(r)

	if count > 1000 {
		count = 1000
	}

	for i := int32(0); i < count; i++ {
		item, _ := readSlot(r)
		if item != nil {
			pItem := &player.Item{ID: fmt.Sprintf("%d", item.ID), Count: item.Count}
			c.Player.UpdateInventorySlot(int8(windowID), stateID, int8(i), pItem)
		} else {
			c.Player.UpdateInventorySlot(int8(windowID), stateID, int8(i), nil)
		}
	}

	// Skip cursor item
	readSlot(r)

	return nil
}

func (c *Client) handleContainerSlotPacket(data []byte) error {
	r := bytes.NewReader(data)
	if r.Len() < 3 {
		return nil
	}

	windowID, _ := readU8(r)
	stateID, _ := readVarInt(r)
	slot, _ := readU8(r)

	item, _ := readSlot(r)

	if item != nil {
		pItem := &player.Item{ID: fmt.Sprintf("%d", item.ID), Count: item.Count}
		c.Player.UpdateInventorySlot(int8(windowID), stateID, int8(slot), pItem)
	} else {
		c.Player.UpdateInventorySlot(int8(windowID), stateID, int8(slot), nil)
	}

	return nil
}

func (c *Client) handleGameEventPacket(data []byte) error {
	r := bytes.NewReader(data)
	if r.Len() < 5 {
		return nil
	}

	eventType, _ := readU8(r)
	value, _ := readFloat32(r)

	if eventType == 3 && value >= 0 && value <= 3 {
		mode := player.GameMode(int(value))
		c.Player.SetGameMode(mode)
		logx.Infof("游戏模式变更: %s", mode.String())
	}

	return nil
}

func readFloat32(r *bytes.Reader) (float32, error) {
	var v float32
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

func readSlot(r *bytes.Reader) (*Slot, error) {
	hasItem, err := readBool(r)
	if err != nil {
		return nil, err
	}
	if !hasItem {
		return nil, nil
	}

	if r.Len() < 3 {
		return nil, nil
	}

	itemID, err := readVarInt(r)
	if err != nil {
		return nil, nil
	}

	if r.Len() < 1 {
		return nil, nil
	}

	count, err := readU8(r)
	if err != nil {
		return nil, nil
	}

	// Skip data component mask
	if r.Len() < 1 {
		return nil, nil
	}
	mask, err := readU8(r)
	if err != nil {
		return nil, nil
	}

	// Skip components
	if mask&0x01 != 0 {
		compCount, err := readVarInt(r)
		if err != nil {
			return nil, nil
		}
		for i := int32(0); i < compCount; i++ {
			_, err := readVarInt(r) // component type
			if err != nil {
				return nil, nil
			}
			// Skip component data
			_, _ = readVarInt(r) // data length
		}
	}

	// Skip removed components
	if mask&0x02 != 0 {
		removedCount, err := readVarInt(r)
		if err != nil {
			return nil, nil
		}
		for i := int32(0); i < removedCount; i++ {
			_, _ = readVarInt(r)
		}
	}

	return &Slot{ID: itemID, Count: int32(count)}, nil
}

type Slot struct {
	ID    int32
	Count int32
}
