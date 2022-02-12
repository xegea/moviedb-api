package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisHost     string
	RedisPassword string
}

func LoadConfig(env string) (Config, error) {

	err := godotenv.Load(env)
	if err != nil {
		log.Printf("Error loading %s file", env)
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		return Config{}, fmt.Errorf("REDIS_HOST cannot be empty")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		return Config{}, fmt.Errorf("REDIS_PASSWORD cannot be empty")
	}

	return Config{
		RedisHost:     redisHost,
		RedisPassword: redisPassword,
	}, nil
}
