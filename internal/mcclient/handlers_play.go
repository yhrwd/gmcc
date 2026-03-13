package mcclient

import (
	"bytes"
	"fmt"
	"time"

	"gmcc/internal/logx"
)

func (c *Client) handlePlayPacket(pkt packet) error {
	switch pkt.ID {
	case playClientDisconnect:
		return fmt.Errorf("Play 阶段被服务器断开: %s", rawPreview(pkt.Data))

	case playClientKeepAlive:
		r := bytes.NewReader(pkt.Data)
		id, err := readInt64(r)
		if err != nil {
			return fmt.Errorf("读取 play keep_alive 失败: %w", err)
		}
		return c.conn.WritePacket(playServerKeepAlive, encodeInt64(id))

	case playClientPing:
		r := bytes.NewReader(pkt.Data)
		id, err := readInt32(r)
		if err != nil {
			return fmt.Errorf("读取 play ping 失败: %w", err)
		}
		return c.conn.WritePacket(playServerPong, encodeInt32(id))

	case playClientDeclareCommands:
		return c.handleDeclareCommandsPacket(pkt.Data)

	case playClientProfilelessChat:
		return c.handleProfilelessChatPacket(pkt.Data)

	case playClientPlayerChat:
		return c.handlePlayerChatPacket(pkt.Data)

	case playClientPosition:
		r := bytes.NewReader(pkt.Data)

		teleportID, err := readVarInt(r)
		if err != nil {
			return fmt.Errorf("读取 player_position teleport id 失败: %w", err)
		}

		x, _ := readFloat64FromReader(r)
		y, _ := readFloat64FromReader(r)
		z, _ := readFloat64FromReader(r)

		deltaX, _ := readFloat64FromReader(r)
		deltaY, _ := readFloat64FromReader(r)
		deltaZ, _ := readFloat64FromReader(r)

		yRot, _ := readFloat32(r)
		xRot, _ := readFloat32(r)

		relBits, _ := readInt32(r)

		logx.Debugf("player_position: teleportID=%d, pos=(%.2f,%.2f,%.2f), rot=(%.2f,%.2f), rel=0x%x",
			teleportID, x, y, z, yRot, xRot, relBits)

		c.Player.UpdatePosition(x, y, z, yRot, xRot, int8(relBits))

		_ = deltaX
		_ = deltaY
		_ = deltaZ

		if err := c.conn.WritePacket(playServerAcceptTeleport, encodeVarInt(teleportID)); err != nil {
			return fmt.Errorf("发送 accept_teleportation 失败: %w", err)
		}
		return nil

	case playClientActionBar:
		return c.handleActionBarPacket(pkt.Data)

	case playClientCookieReq:
		r := bytes.NewReader(pkt.Data)
		key, err := readString(r, r)
		if err != nil {
			return fmt.Errorf("读取 play cookie key 失败: %w", err)
		}
		return c.sendCookieResponse(playServerCookieResp, key)

	case playClientPackPush:
		id, err := readPacketUUID(pkt.Data)
		if err != nil {
			return fmt.Errorf("读取 play resource_pack_push 失败: %w", err)
		}
		return c.sendResourcePackResponses(playServerResource, id)

	case playClientLogin:
		if err := c.handlePlayLoginPacket(pkt.Data); err != nil {
			return err
		}
		if !c.inPlay {
			c.inPlay = true
			logx.Infof("已进入服务器, 开始挂机: %s (%s)", c.username, formatUUID(c.uuid))
			if err := c.initSecureChatSession(); err != nil {
				logx.Warnf("初始化 secure chat 会话失败: %v", err)
			}
			c.runOnJoinActions()
		}
		return nil

	case playClientSystemChat:
		return c.handleSystemChatPacket(pkt.Data)

	case playClientPackPop:
		return nil

	case playClientSetHealth:
		return c.handleSetHealthPacket(pkt.Data)

	case playClientSetExperience:
		return c.handleSetExperiencePacket(pkt.Data)

	case playClientSetHeldSlot:
		return c.handleSetHeldSlotPacket(pkt.Data)

	case playClientContainerContent:
		return c.handleContainerContentPacket(pkt.Data)

	case playClientContainerSlot:
		return c.handleContainerSlotPacket(pkt.Data)

	case playClientGameEvent:
		return c.handleGameEventPacket(pkt.Data)

	case playClientEntityData:
		return c.handleEntityDataPacket(pkt.Data)

	default:
		logx.PacketLogf("未处理的 Play 数据包: id=0x%02X (%s) len=%d", pkt.ID, packetName(statePlay, pkt.ID), len(pkt.Data))
		return nil
	}
}

func (c *Client) sendAFKHeartbeatIfNeeded() error {
	if time.Since(c.lastAFKPacket) < 15*time.Second {
		return nil
	}
	c.lastAFKPacket = time.Now()
	return c.conn.WritePacket(playServerMoveStatus, []byte{0x01})
}

func (c *Client) sendCookieResponse(packetID int32, key string) error {
	payload := make([]byte, 0, len(key)+16)
	payload = append(payload, encodeString(key)...)
	payload = append(payload, encodeBool(false)...)
	return c.conn.WritePacket(packetID, payload)
}

func (c *Client) sendClientInformation(packetID int32) error {
	payload := make([]byte, 0, 64)
	payload = append(payload, encodeString("zh_cn")...)
	payload = append(payload, byte(8))
	payload = append(payload, encodeVarInt(0)...)
	payload = append(payload, encodeBool(true)...)
	payload = append(payload, byte(0x7F))
	payload = append(payload, encodeVarInt(1)...)
	payload = append(payload, encodeBool(false)...)
	payload = append(payload, encodeBool(true)...)
	payload = append(payload, encodeVarInt(0)...)
	return c.conn.WritePacket(packetID, payload)
}

func (c *Client) sendResourcePackResponses(packetID int32, id [16]byte) error {
	for _, action := range []int32{resourcePackAccepted, resourcePackDownloaded, resourcePackLoaded} {
		payload := make([]byte, 0, 24)
		payload = append(payload, id[:]...)
		payload = append(payload, encodeVarInt(action)...)
		if err := c.conn.WritePacket(packetID, payload); err != nil {
			return err
		}
	}
	return nil
}

func readPacketUUID(data []byte) ([16]byte, error) {
	r := bytes.NewReader(data)
	return readUUID(r)
}
