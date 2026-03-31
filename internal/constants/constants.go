package constants

import "time"

const (
	// 网络配置
	DialTimeout      = 10 * time.Second
	ReadTimeout      = 20 * time.Second
	AFKCheckInterval = 15 * time.Second

	// 游戏刻配置
	TickInterval = 50 * time.Millisecond // 20Hz = 每秒20次tick

	// 重试/延迟配置
	AuthRetryDelay = 3 * time.Second
	ReconnectDelay = 5 * time.Second

	// Token 配置
	TokenExpirySkew = 30 * time.Second
	ChatSessionTTL  = 2 * time.Hour

	// 缓冲区配置
	MaxPacketSize    = 2 * 1024 * 1024 // 2MB
	SocketBufferSize = 4096
)
