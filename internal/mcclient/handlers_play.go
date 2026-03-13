package mcclient

import (
	"bytes"
	"fmt"
	"time"

	"gmcc/internal/logx"
	"gmcc/internal/mcclient/packet"
	"gmcc/internal/mcclient/protocol"
)

func (c *Client) handlePlayPacket(pkt packet.Packet) error {
	switch pkt.ID {
	case protocol.PlayClientDisconnect:
		return fmt.Errorf("Play 阶段被服务器断开: %s", packet.RawPreview(pkt.Data))

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

		x, _ := packet.ReadFloat64FromReader(r)
		y, _ := packet.ReadFloat64FromReader(r)
		z, _ := packet.ReadFloat64FromReader(r)

		deltaX, _ := packet.ReadFloat64FromReader(r)
		deltaY, _ := packet.ReadFloat64FromReader(r)
		deltaZ, _ := packet.ReadFloat64FromReader(r)

		yRot := packet.ReadFloat32FromBytes(packet.ReadBytes(r, 4))
		xRot := packet.ReadFloat32FromBytes(packet.ReadBytes(r, 4))

		relBits, _ := packet.ReadInt32FromReader(r)

		logx.Debugf("player_position: teleportID=%d, pos=(%.2f,%.2f,%.2f), rot=(%.2f,%.2f), rel=0x%x",
			teleportID, x, y, z, yRot, xRot, relBits)

		c.Player.UpdatePosition(x, y, z, yRot, xRot, int8(relBits))

		_ = deltaX
		_ = deltaY
		_ = deltaZ

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
		return c.handleContainerContentPacket(pkt.Data)

	case protocol.PlayClientContainerSlot:
		return c.handleContainerSlotPacket(pkt.Data)

	case protocol.PlayClientGameEvent:
		return c.handleGameEventPacket(pkt.Data)

	case protocol.PlayClientEntityData:
		return c.handleEntityDataPacket(pkt.Data)

	default:
		logx.PacketLogf("未处理的 Play 数据包: id=0x%02X (%s) len=%d", pkt.ID, protocol.PacketName(protocol.StatePlay, pkt.ID), len(pkt.Data))
		return nil
	}
}

func (c *Client) sendAFKHeartbeatIfNeeded() error {
	if time.Since(c.lastAFKPacket) < 15*time.Second {
		return nil
	}
	c.lastAFKPacket = time.Now()
	return c.conn.WritePacket(protocol.PlayServerMoveStatus, []byte{0x01})
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
