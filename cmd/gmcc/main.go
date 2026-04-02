package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gmcc/internal/cluster"
	"gmcc/internal/config"
	"gmcc/internal/logx"
	"gmcc/internal/web"
	"gmcc/internal/webtypes"
)

var Version = "dev"

func main() {
	configPath := "config.yaml"
	if v := os.Getenv("GMCC_CONFIG"); v != "" {
		configPath = v
	}

	// 加载统一配置
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[错误] 配置加载失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logx.Init(cfg.Log.LogDir, cfg.Log.EnableFile, cfg.Log.MaxSize*1024, cfg.Log.Debug); err != nil {
		fmt.Fprintf(os.Stderr, "[错误] 日志初始化失败: %v\n", err)
		os.Exit(1)
	}
	defer logx.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 创建集群管理器
	clusterCfg := cluster.ClusterConfig{
		Global: cluster.GlobalConfig{
			MaxInstances: cfg.Cluster.Global.MaxInstances,
			ReconnectPolicy: cluster.ReconnectPolicy{
				Enabled:    cfg.Cluster.Global.ReconnectPolicy.Enabled,
				MaxRetries: cfg.Cluster.Global.ReconnectPolicy.MaxRetries,
				BaseDelay:  cfg.Cluster.Global.ReconnectPolicy.BaseDelay,
				MaxDelay:   cfg.Cluster.Global.ReconnectPolicy.MaxDelay,
				Multiplier: cfg.Cluster.Global.ReconnectPolicy.Multiplier,
			},
		},
		Accounts: convertAccounts(cfg.Cluster.Accounts),
	}

	clusterManager := cluster.NewManager(clusterCfg, configPath)
	if err := clusterManager.Start(); err != nil {
		logx.Errorf("集群管理器启动失败: %v", err)
		os.Exit(1)
	}

	// 创建Web配置
	webCfg := convertToWebConfig(&cfg.Web)

	server, err := web.NewServer(webCfg, configPath, clusterManager)
	if err != nil {
		logx.Errorf("Web服务器创建失败: %v", err)
		os.Exit(1)
	}

	// 启动Web服务器
	go func() {
		if err := server.Run(); err != nil {
			logx.Errorf("Web服务器错误: %v", err)
		}
	}()

	logx.Infof("GMCC Web 面板已启动，版本: %s", Version)
	logx.Infof("配置文件: %s", configPath)

	// 等待退出信号
	<-ctx.Done()
	logx.Infof("正在关闭...")

	if err := clusterManager.Stop(); err != nil {
		logx.Errorf("集群管理器停止错误: %v", err)
	}
}

// convertAccounts 将 config.AccountEntry 转换为 cluster.AccountEntry
func convertAccounts(accounts []config.AccountEntry) []cluster.AccountEntry {
	result := make([]cluster.AccountEntry, len(accounts))
	for i, acc := range accounts {
		result[i] = cluster.AccountEntry{
			ID:              acc.ID,
			PlayerID:        acc.PlayerID,
			UseOfficialAuth: acc.UseOfficialAuth,
			ServerAddress:   acc.ServerAddress,
			Enabled:         acc.Enabled,
		}
	}
	return result
}

// convertToWebConfig 将 config.WebConfig 转换为 webtypes.WebConfig
func convertToWebConfig(webCfg *config.WebConfig) webtypes.WebConfig {
	passwords := make([]webtypes.PasswordEntry, len(webCfg.Auth.Passwords))
	for i, p := range webCfg.Auth.Passwords {
		passwords[i] = webtypes.PasswordEntry{
			ID:        p.ID,
			Hash:      p.Hash,
			Enabled:   p.Enabled,
			CreatedAt: p.CreatedAt,
			Note:      p.Note,
		}
	}

	return webtypes.WebConfig{
		Bind: webCfg.Bind,
		Auth: webtypes.AuthConfig{
			TokenExpiry:           webCfg.Auth.TokenExpiry,
			AuditLogRetentionDays: webCfg.Auth.AuditLogRetentionDays,
			Passwords:             passwords,
		},
		TokenVault: webtypes.TokenVaultConfig{
			StoragePath: webCfg.TokenVault.StoragePath,
			ScryptN:     webCfg.TokenVault.ScryptN,
			ScryptR:     webCfg.TokenVault.ScryptR,
			ScryptP:     webCfg.TokenVault.ScryptP,
			SaltLen:     webCfg.TokenVault.SaltLen,
		},
		CORS: webtypes.CORSConfig{
			Enabled: webCfg.CORS.Enabled,
			Origins: webCfg.CORS.Origins,
		},
	}
}
