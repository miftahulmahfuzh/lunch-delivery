// internal/config/config.go
package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// LLM Configuration
	LLMType                string
	DeepseekTencentAPIKey  string
	DeepseekTencentModel   string
	DeepseekTencentBaseURL string
	LLMRequestTimeout      time.Duration
}

func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Debug().Err(err).Msg(".env file not found or failed to load (this is normal in production)")
	} else {
		log.Info().Msg(".env file loaded successfully")
	}

	cfg := &Config{
		// Database defaults
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "lunch_user"),
		DBPassword: getEnv("DB_PASSWORD", "1234"),
		DBName:     getEnv("DB_NAME", "lunch_delivery"),

		// LLM defaults
		LLMType:                getEnv("LLM_TYPE", "DEEPSEEK_TENCENT"),
		DeepseekTencentAPIKey:  getEnv("DEEPSEEK_TENCENT_API_KEY", ""),
		DeepseekTencentModel:   getEnv("DEEPSEEK_TENCENT_MODEL", "deepseek-v3"),
		DeepseekTencentBaseURL: getEnv("DEEPSEEK_TENCENT_BASE_URL", "https://api.lkeap.tencentcloud.com/v1"),
		LLMRequestTimeout:      5 * time.Minute,
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}