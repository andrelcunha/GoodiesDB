package server

type Config struct {
	Port     string
	Password string
}

func NewConfig() *Config {
	return &Config{
		Port:     "6379",
		Password: "guest",
	}
}
