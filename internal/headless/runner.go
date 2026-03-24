package headless

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"gmcc/internal/config"
	"gmcc/internal/logx"
	"gmcc/internal/mcclient"
	"gmcc/internal/mcclient/chat"
)

// Runner 无界面运行器
type Runner struct {
	cfg    *config.Config
	client *mcclient.Client
}

// New 创建无界面运行器
func New(cfg *config.Config) *Runner {
	return &Runner{
		cfg: cfg,
	}
}

// Run 运行无界面模式
func (r *Runner) Run(ctx context.Context) error {
	logx.Infof("gmcc 启动 (无界面模式)")
	logx.Infof("玩家: %s", r.cfg.Account.PlayerID)
	logx.Infof("服务器: %s", r.cfg.Server.Address)
	logx.Infof("正在连接...")

	r.client = mcclient.New(r.cfg)
	r.client.SetChatHandler(func(msg mcclient.ChatMessage) {
		var text string
		if msg.RawJSON != "" {
			comp, err := chat.ParseTextComponent(msg.RawJSON)
			if err == nil {
				text = comp.ToANSI()
			} else {
				text = msg.PlainText
			}
		} else {
			text = msg.PlainText
		}
		if text != "" {
			logx.Infof("[聊天] %s", text)
		}
	})

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- r.client.Run(ctx)
	}()

	for {
		select {
		case <-ctx.Done():
			logx.Infof("收到退出信号，正在断开连接...")
			time.Sleep(500 * time.Millisecond)
			return nil
		case err := <-errCh:
			if err != nil {
				logx.Errorf("运行错误: %v", err)
			}
			return err
		}
	}
}

// SendCommand 发送命令
func (r *Runner) SendCommand(cmd string) error {
	if r.client == nil || !r.client.IsReady() {
		return fmt.Errorf("客户端未就绪")
	}
	return r.client.SendCommand(cmd)
}

// SendMessage 发送消息
func (r *Runner) SendMessage(msg string) error {
	if r.client == nil || !r.client.IsReady() {
		return fmt.Errorf("客户端未就绪")
	}
	return r.client.SendMessage(msg)
}

// IsReady 检查客户端是否就绪
func (r *Runner) IsReady() bool {
	return r.client != nil && r.client.IsReady()
}
