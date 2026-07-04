package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DSN       string
	JWTSecret string
}

func LoadEnv() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on system environment variables")
	}

	return &Config{
		Port:      os.Getenv("PORT"),
		DSN:       os.Getenv("DSN"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
