package main

import (
	"os"
	"path/filepath"
	"time"

	"gmcc/internal/auth"
	"gmcc/internal/config"
	"gmcc/internal/protocol"
	"gmcc/internal/tui"
	"gmcc/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	exeDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	logDir := cfg.Log.LogDir
	if !filepath.IsAbs(logDir) {
		logDir = filepath.Join(exeDir, logDir)
	}

	err = logger.InitLogger(
		logDir,
		cfg.Log.MaxSize,    // 单文件最大大小
		cfg.Log.Debug,      // 是否开启 debug
		cfg.Log.EnableFile, // 是否写文件
	)
	if err != nil {
		panic(err)
	}

	go func() {
		time.Sleep(500 * time.Millisecond)

		logger.Info("配置加载完成")
		// logger.Debugf("日志目录: %s", logDir)

		_, err := auth.LoginAccount(cfg.Account.PlayerID)
		if err != nil {
			logger.Errorf("登录错误: %v", err)
			return
		}

		addr, err := protocol.CheckAddr(cfg.Server.Address)
		if err != nil {
			logger.Errorf("服务器地址填写错误: %v", err)
			return
		}
		logger.Infof("服务器地址: %s", addr)
		tui.Push(tui.MsgUpdateStatus{Addr: addr})
		tui.Push(tui.MsgUpdateUserID(cfg.Account.PlayerID))

		path := filepath.Join(
			exeDir,
			".session",
			"e568f2a2e6adc9bba3ddd59c337013c60fcb8692c05cf248cb1cc114f3f96e4f.png",
		)

		tui.Push(tui.MsgUpdateImage(path))

		protocol.Startmc(addr)

		// 5. 如果 Start 返回了，说明连接断了，更新一下 UI 状态
		tui.Push(tui.MsgUpdateStatus{Addr: "已断开"})

	}()

	tui.Start()
}
