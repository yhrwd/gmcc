package config

// Config 包含程序的配置结构体定义，按照仓库根目录 config.yaml 的字段组织。
type AccountConfig struct {
	PlayerID string `yaml:"player_id"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

type LogConfig struct {
	LogDir     string `yaml:"log_dir"`
	MaxSize    int64  `yaml:"max_size"`
	Debug      bool   `yaml:"debug"`
	EnableFile bool   `yaml:"enable_file"`
}

type Config struct {
	Account AccountConfig `yaml:"account"`
	Server  ServerConfig  `yaml:"server"`
	Log     LogConfig     `yaml:"log"`
}
