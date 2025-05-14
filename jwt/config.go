package jwt

import (
	"os"
	"strconv"
)

type Config struct {
	SecretKey string `yaml:"secretKey"`
	ExpIn     int    `yaml:"expIn"`
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
