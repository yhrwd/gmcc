package commands

import (
	"fmt"
	"strings"
	"testing"
)

type mockBotAdapter struct {
	online      bool
	playerName  string
	messages    []string
	commands    []string
	privateMsgs []struct {
		target, msg string
	}
}

func (m *mockBotAdapter) GetPlayerID() string               { return m.playerName }
func (m *mockBotAdapter) GetUUID() string                   { return "test-uuid" }
func (m *mockBotAdapter) GetPosition() (x, y, z float64)    { return 0, 64, 0 }
func (m *mockBotAdapter) GetRotation() (yaw, pitch float32) { return 0, 0 }
func (m *mockBotAdapter) SendChat(msg string) error         { m.messages = append(m.messages, msg); return nil }
func (m *mockBotAdapter) SendCommand(cmd string) error {
	m.commands = append(m.commands, cmd)
	return nil
}
func (m *mockBotAdapter) SendPrivateMessage(target, msg string) error {
	m.privateMsgs = append(m.privateMsgs, struct{ target, msg string }{target, msg})
	return nil
}
func (m *mockBotAdapter) IsOnline() bool                                 { return m.online }
func (m *mockBotAdapter) GetNearbyPlayers() []PlayerInfo                 { return nil }
func (m *mockBotAdapter) GetPlayerByName(name string) (PlayerInfo, bool) { return PlayerInfo{}, false }
func (m *mockBotAdapter) DistanceTo(x, y, z float64) float64             { return 0 }
func (m *mockBotAdapter) SetYawPitch(yaw, pitch float32) error           { return nil }
func (m *mockBotAdapter) LookAt(x, y, z float64) error                   { return nil }
func (m *mockBotAdapter) SetHeldSlot(slot int16) error                   { return nil }
func (m *mockBotAdapter) InteractEntity(entityID int32) error            { return nil }

type mockCommand struct {
	name          string
	executeResult *CommandResult
	tickResult    *CommandResult
}

func (m *mockCommand) Name() string                            { return m.name }
func (m *mockCommand) Description() string                     { return "mock command" }
func (m *mockCommand) Usage() string                           { return "mock" }
func (m *mockCommand) Execute(ctx *ChatContext) *CommandResult { return m.executeResult }
func (m *mockCommand) Tick(ctx *ChatContext) *CommandResult    { return m.tickResult }
func (m *mockCommand) Cleanup()                                {}
func (m *mockCommand) Stop()                                   {}
func (m *mockCommand) State() StateType                        { return StateIdle }
func (m *mockCommand) Target() string                          { return "" }

func TestRouter_RegisterCommand(t *testing.T) {
	bot := &mockBotAdapter{online: true, playerName: "TestBot"}
	router := NewRouter(bot, "!")

	cmd := &mockCommand{name: "test"}
	router.RegisterCommand(cmd)

	if _, ok := router.GetCommand("test"); !ok {
		t.Error("command should be registered")
	}

	router.UnregisterCommand("test")
	if _, ok := router.GetCommand("test"); ok {
		t.Error("command should be unregistered")
	}
}

func TestRouter_ListCommands(t *testing.T) {
	bot := &mockBotAdapter{online: true, playerName: "TestBot"}
	router := NewRouter(bot, "!")

	router.RegisterCommand(&mockCommand{name: "ride"})
	router.RegisterCommand(&mockCommand{name: "goto"})
	router.RegisterCommand(&mockCommand{name: "follow"})

	cmds := router.ListCommands()
	if len(cmds) != 3 {
		t.Errorf("ListCommands() returned %d commands, expected 3", len(cmds))
	}
}

func TestRouter_HandleRawChat(t *testing.T) {
	bot := &mockBotAdapter{online: true, playerName: "TestBot"}
	router := NewRouter(bot, "!")
	parser := NewDefaultParser("!", "TestBot")
	router.SetParser(parser)

	auth := NewAuthManager()
	auth.SetAllowAll(true)
	router.SetAuth(auth)

	router.RegisterCommand(&mockCommand{name: "test", executeResult: &CommandResult{
		Success: true,
		Message: "test ok",
	}})

	raw := RawChat{
		PlainText: "[Player1 ➥ TestBot] !test",
	}

	router.HandleRawChat(raw)

	if len(bot.privateMsgs) != 1 {
		t.Errorf("expected 1 private message, got %d", len(bot.privateMsgs))
	}
	if bot.privateMsgs[0].msg != "test ok" {
		t.Errorf("unexpected message: %q", bot.privateMsgs[0].msg)
	}
}

func TestRouter_HandleRawChat_NotOnline(t *testing.T) {
	bot := &mockBotAdapter{online: false, playerName: "TestBot"}
	router := NewRouter(bot, "!")

	parser := NewDefaultParser("!", "TestBot")
	router.SetParser(parser)

	router.RegisterCommand(&mockCommand{name: "test"})

	raw := RawChat{
		PlainText: "[Player1 ➥ TestBot] !test",
	}

	router.HandleRawChat(raw)

	if len(bot.privateMsgs) > 0 {
		t.Error("should not send message when offline")
	}
}

