package ride

import (
	"math"
	"testing"
	"time"

	"gmcc/internal/commands"
)

func TestRideCommand_Name(t *testing.T) {
	cmd := NewRideCommand(nil)
	if cmd.Name() != "ride" {
		t.Errorf("Name() = %q, want %q", cmd.Name(), "ride")
	}
}

func TestRideCommand_Description(t *testing.T) {
	cmd := NewRideCommand(nil)
	if cmd.Description() == "" {
		t.Error("Description() returned empty string")
	}
}

func TestRideCommand_Usage(t *testing.T) {
	cmd := NewRideCommand(nil)
	if cmd.Usage() == "" {
		t.Error("Usage() returned empty string")
	}
}

func TestRideCommand_State(t *testing.T) {
	cmd := NewRideCommand(nil)
	if cmd.State() != commands.StateIdle {
		t.Errorf("State() = %v, want %v", cmd.State(), commands.StateIdle)
	}
}

func TestRideCommand_Target(t *testing.T) {
	cmd := NewRideCommand(nil)
	if cmd.Target() != "" {
		t.Errorf("Target() = %q, want empty", cmd.Target())
	}
}

func TestRideCommand_Cleanup(t *testing.T) {
	cmd := NewRideCommand(nil)
	cmd.Cleanup()
	if cmd.State() != commands.StateIdle {
		t.Errorf("State() after Cleanup() = %v, want %v", cmd.State(), commands.StateIdle)
	}
}

func TestRideCommand_Execute_TargetNotFound(t *testing.T) {
	cmd := NewRideCommand(nil)
	bot := &mockBot{online: true}
	cmd.Init(bot, nil)

	ctx := &commands.ChatContext{
		Bot:    bot,
		Sender: "TestPlayer",
		Args:   []string{"NonExistent"},
	}

	result := cmd.Execute(ctx)

	if !result.Success {
		t.Error("Execute() should succeed even when target not found")
	}
	if cmd.State() != commands.StatePreparing {
		t.Errorf("State() = %v, want %v", cmd.State(), commands.StatePreparing)
	}
}

func TestRideCommand_Execute_AlreadyRunning(t *testing.T) {
	cmd := NewRideCommand(nil)
	bot := &mockBot{online: true}
	cmd.Init(bot, nil)

	ctx := &commands.ChatContext{
		Bot:    bot,
		Sender: "TestPlayer",
	}

	cmd.Execute(ctx)
	result := cmd.Execute(ctx)

	if result.Success {
		t.Error("Execute() should fail when already running")
	}
	if result.Message == "" {
		t.Error("Execute() should return message when already running")
	}
}

func TestRideCommand_Execute_TargetInRange(t *testing.T) {
	cmd := &RideCommand{
		config: &Config{RangeLimit: 10.0},
		state:  commands.StateIdle,
	}
	bot := &mockBot{
		online:   true,
		position: [3]float64{0, 64, 0},
		players: []mockPlayer{
			{name: "NearbyPlayer", entityID: 100, position: [3]float64{2, 64, 2}},
		},
	}
	cmd.Init(bot, nil)

	ctx := &commands.ChatContext{
		Bot:    bot,
		Sender: "TestPlayer",
		Args:   []string{"NearbyPlayer"},
	}

	result := cmd.Execute(ctx)

	if !result.Success {
		t.Errorf("Execute() failed: %s", result.Message)
	}
	if cmd.State() != commands.StateCooldown {
		t.Errorf("State() = %v, want %v", cmd.State(), commands.StateCooldown)
	}
}

func TestRideCommand_Execute_TargetOutOfRange(t *testing.T) {
	cmd := &RideCommand{
		config: &Config{RangeLimit: 3.0},
		state:  commands.StateIdle,
	}
	bot := &mockBot{
		online:   true,
		position: [3]float64{0, 64, 0},
		players: []mockPlayer{
			{name: "FarPlayer", entityID: 200, position: [3]float64{50, 64, 50}},
		},
	}
	cmd.Init(bot, nil)

	ctx := &commands.ChatContext{
		Bot:    bot,
		Sender: "TestPlayer",
		Args:   []string{"FarPlayer"},
	}

	result := cmd.Execute(ctx)

	if !result.Success {
		t.Errorf("Execute() failed: %s", result.Message)
	}
	if cmd.State() != commands.StateExecuting {
		t.Errorf("State() = %v, want %v", cmd.State(), commands.StateExecuting)
	}
}

func TestRideCommand_Tick_Timeout(t *testing.T) {
	cmd := &RideCommand{
		config:    &Config{Timeout: 1 * time.Nanosecond}, // Very short timeout
		state:     commands.StatePreparing,
		target:    "NonExistent",
		startTime: time.Now().Add(-1 * time.Hour), // Already expired
	}
	bot := &mockBot{online: true, players: []mockPlayer{}}
	cmd.Init(bot, nil)

	ctx := &commands.ChatContext{Bot: bot}
	result := cmd.Tick(ctx)

	if result == nil {
		t.Error("Tick() should return error result when timed out")
	}
	if result.Success {
		t.Error("Timeout result should not be success")
	}
}

