package jwt

import (
	"github.com/spf13/viper"
	"os"
	"strconv"
)

type Config struct {
	SecretKey string `mapstructure:"secretKey" yaml:"secretKey"`
	ExpIn     int    `mapstructure:"expIn" yaml:"expIn"`
}

func DefaultConfig() *Config {
	return &Config{
		SecretKey: getEnv("JWT_SECRET_KEY", "Matkhau@1234Nam"),
		ExpIn:     getEnvAsInt("JWT_EXP_IN", 3600),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsedValue
}

func GetConfig(configs ...*Config) *Config {
	if len(configs) > 0 && configs[0] != nil {
		return configs[0]
	}
	viper.SetDefault("jwt.secretKey", "Matkhau@1234Nam")
	viper.SetDefault("jwt.expIn", 3600)
	return &Config{
		SecretKey: viper.GetString("jwt.secretKey"),
		ExpIn:     viper.GetInt("jwt.expIn"),
	}
}