func TestRouter_HandleRawChat_AuthFailed(t *testing.T) {
	bot := &mockBotAdapter{online: true, playerName: "TestBot"}
	router := NewRouter(bot, "!")
	parser := NewDefaultParser("!", "TestBot")
	router.SetParser(parser)

	auth := NewAuthManager()
	auth.SetAllowAll(false)
	auth.SetWhitelist([]string{"WhitelistedUser"})
	router.SetAuth(auth)

	router.RegisterCommand(&mockCommand{name: "test"})

	raw := RawChat{
		PlainText: "[Player1 ➥ TestBot] !test",
	}

	router.HandleRawChat(raw)

	if len(bot.privateMsgs) != 1 {
		t.Errorf("expected 1 private message, got %d", len(bot.privateMsgs))
	}
	if bot.privateMsgs[0].msg != "你没有权限使用此机器人" {
		t.Errorf("unexpected message: %q", bot.privateMsgs[0].msg)
	}
}

func TestRouter_HandleRawChat_UnknownCommand(t *testing.T) {
	bot := &mockBotAdapter{online: true, playerName: "TestBot"}
	router := NewRouter(bot, "!")
	parser := NewDefaultParser("!", "TestBot")
	router.SetParser(parser)

	auth := NewAuthManager()
	auth.SetAllowAll(true)
	router.SetAuth(auth)

	raw := RawChat{
		PlainText: "[Player1 ➥ TestBot] !unknown",
	}

	router.HandleRawChat(raw)

	if len(bot.privateMsgs) != 1 {
		t.Errorf("expected 1 private message, got %d", len(bot.privateMsgs))
	}
	if bot.privateMsgs[0].msg != "未知指令: unknown" {
		t.Errorf("unexpected message: %q", bot.privateMsgs[0].msg)
	}
}

func TestRouter_HandleRawChat_NoPrefix(t *testing.T) {
	bot := &mockBotAdapter{online: true, playerName: "TestBot"}
	router := NewRouter(bot, "!")
	parser := NewDefaultParser("!", "TestBot")
	router.SetParser(parser)

	auth := NewAuthManager()
	auth.SetAllowAll(true)
	router.SetAuth(auth)

	router.RegisterCommand(&mockCommand{name: "test"})

	raw := RawChat{
		PlainText: "[Player1 ➥ TestBot] test",
	}

	router.HandleRawChat(raw)

	if len(bot.privateMsgs) > 0 {
		t.Error("should not respond to message without prefix")
	}
}

func TestRouter_HandleRawChat_NotPrivate(t *testing.T) {
	bot := &mockBotAdapter{online: true, playerName: "TestBot"}
	router := NewRouter(bot, "!")
	parser := NewDefaultParser("!", "TestBot")
	router.SetParser(parser)

	auth := NewAuthManager()
	auth.SetAllowAll(true)
	router.SetAuth(auth)

	router.RegisterCommand(&mockCommand{name: "test"})

	raw := RawChat{
		PlainText: "some random chat message",
	}

	router.HandleRawChat(raw)

	if len(bot.privateMsgs) > 0 {
		t.Error("should not respond to non-private message")
	}
}

func TestRouter_RegisterSimpleCommand(t *testing.T) {
	bot := &mockBotAdapter{online: true, playerName: "TestBot"}
	router := NewRouter(bot, "!")

	router.RegisterSimpleCommand("pos", func(ctx *ChatContext) *CommandResult {
		x, y, z := ctx.Bot.GetPosition()
		return &CommandResult{
			Success: true,
			Message: fmt.Sprintf("%.1f, %.1f, %.1f", x, y, z),
		}
	})

	cmd, ok := router.GetCommand("pos")
	if !ok {
		t.Error("simple command should be registered")
	}
	if cmd.Name() != "pos" {
		t.Errorf("Name() = %q, want %q", cmd.Name(), "pos")
	}
}

func TestRouter_SimpleCommand_Execute(t *testing.T) {
	bot := &mockBotAdapter{online: true, playerName: "TestBot"}
	router := NewRouter(bot, "!")
	parser := NewDefaultParser("!", "TestBot")
	router.SetParser(parser)

	auth := NewAuthManager()
	auth.SetAllowAll(true)
	router.SetAuth(auth)

	router.RegisterSimpleCommand("echo", func(ctx *ChatContext) *CommandResult {
		if len(ctx.Args) == 0 {
			return &CommandResult{Success: false, Message: "usage: !echo <text>"}
		}
		return &CommandResult{Success: true, Message: strings.Join(ctx.Args, " ")}
	})

	raw := RawChat{
		PlainText: "[Player1 ➥ TestBot] !echo hello world",
	}

	router.HandleRawChat(raw)

	if len(bot.privateMsgs) != 1 {
		t.Fatalf("expected 1 private message, got %d", len(bot.privateMsgs))
	}
	if bot.privateMsgs[0].msg != "hello world" {
		t.Errorf("unexpected message: %q", bot.privateMsgs[0].msg)
	}
}
