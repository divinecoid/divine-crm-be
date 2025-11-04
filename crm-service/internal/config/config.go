package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	CORS      CORSConfig
	WhatsApp  WhatsAppConfig
	RateLimit RateLimitConfig
	Logging   LoggingConfig
	AI        AIConfig
}

type ServerConfig struct {
	Port        string
	AppName     string
	Environment string
}

type DatabaseConfig struct {
	URL             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type JWTConfig struct {
	Secret     []byte
	Expiration time.Duration
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

type WhatsAppConfig struct {
	VerifyToken string
	APIVersion  string
}

type RateLimitConfig struct {
	Enabled bool
	Max     int
	Window  time.Duration
}

type LoggingConfig struct {
	Level  string
	Format string
}

type AIConfig struct {
	RequestTimeout time.Duration
	MaxRetries     int
}

var AppConfig *Config

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port:        getEnv("PORT", "3002"),
			AppName:     getEnv("APP_NAME", "Divine CRM"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", "host=localhost user=divine_user password=divine_password_123 dbname=divine_crm port=5432 sslmode=disable"),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 1*time.Hour),
		},
		JWT: JWTConfig{
			Secret:     []byte(getEnv("JWT_SECRET", "divine-crm-secret-key-change-in-production")),
			Expiration: getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
			AllowedMethods: getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders: getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization"}),
		},
		WhatsApp: WhatsAppConfig{
			VerifyToken: getEnv("WHATSAPP_VERIFY_TOKEN", "divine-crm-webhook-2025"),
			APIVersion:  getEnv("WHATSAPP_API_VERSION", "v18.0"),
		},
		RateLimit: RateLimitConfig{
			Enabled: getEnvAsBool("RATE_LIMIT_ENABLED", true),
			Max:     getEnvAsInt("RATE_LIMIT_MAX", 100),
			Window:  getEnvAsDuration("RATE_LIMIT_WINDOW", 1*time.Minute),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		AI: AIConfig{
			RequestTimeout: getEnvAsDuration("AI_REQUEST_TIMEOUT", 30*time.Second),
			MaxRetries:     getEnvAsInt("AI_MAX_RETRIES", 3),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	AppConfig = config
	return config, nil
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("PORT is required")
	}

	if len(c.JWT.Secret) < 32 {
		log.Println("⚠️  WARNING: JWT_SECRET should be at least 32 characters for production")
	}

	if c.Server.Environment == "production" {
		if string(c.JWT.Secret) == "divine-crm-secret-key-change-in-production" {
			return fmt.Errorf("JWT_SECRET must be changed in production")
		}
	}

	return nil
}

func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