func TestCalculateLookAt(t *testing.T) {
	tests := []struct {
		name           string
		bx, by, bz     float64
		tx, ty, tz     float64
		expectYawDeg   float64
		expectPitchDeg float64
	}{
		{"look forward", 0, 64, 0, 0, 64, 10, 0, 0},
		{"look back", 0, 64, 0, 0, 64, -10, -180, 0},
		{"look down", 0, 64, 0, 0, 60, 0, 0, 90},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yaw, pitch := calculateLookAt(tt.bx, tt.by, tt.bz, tt.tx, tt.ty, tt.tz)
			delta := float32(5.0)
			if diff := abs32(yaw - float32(tt.expectYawDeg)); diff > delta {
				t.Errorf("yaw = %v, want ~%v (delta %v)", yaw, tt.expectYawDeg, delta)
			}
			if diff := abs32(pitch - float32(tt.expectPitchDeg)); diff > delta {
				t.Errorf("pitch = %v, want ~%v (delta %v)", pitch, tt.expectPitchDeg, delta)
			}
		})
	}
}

func TestSmoothLook(t *testing.T) {
	tests := []struct {
		name         string
		currentYaw   float32
		currentPitch float32
		targetYaw    float32
		targetPitch  float32
		smoothing    float32
	}{
		{"no smoothing", 0, 0, 90, 45, 1.0},
		{"half smoothing", 0, 0, 90, 45, 0.5},
		{"min smoothing", 0, 0, 90, 45, 0.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newYaw, newPitch := smoothLook(tt.currentYaw, tt.currentPitch, tt.targetYaw, tt.targetPitch, tt.smoothing)

			if tt.smoothing == 1.0 {
				if newYaw != tt.targetYaw {
					t.Errorf("newYaw = %v, want %v", newYaw, tt.targetYaw)
				}
				if newPitch != tt.targetPitch {
					t.Errorf("newPitch = %v, want %v", newPitch, tt.targetPitch)
				}
			}

			// 检查角度是否在合理范围内 [-180, 180]
			if newYaw < -180 || newYaw > 180 {
				t.Errorf("newYaw %v out of valid range [-180, 180]", newYaw)
			}
		})
	}
}

func TestSmoothLook_CrossBoundary(t *testing.T) {
	tests := []struct {
		name         string
		currentYaw   float32
		targetYaw    float32
		smoothing    float32
		checkShorter bool
	}{
		{"cross 180 clockwise", 170, -170, 1.0, true},
		{"cross 180 counterclockwise", -170, 170, 1.0, true},
		{"within same hemisphere", 10, 30, 0.5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newYaw, _ := smoothLook(tt.currentYaw, 0, tt.targetYaw, 0, tt.smoothing)

			// 检查结果在有效范围内
			if newYaw < -180 || newYaw > 180 {
				t.Errorf("newYaw %v out of valid range [-180, 180]", newYaw)
			}
		})
	}
}

func TestNormalizeAngle(t *testing.T) {
	tests := []struct {
		input    float32
		expected float32
	}{
		{0, 0},
		{180, 180},
		{-180, -180},
		{360, 0},
		{-360, 0},
		{540, 180},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := normalizeAngle(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeAngle(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

var testNow = mustParseTime("2024-01-01T00:00:00Z")

func mustParseTime(s string) (t time.Time) {
	t, _ = time.Parse(time.RFC3339, s)
	return
}

func abs32(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

type mockBot struct {
	online      bool
	position    [3]float64
	rotation    [2]float32
	messages    []string
	commands    []string
	privateMsgs []struct{ target, msg string }
	players     []mockPlayer
}

type mockPlayer struct {
	name     string
	entityID int32
	position [3]float64
}

func (m *mockBot) GetPlayerID() string               { return "MockBot" }
func (m *mockBot) GetUUID() string                   { return "mock-uuid" }
func (m *mockBot) GetPosition() (x, y, z float64)    { return m.position[0], m.position[1], m.position[2] }
func (m *mockBot) GetRotation() (yaw, pitch float32) { return m.rotation[0], m.rotation[1] }
func (m *mockBot) SendChat(msg string) error         { m.messages = append(m.messages, msg); return nil }
func (m *mockBot) SendCommand(cmd string) error      { m.commands = append(m.commands, cmd); return nil }
func (m *mockBot) SendPrivateMessage(target, msg string) error {
	m.privateMsgs = append(m.privateMsgs, struct{ target, msg string }{target, msg})
	return nil
}
func (m *mockBot) SetYawPitch(yaw, pitch float32) error {
	m.rotation = [2]float32{yaw, pitch}
	return nil
}
func (m *mockBot) LookAt(x, y, z float64) error            { return nil }
func (m *mockBot) IsOnline() bool                          { return m.online }
func (m *mockBot) GetNearbyPlayers() []commands.PlayerInfo { return nil }
func (m *mockBot) DistanceTo(x, y, z float64) float64 {
	dx := m.position[0] - x
	dy := m.position[1] - y
	dz := m.position[2] - z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
func (m *mockBot) GetPlayerByName(name string) (commands.PlayerInfo, bool) {
	for _, p := range m.players {
		if p.name == name {
			return commands.PlayerInfo{
				Name:     p.name,
				EntityID: p.entityID,
				Position: struct{ X, Y, Z float64 }{p.position[0], p.position[1], p.position[2]},
			}, true
		}
	}
	return commands.PlayerInfo{}, false
}
func (m *mockBot) SetHeldSlot(slot int16) error        { return nil }
func (m *mockBot) InteractEntity(entityID int32) error { return nil }
