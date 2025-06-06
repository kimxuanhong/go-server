package core

import (
	"github.com/spf13/viper"
	"os"
)

// Config defines server configuration.
type Config struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     string `mapstructure:"port" yaml:"port"`
	Mode     string `mapstructure:"mode" yaml:"mode"`
	RootPath string `mapstructure:"root-path" yaml:"root-path"`
	Engine   string `mapstructure:"engine" yaml:"engine"` //gin, fiber, echo
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
		Engine:   getEnv("SERVER_ENGINE", "gin"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetConfig(configs ...*Config) *Config {
	if len(configs) > 0 && configs[0] != nil {
		return configs[0]
	}
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.root-path", "")
	viper.SetDefault("server.engine", "gin")
	return &Config{
		Host:     viper.GetString("server.host"),
		Port:     viper.GetString("server.port"),
		Mode:     viper.GetString("server.mode"),
		RootPath: viper.GetString("server.root-path"),
		Engine:   viper.GetString("server.engine"),
	}
}
