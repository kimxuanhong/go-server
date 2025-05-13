package core

import "os"

// Config defines server configuration.
type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Mode     string `yaml:"mode"`
	RootPath string `yaml:"root-path"`
}

func (c *Config) GetAddr() string {
	return c.Host + ":" + c.Port
}

func NewConfig() *Config {
	return &Config{
		Host:     getEnv("SERVER_HOST", "localhost"),
		Port:     getEnv("SERVER_PORT", "8080"),
		Mode:     getEnv("SERVER_MODE", "debug"),
		RootPath: getEnv("SERVER_ROOT_PATH", ""),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
