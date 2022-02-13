package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env           string
	Port          string
	RedisHost     string
	RedisPassword string
}

func LoadConfig(env *string) (Config, error) {

	err := godotenv.Load(*env)
	if err != nil {
		log.Printf("Error loading %s file", *env)
	}

	environment := os.Getenv("ENV")

	port := os.Getenv("PORT")
	if port == "" {
		return Config{}, fmt.Errorf("PORT cannot be empty")
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
		Env:           environment,
		Port:          port,
		RedisHost:     redisHost,
		RedisPassword: redisPassword,
	}, nil
}
