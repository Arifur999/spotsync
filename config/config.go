package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DSN            string
	JWTSecret      string
	AllowedOrigins []string
}

func LoadEnv() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on system environment variables")
	}

	return &Config{
		Port:           os.Getenv("PORT"),
		DSN:            os.Getenv("DSN"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		AllowedOrigins: parseOrigins(os.Getenv("ALLOWED_ORIGINS")),
	}
}

// parseOrigins turns a comma-separated ALLOWED_ORIGINS env value into a
// slice, defaulting to "*" (all origins) when unset.
func parseOrigins(raw string) []string {
	if raw == "" {
		return []string{"*"}
	}

	origins := strings.Split(raw, ",")
	for i, o := range origins {
		origins[i] = strings.TrimSpace(o)
	}

	return origins
}
