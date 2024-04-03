package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost   string
	Port         string
	DATABASE_URL string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost:   getEnv("PUBLIC_HOST", "http://localhost"),
		Port:         getEnv("PORT", "8080"),
		DATABASE_URL: getEnv("DATABASE_URL", "panic"),
	}
}

// Gets the env by key or fallbacks
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if fallback == "panic" {
		log.Panic("Env not found: " + key)
	}
	return fallback
}
