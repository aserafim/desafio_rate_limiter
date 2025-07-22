package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RateLimitIP        int
	BlockDurationIP    int
	RateLimitToken     int
	BlockDurationToken int
	RedisAddr          string
	RedisDB            int
	RedisPassword      string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: .env não encontrado, usando variáveis do ambiente.")
	}

	getEnvInt := func(key string, defaultVal int) int {
		valStr := os.Getenv(key)
		val, err := strconv.Atoi(valStr)
		if err != nil {
			return defaultVal
		}
		return val
	}

	return &Config{
		RateLimitIP:        getEnvInt("RATE_LIMIT_IP", 5),
		BlockDurationIP:    getEnvInt("BLOCK_DURATION_IP", 300),
		RateLimitToken:     getEnvInt("RATE_LIMIT_TOKEN", 10),
		BlockDurationToken: getEnvInt("BLOCK_DURATION_TOKEN", 300),
		RedisAddr:          os.Getenv("REDIS_ADDR"),
		RedisPassword:      os.Getenv("REDIS_PASSWORD"),
		RedisDB:            getEnvInt("REDIS_DB", 0),
	}
}
