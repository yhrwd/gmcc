package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	authsession "gmcc/internal/auth/session"
	authvault "gmcc/internal/auth/vault"
	"gmcc/internal/cluster"
	"gmcc/internal/config"
	"gmcc/internal/logx"
	"gmcc/internal/resource"
	"gmcc/internal/state"
	"gmcc/internal/web"
	"gmcc/internal/webtypes"
)

var Version = "dev"

type runtimeDeps struct {
	VaultRepository    *authvault.Repository
	AccountRepository  *state.AccountRepository
	InstanceRepository *state.InstanceRepository
	AuthManager        *authsession.AuthManager
	ResourceManager    *resource.Manager
	ClusterManager     *cluster.Manager
	WebConfig          webtypes.WebConfig
}

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

	runtime, err := buildRuntime(configPath, cfg)
	if err != nil {
		logx.Errorf("运行时初始化失败: %v", err)
		os.Exit(1)
	}

	if err := runtime.ClusterManager.Start(); err != nil {
		logx.Errorf("集群管理器启动失败: %v", err)
		os.Exit(1)
	}

	server, err := web.NewServer(runtime.WebConfig, configPath, runtime.ClusterManager, runtime.ResourceManager, runtime.AuthManager)
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

	if err := runtime.ClusterManager.Stop(); err != nil {
		logx.Errorf("集群管理器停止错误: %v", err)
	}
}

func buildRuntime(configPath string, cfg *config.Config) (*runtimeDeps, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	vaultKey := os.Getenv(cfg.Auth.Vault.KeyEnv)
	if vaultKey == "" {
		return nil, fmt.Errorf("missing vault key env %s", cfg.Auth.Vault.KeyEnv)
	}

	baseDir := runtimeBaseDir(configPath)
	vaultRepo, err := authvault.NewRepository(authvault.Config{
		Dir:       resolveRuntimePath(baseDir, cfg.Auth.Vault.Path),
		MasterKey: []byte(vaultKey),
		ScryptN:   cfg.Auth.Vault.ScryptN,
		ScryptR:   cfg.Auth.Vault.ScryptR,
		ScryptP:   cfg.Auth.Vault.ScryptP,
		SaltLen:   cfg.Auth.Vault.SaltLen,
	})
	if err != nil {
		return nil, fmt.Errorf("create auth vault repository: %w", err)
	}

	accountRepo := state.NewAccountRepository(filepath.Join(baseDir, ".state", "accounts.yaml"))
	instanceRepo := state.NewInstanceRepository(filepath.Join(baseDir, ".state", "instances.yaml"))
	authManager := authsession.NewAuthManager(vaultRepo, authsession.NewLiveProviderSet())
	resourceManager := resource.NewManager(accountRepo, instanceRepo, authManager)

	if restored, err := resourceManager.RestoreResources(); err != nil {
		return nil, fmt.Errorf("restore resource metadata: %w", err)
	} else {
		logx.Infof("资源元数据已加载: restored=%d skipped=%d", restored.RestoredCount, restored.SkippedCount)
	}

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

	clusterManager := cluster.NewManager(clusterCfg, authManager, configPath)
	clusterManager.SetResourceManager(resourceManager)

	return &runtimeDeps{
		VaultRepository:    vaultRepo,
		AccountRepository:  accountRepo,
		InstanceRepository: instanceRepo,
		AuthManager:        authManager,
		ResourceManager:    resourceManager,
		ClusterManager:     clusterManager,
		WebConfig:          convertToWebConfig(&cfg.Web),
	}, nil
}

func runtimeBaseDir(configPath string) string {
	if configPath == "" {
		return "."
	}
	return filepath.Dir(configPath)
}

func resolveRuntimePath(baseDir, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(baseDir, path)
}

// convertAccounts 将 config.AccountEntry 转换为 cluster.AccountEntry
func convertAccounts(accounts []config.AccountEntry) []cluster.AccountEntry {
	result := make([]cluster.AccountEntry, len(accounts))
	for i, acc := range accounts {
		result[i] = cluster.AccountEntry{
			ID:            acc.ID,
			ServerAddress: acc.ServerAddress,
			Enabled:       acc.Enabled,
		}
	}
	return result
}

// convertToWebConfig 将 config.WebConfig 转换为 webtypes.WebConfig
func convertToWebConfig(webCfg *config.WebConfig) webtypes.WebConfig {
	return webtypes.WebConfig{
		Bind: webCfg.Bind,
		Auth: webtypes.AuthConfig{
			AuditLogRetentionDays: webCfg.Auth.AuditLogRetentionDays,
		},
		CORS: webtypes.CORSConfig{
			Enabled: webCfg.CORS.Enabled,
			Origins: webCfg.CORS.Origins,
		},
	}
}
