package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	OpenAI    OpenAIConfig
	DeepSeek  AIProviderConfig
	Grok      AIProviderConfig
	Gemini    GeminiConfig
	WhatsApp  WhatsAppConfig
	Instagram InstagramConfig
	Telegram  TelegramConfig
	JWT       JWTConfig
	CORS      CORSConfig
	RateLimit RateLimitConfig
	Logging   LoggingConfig
}

type ServerConfig struct {
	Port    string
	Host    string
	Env     string
	AppName string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime string
}

type OpenAIConfig struct {
	APIKey   string
	Model    string
	Endpoint string
}

type AIProviderConfig struct {
	APIKey   string
	Model    string
	Endpoint string
}

type GeminiConfig struct {
	APIKey   string
	Model    string
	Endpoint string
}

type WhatsAppConfig struct {
	AccessToken       string
	PhoneNumberID     string
	BusinessAccountID string
	VerifyToken       string
	APIVersion        string
}

type InstagramConfig struct {
	AccessToken   string
	PageID        string
	VerifyToken   string
	WebhookSecret string
}

type TelegramConfig struct {
	BotToken      string
	WebhookURL    string
	VerifyToken   string
	WebhookSecret string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

type CORSConfig struct {
	AllowOrigins     string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials bool
	MaxAge           int
}

type RateLimitConfig struct {
	Enabled           bool
	RequestsPerMinute int
	Max               int
	Window            string
}

type LoggingConfig struct {
	Level  string
	Format string
}

func LoadConfig() (*Config, error) {
	// Load .env file if exists
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port:    getEnv("SERVER_PORT", getEnv("PORT", "3002")),
			Host:    getEnv("SERVER_HOST", "0.0.0.0"),
			Env:     getEnv("ENVIRONMENT", "development"),
			AppName: getEnv("APP_NAME", "Divine CRM"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "dahsyan80"),
			DBName:          getEnv("DB_NAME", "divine_crm"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 100),
			ConnMaxLifetime: getEnv("DB_CONN_MAX_LIFETIME", "1h"),
		},
		OpenAI: OpenAIConfig{
			APIKey:   getEnv("OPENAI_API_KEY", ""),
			Model:    getEnv("OPENAI_MODEL", "gpt-3.5-turbo"),
			Endpoint: getEnv("OPENAI_ENDPOINT", "https://api.openai.com/v1/chat/completions"),
		},
		DeepSeek: AIProviderConfig{
			APIKey:   getEnv("DEEPSEEK_API_KEY", ""),
			Model:    getEnv("DEEPSEEK_MODEL", "deepseek-chat"),
			Endpoint: getEnv("DEEPSEEK_ENDPOINT", "https://api.deepseek.com/v1/chat/completions"),
		},
		Grok: AIProviderConfig{
			APIKey:   getEnv("GROK_API_KEY", ""),
			Model:    getEnv("GROK_MODEL", "grok-beta"),
			Endpoint: getEnv("GROK_ENDPOINT", "https://api.x.ai/v1/chat/completions"),
		},
		Gemini: GeminiConfig{
			APIKey:   getEnv("GEMINI_API_KEY", ""),
			Model:    getEnv("GEMINI_MODEL", "gemini-pro"),
			Endpoint: getEnv("GEMINI_ENDPOINT", "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent"),
		},
		WhatsApp: WhatsAppConfig{
			AccessToken:       getEnv("WHATSAPP_ACCESS_TOKEN", ""),
			PhoneNumberID:     getEnv("WHATSAPP_PHONE_NUMBER_ID", ""),
			BusinessAccountID: getEnv("WHATSAPP_BUSINESS_ACCOUNT_ID", ""),
			VerifyToken:       getEnv("WHATSAPP_VERIFY_TOKEN", "divine-crm-webhook-2025"),
			APIVersion:        getEnv("WHATSAPP_API_VERSION", "v18.0"),
		},
		Instagram: InstagramConfig{
			AccessToken:   getEnv("INSTAGRAM_ACCESS_TOKEN", ""),
			PageID:        getEnv("INSTAGRAM_PAGE_ID", ""),
			VerifyToken:   getEnv("INSTAGRAM_VERIFY_TOKEN", "divine-crm-webhook-2025"),
			WebhookSecret: getEnv("INSTAGRAM_WEBHOOK_SECRET", ""),
		},
		Telegram: TelegramConfig{
			BotToken:      getEnv("TELEGRAM_BOT_TOKEN", ""),
			WebhookURL:    getEnv("TELEGRAM_WEBHOOK_URL", ""),
			VerifyToken:   getEnv("TELEGRAM_VERIFY_TOKEN", "divine-crm-webhook-2025"),
			WebhookSecret: getEnv("TELEGRAM_WEBHOOK_SECRET", ""),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "divine-crm-super-secret-key"),
			ExpirationHours: getEnvInt("JWT_EXPIRATION_HOURS", 24),
		},
		CORS: CORSConfig{
			AllowOrigins:     getEnv("CORS_ALLOW_ORIGINS", "*"),
			AllowMethods:     getEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,DELETE,PATCH,OPTIONS"),
			AllowHeaders:     getEnv("CORS_ALLOW_HEADERS", "Origin,Content-Type,Accept,Authorization"),
			AllowCredentials: getEnvBool("CORS_ALLOW_CREDENTIALS", true),
			MaxAge:           getEnvInt("CORS_MAX_AGE", 3600),
		},
		RateLimit: RateLimitConfig{
			Enabled:           getEnvBool("RATE_LIMIT_ENABLED", true),
			RequestsPerMinute: getEnvInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 60),
			Max:               getEnvInt("RATE_LIMIT_MAX", 100),
			Window:            getEnv("RATE_LIMIT_WINDOW", "1m"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}
