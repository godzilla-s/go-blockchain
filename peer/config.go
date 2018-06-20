package peer

type Config struct {
	Addresses   []string // 节点 "127.0.0.1:3001"
	ConnTimeout int64
}

func LoadConfig() *Config {
	var cfg Config

	return &cfg
}
