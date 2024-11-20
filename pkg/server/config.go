package server

import "os"

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

// LoadFromEnv loads the configuration from environment variables
func (c *Config) LoadFromEnv() {
	if port := os.Getenv("PORT"); port != "" {
		c.Port = port
	}
	if password := os.Getenv("PASSWORD"); password != "" {
		c.Password = password
	}
}
