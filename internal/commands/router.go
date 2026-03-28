package commands

import (
	"fmt"
	"strings"
	"sync"
)

type Router struct {
	bot      BotAdapter
	commands map[string]Command
	parser   MessageParser
	prefix   string
	auth     *AuthManager

	currentCmd Command
	mu         sync.RWMutex
}

func NewRouter(bot BotAdapter, prefix string) *Router {
	return &Router{
		bot:      bot,
		commands: make(map[string]Command),
		prefix:   prefix,
		auth:     NewAuthManager(),
	}
}

func (r *Router) SetParser(parser MessageParser) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.parser = parser
}

func (r *Router) SetAuth(auth *AuthManager) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.auth = auth
}

func (r *Router) GetAuth() *AuthManager {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.auth
}

func (r *Router) RegisterCommand(cmd Command) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.commands[strings.ToLower(cmd.Name())] = cmd
}

func (r *Router) UnregisterCommand(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.commands, strings.ToLower(name))
}

func (r *Router) GetCommand(name string) (Command, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cmd, ok := r.commands[strings.ToLower(name)]
	return cmd, ok
}

func (r *Router) ListCommands() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.commands))
	for name := range r.commands {
		names = append(names, name)
	}
	return names
}

func (r *Router) HandleRawChat(raw RawChat) {
	if !r.bot.IsOnline() {
		return
	}

	r.mu.RLock()
	parser := r.parser
	auth := r.auth
	r.mu.RUnlock()

	if parser == nil {
		return
	}

	msg := parser.Parse(raw)
	if msg == nil {
		return
	}

	if !strings.HasPrefix(msg.PlainText, r.prefix) {
		return
	}

	cmdName, args := r.parseCommand(msg.PlainText)
	if cmdName == "" {
		return
	}

	if !auth.Check(msg) {
		r.bot.SendPrivateMessage(msg.Sender, "你没有权限使用此机器人")
		return
	}

	cmd, ok := r.GetCommand(cmdName)
	if !ok {
		r.bot.SendPrivateMessage(msg.Sender, fmt.Sprintf("未知指令: %s", cmdName))
		return
	}

	ctx := &ChatContext{
		Bot:     r.bot,
		Message: *msg,
		Sender:  msg.Sender,
		Args:    args,
	}

	result := cmd.Execute(ctx)
	if result != nil && result.Message != "" {
		r.bot.SendPrivateMessage(msg.Sender, result.Message)
	}

	if result != nil && result.NextState != StateIdle {
		r.mu.Lock()
		r.currentCmd = cmd
		r.mu.Unlock()
	}
}

func (r *Router) parseCommand(text string) (string, []string) {
	content := strings.TrimPrefix(text, r.prefix)
	content = strings.TrimSpace(content)

	parts := strings.Fields(content)
	if len(parts) == 0 {
		return "", nil
	}

	cmd := strings.ToLower(parts[0])
	args := parts[1:]

	return cmd, args
}

func (r *Router) Tick() {
	r.mu.RLock()
	cmd := r.currentCmd
	r.mu.RUnlock()

	if cmd == nil {
		return
	}

	ctx := &ChatContext{
		Bot: r.bot,
	}

	result := cmd.Tick(ctx)
	if result == nil {
		return
	}

	if result.Message != "" && cmd.Target() != "" {
		r.bot.SendPrivateMessage(cmd.Target(), result.Message)
	}
}

func (r *Router) CurrentCommand() Command {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.currentCmd
}

func (r *Router) ClearCurrentCommand() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.currentCmd != nil {
		r.currentCmd.Stop()
		r.currentCmd = nil
	}
}

type SimpleHandler func(ctx *ChatContext) *CommandResult

type simpleCommand struct {
	name    string
	handler SimpleHandler
}

func (s *simpleCommand) Name() string                            { return s.name }
func (s *simpleCommand) Description() string                     { return "" }
func (s *simpleCommand) Usage() string                           { return s.name }
func (s *simpleCommand) Execute(ctx *ChatContext) *CommandResult { return s.handler(ctx) }
func (s *simpleCommand) Tick(ctx *ChatContext) *CommandResult    { return nil }
func (s *simpleCommand) Cleanup()                                {}
func (s *simpleCommand) Stop()                                   {}
func (s *simpleCommand) State() StateType                        { return StateIdle }
func (s *simpleCommand) Target() string                          { return "" }

func (r *Router) RegisterSimpleCommand(name string, handler SimpleHandler) {
	r.RegisterCommand(&simpleCommand{name: strings.ToLower(name), handler: handler})
}
