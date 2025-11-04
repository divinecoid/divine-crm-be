package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"divine-crm/internal/config"
	"divine-crm/internal/database"
	"divine-crm/internal/handlers"
	"divine-crm/internal/middleware"
	"divine-crm/internal/repository"
	"divine-crm/internal/routes"
	"divine-crm/internal/services"
	"divine-crm/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize logger
	appLogger := utils.NewLogger(cfg.Logging.Level, cfg.Logging.Format)

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		appLogger.Fatal("Failed to initialize database", "error", err)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		appLogger.Fatal("Failed to run migrations", "error", err)
	}

	appLogger.Info("Database initialized successfully")

	// Initialize repositories
	contactRepo := repository.NewContactRepository(db)
	leadRepo := repository.NewLeadRepository(db)
	productRepo := repository.NewProductRepository(db)
	chatRepo := repository.NewChatRepository(db)
	aiRepo := repository.NewAIRepository(db)
	platformRepo := repository.NewPlatformRepository(db)

	// Initialize services
	aiService := services.NewAIService(aiRepo, cfg)
	whatsappService := services.NewWhatsAppService(platformRepo, appLogger)
	contactService := services.NewContactService(contactRepo, appLogger)
	chatService := services.NewChatService(chatRepo, contactService, aiService, appLogger)
	webhookService := services.NewWebhookService(chatService, whatsappService, appLogger)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(db)
	contactHandler := handlers.NewContactHandler(contactService)
	leadHandler := handlers.NewLeadHandler(leadRepo)
	productHandler := handlers.NewProductHandler(productRepo)
	chatHandler := handlers.NewChatHandler(chatService)
	aiHandler := handlers.NewAIHandler(aiService, aiRepo)
	webhookHandler := handlers.NewWebhookHandler(webhookService, cfg)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      cfg.Server.AppName,
		ErrorHandler: middleware.CustomErrorHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} ${latency}\n",
	}))

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins[0],
		AllowHeaders:     cfg.CORS.AllowedHeaders[0],
		AllowMethods:     cfg.CORS.AllowedMethods[0],
		AllowCredentials: true,
	}))

	// Rate limiting (if enabled)
	if cfg.RateLimit.Enabled {
		app.Use(middleware.RateLimiter(cfg.RateLimit.Max, cfg.RateLimit.Window))
	}

	// Skip browser warning header
	app.Use(func(c *fiber.Ctx) error {
		c.Set("ngrok-skip-browser-warning", "true")
		return c.Next()
	})

	// Setup routes
	routes.SetupRoutes(app, routes.Handlers{
		Health:   healthHandler,
		Contact:  contactHandler,
		Lead:     leadHandler,
		Product:  productHandler,
		Chat:     chatHandler,
		AI:       aiHandler,
		Webhook:  webhookHandler,
	}, cfg)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		appLogger.Info("Shutting down server...")

		// Shutdown with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(ctx); err != nil {
			appLogger.Error("Server forced to shutdown", "error", err)
		}

		// Close database connection
		sqlDB, _ := db.DB()
		if err := sqlDB.Close(); err != nil {
			appLogger.Error("Failed to close database", "error", err)
		}

		appLogger.Info("Server shutdown complete")
		os.Exit(0)
	}()

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	appLogger.Info("Starting server", 
		"port", cfg.Server.Port,
		"environment", cfg.Server.Environment,
	)

	if err := app.Listen(addr); err != nil {
		appLogger.Fatal("Failed to start server", "error", err)
	}
}
