package config

type Config struct {
	Log struct {
		LogDir     string
		MaxSize    int
		Debug      bool
		EnableFile bool
	}
	Account struct {
		PlayerID string
	}
	Server struct {
		Address string
	}
}

// LoadConfig loads configuration from config.yaml
func LoadConfig() (*Config, error) {
	// TODO: Implement config loading
	return &Config{}, nil
}
