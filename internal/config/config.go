package config

// Config 包含程序的配置结构体定义，按照仓库根目录 config.yaml 的字段组织。
type AccountConfig struct {
	PlayerID        string `yaml:"player_id"`
	UseOfficialAuth bool   `yaml:"use_official_auth"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

type ActionsConfig struct {
	OnJoinCommands []string `yaml:"on_join_commands"`
	OnJoinMessages []string `yaml:"on_join_messages"`
	DelayMs        int      `yaml:"delay_ms"`
	SignCommands   bool     `yaml:"sign_commands"`
}

type LogConfig struct {
	LogDir     string `yaml:"log_dir"`
	MaxSize    int64  `yaml:"max_size"`
	Debug      bool   `yaml:"debug"`
	EnableFile bool   `yaml:"enable_file"`
}

type RuntimeConfig struct {
	Headless bool `yaml:"headless"` // 无界面模式，只记录日志
}

type PacketConfig struct {
	HandleContainer bool `yaml:"handle_container"` // 是否处理背包数据包
}

type Config struct {
	Account AccountConfig `yaml:"account"`
	Server  ServerConfig  `yaml:"server"`
	Actions ActionsConfig `yaml:"actions"`
	Log     LogConfig     `yaml:"log"`
	Runtime RuntimeConfig `yaml:"runtime"`
	Packets PacketConfig  `yaml:"packets"`
}

// Default 返回默认配置模板。
func Default() Config {
	return Config{
		Account: AccountConfig{
			PlayerID:        "your_player_id_here",
			UseOfficialAuth: false,
		},
		Server: ServerConfig{
			Address: "127.0.0.1:25565",
		},
		Actions: ActionsConfig{
			OnJoinCommands: nil,
			OnJoinMessages: nil,
			DelayMs:        1200,
			SignCommands:   false,
		},
		Log: LogConfig{
			LogDir:     "logs",
			MaxSize:    512, // 单位：KB，转换为字节时 *1024
			Debug:      false,
			EnableFile: true,
		},
		Runtime: RuntimeConfig{
			Headless: false,
		},
		Packets: PacketConfig{
			HandleContainer: true,
		},
	}
}

// MaxSizeInBytes 返回转换为字节的最大日志大小 (KB -> 字节)
func (c *LogConfig) MaxSizeInBytes() int64 {
	return c.MaxSize * 1024
}
