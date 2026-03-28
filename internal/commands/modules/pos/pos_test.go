package pos

import (
	"math"
	"testing"

	"gmcc/internal/commands"
)

func TestPosCommand_Name(t *testing.T) {
	cmd := NewPosCommand()
	if cmd.Name() != "pos" {
		t.Errorf("Name() = %q, want %q", cmd.Name(), "pos")
	}
}

func TestPosCommand_Execute(t *testing.T) {
	cmd := NewPosCommand()
	bot := &mockBot{
		position: [3]float64{100.5, 64.0, -200.25},
		rotation: [2]float32{45.0, -30.0},
	}
	cmd.Init(bot, nil)

	ctx := &commands.ChatContext{Bot: bot}
	result := cmd.Execute(ctx)

	if !result.Success {
		t.Errorf("Execute() failed: %s", result.Message)
	}
	if result.Message == "" {
		t.Error("Execute() returned empty message")
	}
}

type mockBot struct {
	position [3]float64
	rotation [2]float32
	online   bool
}

func (m *mockBot) GetPlayerID() string                         { return "MockBot" }
func (m *mockBot) GetUUID() string                             { return "mock-uuid" }
func (m *mockBot) GetPosition() (x, y, z float64)              { return m.position[0], m.position[1], m.position[2] }
func (m *mockBot) GetRotation() (yaw, pitch float32)           { return m.rotation[0], m.rotation[1] }
func (m *mockBot) SendChat(msg string) error                   { return nil }
func (m *mockBot) SendCommand(cmd string) error                { return nil }
func (m *mockBot) SendPrivateMessage(target, msg string) error { return nil }
func (m *mockBot) SetYawPitch(yaw, pitch float32) error        { return nil }
func (m *mockBot) LookAt(x, y, z float64) error                { return nil }
func (m *mockBot) IsOnline() bool                              { return m.online }
func (m *mockBot) GetNearbyPlayers() []commands.PlayerInfo     { return nil }
func (m *mockBot) GetPlayerByName(name string) (commands.PlayerInfo, bool) {
	return commands.PlayerInfo{}, false
}
func (m *mockBot) DistanceTo(x, y, z float64) float64 {
	dx := m.position[0] - x
	dy := m.position[1] - y
	dz := m.position[2] - z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
func (m *mockBot) SetHeldSlot(slot int16) error        { return nil }
func (m *mockBot) InteractEntity(entityID int32) error { return nil }
