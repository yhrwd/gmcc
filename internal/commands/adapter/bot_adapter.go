package adapter

import (
	"fmt"
	"log"
	"math"
	"sync"

	"gmcc/internal/commands"
	"gmcc/internal/mcclient"
	"gmcc/internal/mcclient/packet"
	"gmcc/internal/mcclient/protocol"
)

type ClientAdapter struct {
	client *mcclient.Client

	currentYaw   float32
	currentPitch float32
	mu           sync.RWMutex
}

func NewClientAdapter(client *mcclient.Client) *ClientAdapter {
	return &ClientAdapter{
		client: client,
	}
}

func (c *ClientAdapter) GetPlayerID() string {
	if c.client == nil {
		return ""
	}

	info := c.client.Player.GetInfo()
	if name, ok := info["name"].(string); ok && name != "" {
		return name
	}
	return ""
}

func (c *ClientAdapter) GetUUID() string {
	if c.client == nil {
		return ""
	}

	info := c.client.Player.GetInfo()
	if uuid, ok := info["uuid"].(string); ok && uuid != "" {
		return uuid
	}
	return ""
}

func (c *ClientAdapter) GetPosition() (x, y, z float64) {
	if c.client == nil {
		log.Printf("[WARN] ClientAdapter.GetPosition: client is nil, returning zero position")
		return 0, 0, 0
	}

	info := c.client.Player.GetInfo()
	if pos, ok := info["position"].([]float64); ok && len(pos) >= 3 {
		return pos[0], pos[1], pos[2]
	}
	return 0, 0, 0
}

func (c *ClientAdapter) GetRotation() (yaw, pitch float32) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentYaw, c.currentPitch
}

func (c *ClientAdapter) SendChat(msg string) error {
	if c.client == nil {
		return fmt.Errorf("client not initialized")
	}
	return c.client.SendMessage(msg)
}

func (c *ClientAdapter) SendCommand(cmd string) error {
	if c.client == nil {
		return fmt.Errorf("client not initialized")
	}
	return c.client.SendCommand(cmd)
}

func (c *ClientAdapter) SendPrivateMessage(target, msg string) error {
	if c.client == nil {
		return fmt.Errorf("client not initialized")
	}
	return c.client.SendCommand(fmt.Sprintf("msg %s %s", target, msg))
}

func (c *ClientAdapter) SetYawPitch(yaw, pitch float32) error {
	c.mu.Lock()
	c.currentYaw = yaw
	c.currentPitch = pitch
	c.mu.Unlock()

	return c.sendPlayerRotation(yaw, pitch, true)
}

func (c *ClientAdapter) LookAt(x, y, z float64) error {
	bx, by, bz := c.GetPosition()

	dx := x - bx
	dy := y - by
	dz := z - bz

	yaw := float32(math.Atan2(-dx, dz) * 180 / math.Pi)
	horizDist := math.Sqrt(dx*dx + dz*dz)
	pitch := float32(math.Atan2(-dy, horizDist) * 180 / math.Pi)

	return c.SetYawPitch(yaw, pitch)
}

func (c *ClientAdapter) GetNearbyPlayers() []commands.PlayerInfo {
	if c.client == nil || c.client.NearbyPlayers == nil {
		return nil
	}

	players := c.client.NearbyPlayers.GetNearbyPlayers()
	result := make([]commands.PlayerInfo, 0, len(players))

	for _, p := range players {
		info := commands.PlayerInfo{
			Name:     p.Username,
			UUID:     packet.FormatUUID(p.UUID),
			EntityID: p.ID,
		}
		info.Position.X = p.Position.X
		info.Position.Y = p.Position.Y
		info.Position.Z = p.Position.Z
		result = append(result, info)
	}

	return result
}

func (c *ClientAdapter) GetPlayerByName(name string) (commands.PlayerInfo, bool) {
	if c.client == nil || c.client.NearbyPlayers == nil {
		return commands.PlayerInfo{}, false
	}

	players := c.client.NearbyPlayers.GetNearbyPlayers()
	for _, p := range players {
		if p.Username == name {
			info := commands.PlayerInfo{
				Name:     p.Username,
				UUID:     packet.FormatUUID(p.UUID),
				EntityID: p.ID,
			}
			info.Position.X = p.Position.X
			info.Position.Y = p.Position.Y
			info.Position.Z = p.Position.Z
			return info, true
		}
	}
	return commands.PlayerInfo{}, false
}

func (c *ClientAdapter) DistanceTo(x, y, z float64) float64 {
	px, py, pz := c.GetPosition()
	dx := px - x
	dy := py - y
	dz := pz - z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (c *ClientAdapter) IsOnline() bool {
	return c.client != nil && c.client.IsReady()
}

func (c *ClientAdapter) sendPlayerRotation(yaw, pitch float32, onGround bool) error {
	c.mu.RLock()
	client := c.client
	c.mu.RUnlock()

	if client == nil {
		return fmt.Errorf("client not initialized")
	}

	return client.SendPlayerRotation(yaw, pitch, onGround)
}

func (c *ClientAdapter) SetHeldSlot(slot int16) error {
	c.mu.RLock()
	client := c.client
	c.mu.RUnlock()

	if client == nil {
		return fmt.Errorf("client not initialized")
	}

	return client.SendSetCarriedItem(slot)
}

func (c *ClientAdapter) InteractEntity(entityID int32) error {
	c.mu.RLock()
	client := c.client
	c.mu.RUnlock()

	if client == nil {
		return fmt.Errorf("client not initialized")
	}

	// 先切换到槽位 0 (快捷栏第一个)
	if err := client.SendSetCarriedItem(0); err != nil {
		return fmt.Errorf("切换快捷栏失败: %w", err)
	}

	// 发送交互包 (右键点击实体)
	return client.SendInteract(entityID, protocol.InteractActionInteract, protocol.HandMainHand, false)
}
