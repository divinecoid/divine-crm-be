package database

import (
	"context"
	"divine-crm/internal/config"
	"divine-crm/internal/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os" // ✅ TAMBAHKAN INI!
	"time"
)

// InitDB initializes database connection with proper configuration
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.Database.URL

	customLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info, // ✅ Change to Info
			IgnoreRecordNotFoundError: false,       // ✅ Also show not found temporarily
			Colorful:                  true,
		},
	)

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: customLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB for connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connected successfully")
	log.Printf("   Max Idle Connections: %d", cfg.Database.MaxIdleConns)
	log.Printf("   Max Open Connections: %d", cfg.Database.MaxOpenConns)
	log.Printf("   Connection Max Lifetime: %s", cfg.Database.ConnMaxLifetime)

	return db, nil
}

// RunMigrations runs auto migrations for all models
func RunMigrations(db *gorm.DB) error {
	log.Println("🔄 Running database migrations...")

	// Drop and recreate (ONLY for development!)
	// Comment this out for production
	db.Migrator().DropTable(
		&models.Contact{},
		&models.Lead{},
		&models.Product{},
		&models.ChatMessage{},
		&models.ChatLabel{},
		&models.QuickReply{},
		&models.BroadcastTemplate{},
		&models.AIConfiguration{},
		&models.AIAgent{},
		&models.TokenBalance{},
		&models.ConnectedPlatform{},
	)

	// Create tables one by one with error handling
	models := []interface{}{
		&models.Contact{},
		&models.Lead{},
		&models.Product{},
		&models.ChatMessage{},
		&models.ChatLabel{},
		&models.QuickReply{},
		&models.BroadcastTemplate{},
		&models.AIConfiguration{},
		&models.AIAgent{},
		&models.TokenBalance{},
		&models.ConnectedPlatform{},
	}

	for _, model := range models {
		modelName := fmt.Sprintf("%T", model)
		log.Printf("Creating table for %s...", modelName)

		if err := db.Migrator().CreateTable(model); err != nil {
			return fmt.Errorf("failed to create table for %s: %w", modelName, err)
		}

		log.Printf("✅ Table created for %s", modelName)
	}

	log.Println("✅ Database migrations completed")
	return nil
}

// HealthCheck checks if database is healthy
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return err
	}

	return nil
}

// GetStats returns database statistics
func GetStats(db *gorm.DB) (map[string]interface{}, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"open_connections":    stats.OpenConnections,
		"in_use":              stats.InUse,
		"idle":                stats.Idle,
		"wait_count":          stats.WaitCount,
		"wait_duration":       stats.WaitDuration.String(),
		"max_idle_closed":     stats.MaxIdleClosed,
		"max_lifetime_closed": stats.MaxLifetimeClosed,
	}, nil
}

func getLogLevel(level string) logger.LogLevel {
	switch level {
	case "debug":
		return logger.Info
	case "info":
		return logger.Warn
	case "warn":
		return logger.Error
	case "error":
		return logger.Error
	default:
		return logger.Silent
	}
}
