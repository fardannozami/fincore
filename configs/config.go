package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl       string
	AutoMigrate bool
	Port        string
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing required env: %s", key)
	}
	return val
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Running without .env (production mode)")
	}

	dbUrl := "host=" + mustGetEnv("DB_HOST") +
		" port=" + mustGetEnv("DB_PORT") +
		" user=" + mustGetEnv("DB_USER") +
		" password=" + mustGetEnv("DB_PASSWORD") +
		" dbname=" + mustGetEnv("DB_NAME") +
		" sslmode=" + getEnv("DB_SSLMODE", "disable")

	autoMigrate, err := strconv.ParseBool(getEnv("AUTO_MIGRATE", "false"))
	if err != nil {
		autoMigrate = false
	}

	port := getEnv("PORT", "8080")

	return &Config{
		DBUrl:       dbUrl,
		AutoMigrate: autoMigrate,
		Port:        port,
	}
}
