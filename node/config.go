package node

// 配置
type Config struct {
	edps        []endpoint
	ConnTimeOut int
}

func LoadConfig(file string) *Config {
	var cfg Config

	return &cfg
}
