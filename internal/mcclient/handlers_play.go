package mcclient

import (
	"bytes"
	"fmt"
	"time"

	"gmcc/internal/constants"
	"gmcc/internal/logx"
	"gmcc/internal/mcclient/packet"
	"gmcc/internal/mcclient/protocol"
)

func (c *Client) handlePlayPacket(pkt packet.Packet) error {
	switch pkt.ID {
	case protocol.PlayClientDisconnect:
		return fmt.Errorf("Play 阶段被服务器断开: %s", disconnectReasonFromNBT(pkt.Data))

	case protocol.PlayClientKeepAlive:
		r := bytes.NewReader(pkt.Data)
		id, err := packet.ReadInt64(r)
		if err != nil {
			return fmt.Errorf("读取 play keep_alive 失败: %w", err)
		}
		return c.conn.WritePacket(protocol.PlayServerKeepAlive, packet.EncodeInt64(id))

	case protocol.PlayClientPing:
		r := bytes.NewReader(pkt.Data)
		id, err := packet.ReadInt32(r)
		if err != nil {
			return fmt.Errorf("读取 play ping 失败: %w", err)
		}
		return c.conn.WritePacket(protocol.PlayServerPong, packet.EncodeInt32(id))

	case protocol.PlayClientDeclareCommands:
		return c.handleDeclareCommandsPacket(pkt.Data)

	case protocol.PlayClientProfilelessChat:
		return c.handleProfilelessChatPacket(pkt.Data)

	case protocol.PlayClientPlayerChat:
		return c.handlePlayerChatPacket(pkt.Data)

	case protocol.PlayClientPosition:
		r := bytes.NewReader(pkt.Data)

		teleportID, err := packet.ReadVarInt(r)
		if err != nil {
			return fmt.Errorf("读取 player_position teleport id 失败: %w", err)
		}

		x := packet.MustReadFloat64(r, "player_position.x")
		y := packet.MustReadFloat64(r, "player_position.y")
		z := packet.MustReadFloat64(r, "player_position.z")

		deltaX := packet.MustReadFloat64(r, "player_position.deltaX")
		deltaY := packet.MustReadFloat64(r, "player_position.deltaY")
		deltaZ := packet.MustReadFloat64(r, "player_position.deltaZ")

		yRotBytes := packet.MustReadBytes(r, 4, "player_position.yRot")
		yRot := packet.ReadFloat32FromBytes(yRotBytes)
		xRotBytes := packet.MustReadBytes(r, 4, "player_position.xRot")
		xRot := packet.ReadFloat32FromBytes(xRotBytes)

		relBits := packet.MustReadVarInt(r, "player_position.relBits")

		// 处理相对坐标 (服务器传送时使用相对坐标)
		px, py, pz := c.Player.GetPosition()
		if relBits&0x01 != 0 {
			x = px + x
		}
		if relBits&0x02 != 0 {
			y = py + y
		}
		if relBits&0x04 != 0 {
			z = pz + z
		}

		logx.Debugf("player_position: teleportID=%d, pos=(%.2f,%.2f,%.2f), rot=(%.2f,%.2f), rel=0x%x, delta=(%.4f,%.4f,%.4f)",
			teleportID, x, y, z, yRot, xRot, relBits, deltaX, deltaY, deltaZ)

		c.Player.UpdatePosition(x, y, z, yRot, xRot, 0) // 重置相对标志，因为我们已处理过了

		if err := c.conn.WritePacket(protocol.PlayServerAcceptTeleport, packet.EncodeVarInt(teleportID)); err != nil {
			return fmt.Errorf("发送 accept_teleportation 失败: %w", err)
		}
		return nil

	case protocol.PlayClientActionBar:
		return c.handleActionBarPacket(pkt.Data)

	case protocol.PlayClientCookieReq:
		r := bytes.NewReader(pkt.Data)
		key, err := packet.ReadString(r, r)
		if err != nil {
			return fmt.Errorf("读取 play cookie key 失败: %w", err)
		}
		return c.sendCookieResponse(protocol.PlayServerCookieResp, key)

	case protocol.PlayClientPackPush:
		id, err := readPacketUUID(pkt.Data)
		if err != nil {
			return fmt.Errorf("读取 play resource_pack_push 失败: %w", err)
		}
		return c.sendResourcePackResponses(protocol.PlayServerResource, id)

	case protocol.PlayClientLogin:
		if err := c.handlePlayLoginPacket(pkt.Data); err != nil {
			return err
		}
		if !c.inPlay {
			c.inPlay = true
			logx.Infof("已进入服务器, 开始挂机: %s (%s)", c.username, packet.FormatUUID(c.uuid))
			if err := c.initSecureChatSession(); err != nil {
				logx.Warnf("初始化 secure chat 会话失败: %v", err)
			}
			c.initializeTrackers()
			c.startTicker()
			c.runOnJoinActions()
		}
		return nil

	case protocol.PlayClientSystemChat:
		return c.handleSystemChatPacket(pkt.Data)

	case protocol.PlayClientPackPop:
		return nil

	case protocol.PlayClientSetHealth:
		return c.handleSetHealthPacket(pkt.Data)

	case protocol.PlayClientSetExperience:
		return c.handleSetExperiencePacket(pkt.Data)

	case protocol.PlayClientSetHeldSlot:
		return c.handleSetHeldSlotPacket(pkt.Data)

	case protocol.PlayClientContainerContent:
		if c.cfg.Packets.HandleContainer {
			return c.handleContainerContentPacket(pkt.Data)
		}
		return nil

	case protocol.PlayClientContainerSlot:
		if c.cfg.Packets.HandleContainer {
			return c.handleContainerSlotPacket(pkt.Data)
		}
		return nil

	case protocol.PlayClientContainerClose:
		if c.cfg.Packets.HandleContainer {
			return c.handleContainerClosePacket(pkt.Data)
		}
		return nil

	case protocol.PlayClientContainerSetData:
		if c.cfg.Packets.HandleContainer {
			return c.handleContainerSetDataPacket(pkt.Data)
		}
		return nil

	case protocol.PlayClientOpenScreen:
		return c.handleOpenScreenPacket(pkt.Data)

	case protocol.PlayClientPlayerAbilities:
		return c.handlePlayerAbilitiesPacket(pkt.Data)

	case protocol.PlayClientGameEvent:
		return c.handleGameEventPacket(pkt.Data)

	case protocol.PlayClientEntityData:
		return c.handleEntityDataPacket(pkt.Data)

	case protocol.PlayClientPlayerInfoUpdate:
		return c.handlePlayerInfoUpdate(pkt.Data)

	case protocol.PlayClientPlayerInfoRemove:
		return c.handlePlayerInfoRemove(pkt.Data)

	// 实体跟踪相关包
	case protocol.PlayClientAddEntity:
		return c.handleAddEntity(pkt.Data)

	case protocol.PlayClientTeleportEntity:
		return c.handleTeleportEntity(pkt.Data)

	case protocol.PlayClientMoveEntityPos:
		return c.handleMoveEntityPos(pkt.Data)

	case protocol.PlayClientRemoveEntities:
		return c.handleRemoveEntities(pkt.Data)

	default:
		logx.PacketLogf("未处理的 Play 数据包: id=0x%02X (%s) len=%d", pkt.ID, protocol.PacketName(protocol.StatePlay, pkt.ID), len(pkt.Data))
		return nil
	}
}

