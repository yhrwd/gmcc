package mcclient

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gmcc/internal/logx"
	"gmcc/internal/mcclient/packet"
	"gmcc/internal/mcclient/protocol"
	"gmcc/internal/player"
	"gmcc/internal/registry"
)

// DEBUG_MODE: 临时关闭背包解析，将原始包dump到文件
const DEBUG_DUMP_CONTAINER_PACKETS = false

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
	name, err := c.readAnonymousNBTJSON(r)
	if err != nil {
		return err
	}
	screenHandlerId, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	syncId, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return err
	}

	c.Player.SetOpenContainer(&player.ContainerState{
		WindowID:   syncId,
		WindowType: screenHandlerId,
		Open:       true,
	})

	logx.Infof("open_screen: syncId=%d, screenHandlerId=%d, name=%s", syncId, screenHandlerId, name)
	return nil
}

func (c *Client) handleContainerClosePacket(data []byte) error {
	r := bytes.NewReader(data)
	syncId, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return err
	}

	container := c.Player.GetOpenContainer()
	if container != nil && container.WindowID == syncId {
		c.Player.SetOpenContainer(nil)
	}

	logx.Debugf("container_close: syncId=%d", syncId)
	return nil
}

func (c *Client) handleContainerSetDataPacket(data []byte) error {
	r := bytes.NewReader(data)
	propertyId, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	syncId, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	value, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return err
	}

	logx.Debugf("container_set_data: syncId=%d, propertyId=%d, value=%d", syncId, propertyId, value)
	return nil
}

func (c *Client) handleContainerContentPacket(data []byte) error {
	// DEBUG: 完全dump背包内容包到文件，暂时关闭解析
	if DEBUG_DUMP_CONTAINER_PACKETS {
		return c.dumpContainerPacket("container_content", data)
	}

	r := bytes.NewReader(data)
	syncId, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		logx.PacketError("container_content", data, err)
		return fmt.Errorf("container_content: 读取 syncId 失败: %w", err)
	}
	stateID, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		logx.PacketError("container_content", data, err)
		return nil
	}
	numItems, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		logx.PacketError("container_content", data, err)
		return nil
	}

	c.Player.UpdateContainerStateID(stateID)

	logx.Infof("container_content: syncId=%d, stateID=%d, numItems=%d, remaining=%d bytes", syncId, stateID, numItems, r.Len())

	if numItems > 1000 {
		logx.Warnf("container_content: numItems 过大 (%d), 限制为 1000", numItems)
		numItems = 1000
	}

	reg := registry.GetItemRegistry()
	items := make([]*player.SlotData, numItems)
	for i := int32(0); i < numItems; i++ {
		slot, err := packet.ReadSlotData(r)
		if err != nil {
			logx.PacketError("container_content", data, fmt.Errorf("slot %d: %w", i, err))
			return nil
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

	carriedItem, err := packet.ReadSlotData(r)
	if err != nil {
		logx.PacketError("container_content", data, fmt.Errorf("carriedItem: %w", err))
		return nil
	}
	var carried *player.SlotData
	if carriedItem != nil {
		carried = &player.SlotData{ID: carriedItem.ID, Count: carriedItem.Count}
		itemName := reg.IDToName(carriedItem.ID)
		logx.Infof("container_content: carried item: id=%d (%s), count=%d", carriedItem.ID, itemName, carriedItem.Count)
	}

	c.Player.UpdateInventory(syncId, items, carried)
	return nil
}

// dumpContainerPacket 将容器包完整dump到文件
func (c *Client) dumpContainerPacket(packetName string, data []byte) error {
	// 创建dump目录
	logDir := "logs"
	dumpDir := filepath.Join(logDir, "dumps")
	if err := os.MkdirAll(dumpDir, 0755); err != nil {
		logx.Warnf("无法创建dump目录: %v", err)
		return err
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102-150405")
	filename := filepath.Join(dumpDir, fmt.Sprintf("%s_%s.txt", timestamp, packetName))

	// 创建文件
	f, err := os.Create(filename)
	if err != nil {
		logx.Warnf("无法创建dump文件: %v", err)
		return err
	}
	defer f.Close()

	// 写入头部信息
	fmt.Fprintf(f, "# Packet: %s\n", packetName)
	fmt.Fprintf(f, "# Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(f, "# Length: %d bytes\n", len(data))
	fmt.Fprintln(f, "# --- HEX DUMP ---")

	// 写入十六进制数据
	hexData := hex.EncodeToString(data)
	// 每行64个字符
	for i := 0; i < len(hexData); i += 64 {
		end := i + 64
		if end > len(hexData) {
			end = len(hexData)
		}
		fmt.Fprintln(f, hexData[i:end])
	}

	logx.Infof("dumped packet %s to %s (%d bytes)", packetName, filename, len(data))
	return nil
}

func (c *Client) handleContainerSlotPacket(data []byte) error {
	// DEBUG: 完全dump背包槽位包到文件，暂时关闭解析
	if DEBUG_DUMP_CONTAINER_PACKETS {
		return c.dumpContainerPacket("container_slot", data)
	}

	r := bytes.NewReader(data)
	syncId, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		logx.PacketError("container_slot", data, err)
		return nil
	}
	stateID, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		logx.PacketError("container_slot", data, err)
		return nil
	}

	var slot int16
	if err := binary.Read(r, binary.BigEndian, &slot); err != nil {
		logx.PacketError("container_slot", data, err)
		return nil
	}

	c.Player.UpdateContainerStateID(stateID)

	item, err := packet.ReadSlotData(r)
	if err != nil {
		logx.PacketError("container_slot", data, fmt.Errorf("slot %d: %w", slot, err))
		return nil
	}

	reg := registry.GetItemRegistry()
	var slotItem *player.SlotData
	if item != nil {
		slotItem = &player.SlotData{ID: item.ID, Count: item.Count}
		itemName := reg.IDToName(item.ID)
		localizedName := reg.LocalizedName(item.ID)
		logx.Infof("container_slot: syncId=%d, stateID=%d, slot=%d, item_id=%d (%s), name=%s, count=%d", syncId, stateID, slot, item.ID, itemName, localizedName, item.Count)
	} else {
		logx.Debugf("container_slot: syncId=%d, stateID=%d, slot=%d, item=empty", syncId, stateID, slot)
	}

	c.Player.UpdateSlot(syncId, int32(slot), slotItem)
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
