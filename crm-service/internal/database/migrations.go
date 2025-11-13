package database

import (
	"divine-crm/internal/models"
	"fmt"
	"gorm.io/gorm"
	"log"
)

// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB) error {
	log.Println("🚀 Running database migrations...")

	// ============================================
	// STEP 1: Install pgvector extension
	// ============================================
	log.Println("📦 Installing pgvector extension...")

	// Check if extension exists
	var extExists bool
	err := db.Raw("SELECT EXISTS(SELECT 1 FROM pg_extension WHERE extname = 'vector')").Scan(&extExists).Error
	if err != nil {
		log.Printf("⚠️  Could not check extension: %v", err)
	}

	if !extExists {
		log.Println("Installing vector extension...")
		if err := db.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; err != nil {
			log.Printf("❌ Failed to install pgvector: %v", err)
			log.Println("⚠️  Vector features will be disabled")
			log.Println("💡 Make sure you're using pgvector/pgvector:pg16 Docker image")

			// Continue without vector features
			return migrateWithoutVector(db)
		}
	}

	// Get version
	var version string
	if err := db.Raw("SELECT extversion FROM pg_extension WHERE extname = 'vector'").Scan(&version).Error; err == nil {
		log.Printf("✅ pgvector extension ready (version: %s)", version)
	}

	// Test vector type
	testErr := db.Exec("SELECT 'test'::vector").Error
	if testErr != nil {
		log.Printf("⚠️  Vector type test failed: %v", testErr)
		log.Println("Attempting to fix...")

		// Try to drop and recreate
		db.Exec("DROP EXTENSION IF EXISTS vector CASCADE")
		if err := db.Exec("CREATE EXTENSION vector").Error; err != nil {
			log.Printf("❌ Could not fix vector extension: %v", err)
			return migrateWithoutVector(db)
		}
	}

	// ============================================
	// STEP 2: Migrate all tables
	// ============================================
	log.Println("📋 Migrating database tables...")

	err = db.AutoMigrate(
		// Core entities
		&models.Contact{},
		&models.Product{},
		&models.ChatLabel{},
		&models.ChatMessage{},

		// AI & Platforms
		&models.AIConfiguration{},
		&models.AIAgent{},
		&models.ConnectedPlatform{},
		&models.HumanAgent{},

		// Communication
		&models.BroadcastTemplate{},
		&models.BroadcastHistory{},
		&models.QuickReply{},

		// Settings & Analytics
		&models.APISettings{},
		&models.TokenBalance{},
		&models.Analytics{},

		// Vector Embeddings (with pgvector)
		&models.KnowledgeBase{},
		&models.ChatHistory{},
		&models.ProductEmbedding{},
		&models.FAQEmbedding{},
	)

	if err != nil {
		log.Printf("❌ Migration failed: %v", err)
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

// migrateWithoutVector migrates without vector-dependent tables
func migrateWithoutVector(db *gorm.DB) error {
	log.Println("⚠️  Migrating without vector features...")

	err := db.AutoMigrate(
		// Core entities
		&models.Contact{},
		&models.Product{},
		&models.ChatLabel{},
		&models.ChatMessage{},

		// AI & Platforms
		&models.AIConfiguration{},
		&models.AIAgent{},
		&models.ConnectedPlatform{},
		&models.HumanAgent{},

		// Communication
		&models.BroadcastTemplate{},
		&models.BroadcastHistory{},
		&models.QuickReply{},

		// Settings & Analytics
		&models.APISettings{},
		&models.TokenBalance{},
		&models.Analytics{},
	)

	if err != nil {
		return fmt.Errorf("migration without vector failed: %w", err)
	}

	log.Println("✅ Database migrations completed (without vector features)")
	return nil
}