func (c *Client) sendAFKHeartbeatIfNeeded() error {
	if time.Since(c.lastAFKPacket) < constants.AFKCheckInterval {
		return nil
	}
	c.lastAFKPacket = time.Now()

	x, y, z, _, _, onGround := c.Player.GetMovementState()
	flags := byte(0)
	if onGround {
		flags |= 0x01
	}

	payload := make([]byte, 0, 25)
	payload = append(payload, packet.EncodeFloat64(x)...)
	payload = append(payload, packet.EncodeFloat64(y)...)
	payload = append(payload, packet.EncodeFloat64(z)...)
	payload = append(payload, flags)

	if err := c.conn.WritePacket(protocol.PlayServerMovePlayerPos, payload); err != nil {
		return err
	}

	return c.conn.WritePacket(protocol.PlayServerClientTickEnd, nil)
}

func (c *Client) sendCookieResponse(packetID int32, key string) error {
	payload := make([]byte, 0, len(key)+16)
	payload = append(payload, packet.EncodeString(key)...)
	payload = append(payload, packet.EncodeBool(false)...)
	return c.conn.WritePacket(packetID, payload)
}

func (c *Client) sendClientInformation(packetID int32) error {
	payload := make([]byte, 0, 64)
	payload = append(payload, packet.EncodeString("zh_cn")...)
	payload = append(payload, byte(8))
	payload = append(payload, packet.EncodeVarInt(0)...)
	payload = append(payload, packet.EncodeBool(true)...)
	payload = append(payload, byte(0x7F))
	payload = append(payload, packet.EncodeVarInt(1)...)
	payload = append(payload, packet.EncodeBool(false)...)
	payload = append(payload, packet.EncodeBool(true)...)
	payload = append(payload, packet.EncodeVarInt(0)...)
	return c.conn.WritePacket(packetID, payload)
}

func (c *Client) sendResourcePackResponses(packetID int32, id [16]byte) error {
	for _, action := range []int32{protocol.ResourcePackAccepted, protocol.ResourcePackDownloaded, protocol.ResourcePackLoaded} {
		payload := make([]byte, 0, 24)
		payload = append(payload, id[:]...)
		payload = append(payload, packet.EncodeVarInt(action)...)
		if err := c.conn.WritePacket(packetID, payload); err != nil {
			return err
		}
	}
	return nil
}

func readPacketUUID(data []byte) ([16]byte, error) {
	r := bytes.NewReader(data)
	return packet.ReadUUID(r)
}
