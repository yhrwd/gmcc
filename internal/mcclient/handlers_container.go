package mcclient

import (
	"bytes"
	"encoding/binary"

	"gmcc/internal/logx"
	"gmcc/internal/mcclient/packet"
	"gmcc/internal/mcclient/protocol"
	"gmcc/internal/player"
	"gmcc/internal/registry"
)

const (
	ContainerTypePlayer    int32 = 0
	ContainerTypeChest     int32 = 1
	ContainerTypeCrafting  int32 = 2
	ContainerTypeFurnace   int32 = 3
	ContainerTypeDispenser int32 = 4
	ContainerTypeHopper    int32 = 5
	ContainerTypeAnvil     int32 = 6
	ContainerTypeBeacon    int32 = 7
	ContainerTypeBrewing   int32 = 8
)

var containerTypeNames = map[int32]string{
	ContainerTypePlayer:    "player_inventory",
	ContainerTypeChest:     "chest",
	ContainerTypeCrafting:  "crafting_table",
	ContainerTypeFurnace:   "furnace",
	ContainerTypeDispenser: "dispenser",
	ContainerTypeHopper:    "hopper",
	ContainerTypeAnvil:     "anvil",
	ContainerTypeBeacon:    "beacon",
	ContainerTypeBrewing:   "brewing_stand",
}

func containerTypeName(t int32) string {
	if name, ok := containerTypeNames[t]; ok {
		return name
	}
	return "unknown"
}

func (c *Client) handleOpenScreenPacket(data []byte) error {
	r := bytes.NewReader(data)
	windowID, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	windowType, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	title, err := c.readAnonymousNBTJSON(r)
	if err != nil {
		return err
	}

	c.Player.SetOpenContainer(&player.ContainerState{
		WindowID:   windowID,
		WindowType: windowType,
		Open:       true,
	})

	logx.Infof("open_screen: windowID=%d, type=%s, title=%s", windowID, containerTypeName(windowType), title)
	return nil
}

func (c *Client) handleContainerClosePacket(data []byte) error {
	r := bytes.NewReader(data)
	windowID, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return err
	}

	container := c.Player.GetOpenContainer()
	if container != nil && container.WindowID == windowID {
		c.Player.SetOpenContainer(nil)
	}

	logx.Debugf("container_close: windowID=%d", windowID)
	return nil
}

func (c *Client) handleContainerSetDataPacket(data []byte) error {
	r := bytes.NewReader(data)
	windowID, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	property, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	value, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return err
	}

	logx.Debugf("container_set_data: windowID=%d, property=%d, value=%d", windowID, property, value)
	return nil
}

func (c *Client) handleContainerContentPacket(data []byte) error {
	r := bytes.NewReader(data)
	windowID, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		logx.Warnf("container_content: 读取 windowID 失败: %v", err)
		return nil
	}
	stateID, _ := packet.ReadVarIntFromReader(r)
	numItems, _ := packet.ReadVarIntFromReader(r)

	c.Player.UpdateContainerStateID(stateID)

	logx.Infof("container_content: windowID=%d, stateID=%d, numItems=%d, remaining=%d bytes", windowID, stateID, numItems, r.Len())

	if numItems > 1000 {
		logx.Warnf("container_content: numItems 过大 (%d), 限制为 1000", numItems)
		numItems = 1000
	}

	reg := registry.GetItemRegistry()
	items := make([]*player.SlotData, numItems)
	for i := int32(0); i < numItems; i++ {
		slot, err := packet.ReadSlotData(r)
		if err != nil {
			logx.Warnf("container_content: 读取 slot %d 失败: %v, 剩余 %d 字节", i, err, r.Len())
			break
		}
		if slot != nil {
			items[i] = &player.SlotData{ID: slot.ID, Count: slot.Count}
			if i < 10 || slot.ID != 0 {
				itemName := reg.IDToName(slot.ID)
				localizedName := reg.LocalizedName(slot.ID)
				logx.Debugf("  slot[%d]: id=%d (%s), name=%s, count=%d", i, slot.ID, itemName, localizedName, slot.Count)
			}
		}
	}

	carriedItem, _ := packet.ReadSlotData(r)
	var carried *player.SlotData
	if carriedItem != nil {
		carried = &player.SlotData{ID: carriedItem.ID, Count: carriedItem.Count}
		itemName := reg.IDToName(carriedItem.ID)
		logx.Infof("container_content: carried item: id=%d (%s), count=%d", carriedItem.ID, itemName, carriedItem.Count)
	}

	c.Player.UpdateInventory(windowID, items, carried)
	return nil
}

func (c *Client) handleContainerSlotPacket(data []byte) error {
	r := bytes.NewReader(data)
	windowID, _ := packet.ReadVarIntFromReader(r)
	stateID, _ := packet.ReadVarIntFromReader(r)

	var slot int16
	binary.Read(r, binary.BigEndian, &slot)

	c.Player.UpdateContainerStateID(stateID)

	item, err := packet.ReadSlotData(r)
	if err != nil {
		logx.Warnf("container_slot: 读取物品数据失败: windowID=%d, slot=%d, err=%v", windowID, slot, err)
		return nil
	}

	reg := registry.GetItemRegistry()
	var slotItem *player.SlotData
	if item != nil {
		slotItem = &player.SlotData{ID: item.ID, Count: item.Count}
		itemName := reg.IDToName(item.ID)
		localizedName := reg.LocalizedName(item.ID)
		logx.Infof("container_slot: windowID=%d, stateID=%d, slot=%d, item_id=%d (%s), name=%s, count=%d", windowID, stateID, slot, item.ID, itemName, localizedName, item.Count)
	} else {
		logx.Debugf("container_slot: windowID=%d, stateID=%d, slot=%d, item=empty", windowID, stateID, slot)
	}

	c.Player.UpdateSlot(windowID, int32(slot), slotItem)
	return nil
}

func (c *Client) SendContainerClose(windowID int32) error {
	container := c.Player.GetOpenContainer()
	if container != nil && container.WindowID == windowID {
		c.Player.SetOpenContainer(nil)
	}

	payload := packet.EncodeVarInt(windowID)
	return c.conn.WritePacket(protocol.PlayServerContainerClose, payload)
}

func (c *Client) GetCurrentContainer() *player.ContainerState {
	return c.Player.GetOpenContainer()
}
