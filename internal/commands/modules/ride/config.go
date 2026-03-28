package ride

import "time"

// 骑乘命令配置常量
const (
	// DefaultRangeLimit 默认骑乘距离限制（格）
	// Minecraft 中玩家交互有效距离约为 3-4 格
	// 设置为 3.0 格确保骑乘成功率高
	DefaultRangeLimit = 3.0

	// DefaultTimeout 默认等待目标玩家的超时时间
	// 30 秒足够玩家移动到附近或取消操作
	DefaultTimeout = 30 * time.Second

	// DefaultCooldown 默认命令冷却时间
	// 5 秒冷却防止命令滥用和刷屏
	DefaultCooldown = 5 * time.Second

	// DefaultLookSmoothing 默认视角平滑系数
	// 范围 0-1，值越小越平滑，0.3 提供自然的视角过渡
	DefaultLookSmoothing = 0.3
)

type Config struct {
	// RangeLimit 骑乘距离限制（格）
	// 玩家必须在此距离内才能执行骑乘命令
	RangeLimit float64

	// Timeout 超时时间
	// 从命令发出到骑乘成功的最大等待时间
	Timeout time.Duration

	// Cooldown 冷却时间
	// 骑乘成功后的冷却时间，防止命令滥用
	Cooldown time.Duration

	// LookSmoothing 视角平滑系数 (0-1)
	// 控制视角追踪的平滑程度，越大越快
	LookSmoothing float32
}

func DefaultConfig() *Config {
	return &Config{
		RangeLimit:    DefaultRangeLimit,
		Timeout:       DefaultTimeout,
		Cooldown:      DefaultCooldown,
		LookSmoothing: DefaultLookSmoothing,
	}
}
