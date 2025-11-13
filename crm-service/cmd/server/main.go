package main

import (
	"divine-crm/internal/config"
	"divine-crm/internal/database"
	"divine-crm/internal/handlers"
	"divine-crm/internal/repository"
	"divine-crm/internal/services"
	"divine-crm/internal/utils"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	appLogger := utils.NewLogger(cfg.Logging.Level)
	appLogger.Info("Starting Divine CRM Service with Vector Embeddings...")

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		appLogger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	appLogger.Info("Database connected successfully")

	// Run migrations
	appLogger.Info("Running database migrations...")
	if err := database.RunMigrations(db); err != nil {
		appLogger.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}
	appLogger.Info("Database migrations completed successfully")

	// ==================== INITIALIZE REPOSITORIES ====================
	appLogger.Info("Initializing repositories...")

	contactRepo := repository.NewContactRepository(db)
	productRepo := repository.NewProductRepository(db)
	chatLabelRepo := repository.NewChatLabelRepository(db)
	chatRepo := repository.NewChatRepository(db)
	aiRepo := repository.NewAIRepository(db)
	platformRepo := repository.NewPlatformRepository(db)
	humanAgentRepo := repository.NewHumanAgentRepository(db)
	broadcastRepo := repository.NewBroadcastRepository(db)
	quickReplyRepo := repository.NewQuickReplyRepository(db)
	vectorRepo := repository.NewVectorRepository(db)

	// ==================== INITIALIZE SERVICES ====================
	appLogger.Info("Initializing services...")

	// Core services
	contactService := services.NewContactService(contactRepo, appLogger)
	productService := services.NewProductService(productRepo, appLogger)
	chatLabelService := services.NewChatLabelService(chatLabelRepo, appLogger)
	quickReplyService := services.NewQuickReplyService(quickReplyRepo, appLogger)
	vectorService := services.NewVectorService(vectorRepo, cfg, appLogger)

	// ✅ AI Service now uses vectorService
	aiService := services.NewAIService(aiRepo, vectorService, cfg)

	humanAgentService := services.NewHumanAgentService(humanAgentRepo, appLogger, cfg)

	// Platform services
	whatsappService := services.NewWhatsAppService(platformRepo, appLogger)
	instagramService := services.NewInstagramService(platformRepo, appLogger)
	telegramService := services.NewTelegramService(platformRepo, appLogger)

	// Broadcast service
	broadcastService := services.NewBroadcastService(
		broadcastRepo,
		contactRepo,
		whatsappService,
		instagramService,
		telegramService,
		appLogger,
	)

	// Chat service (with product context)
	chatService := services.NewChatService(
		chatRepo,
		contactService,
		aiService,
		productService,
		appLogger,
	)

	// Webhook service
	webhookService := services.NewWebhookService(
		chatService,
		whatsappService,
		instagramService,
		telegramService,
		appLogger,
	)

	// Platform service
	platformService := services.NewPlatformService(platformRepo, appLogger)

	// AI Agent service
	aiAgentService := services.NewAIAgentService(aiRepo, appLogger)

	// AI Config service
	aiConfigService := services.NewAIConfigService(aiRepo, appLogger)

	// ==================== INITIALIZE HANDLERS ====================
	appLogger.Info("Initializing handlers...")

	contactHandler := handlers.NewContactHandler(contactService)
	productHandler := handlers.NewProductHandler(productService)
	chatLabelHandler := handlers.NewChatLabelHandler(chatLabelService)
	chatHandler := handlers.NewChatHandler(chatService)
	aiConfigHandler := handlers.NewAIConfigHandler(aiConfigService)
	aiAgentHandler := handlers.NewAIAgentHandler(aiAgentService)
	platformHandler := handlers.NewPlatformHandler(platformService)
	humanAgentHandler := handlers.NewHumanAgentHandler(humanAgentService)
	broadcastHandler := handlers.NewBroadcastHandler(broadcastService)
	quickReplyHandler := handlers.NewQuickReplyHandler(quickReplyService)
	webhookHandler := handlers.NewWebhookHandler(webhookService, cfg)
	vectorHandler := handlers.NewVectorHandler(vectorService)

	// ==================== INITIALIZE FIBER APP ====================
	app := fiber.New(fiber.Config{
		AppName:      "Divine CRM v1.0 + Vector RAG",
		ServerHeader: "Divine CRM",
		BodyLimit:    4 * 1024 * 1024, // 4MB
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			appLogger.Error("Request error", "error", err, "path", c.Path())

			return c.Status(code).JSON(fiber.Map{
				"error":   err.Error(),
				"success": false,
			})
		},
	})

	// ==================== MIDDLEWARE ====================

	// Recover from panics
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// Request ID
	app.Use(requestid.New())

	// Logger
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} ${latency}\n",
		TimeFormat: "15:04:05",
		TimeZone:   "Asia/Jakarta",
	}))

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     cfg.CORS.AllowMethods,
		AllowHeaders:     cfg.CORS.AllowHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           cfg.CORS.MaxAge,
	}))

	// Rate limiter
	app.Use(limiter.New(limiter.Config{
		Max:        cfg.RateLimit.RequestsPerMinute,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		},
	}))

	// Custom request logging
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		if duration > 1*time.Second {
			appLogger.Warn("Slow request",
				"method", c.Method(),
				"path", c.Path(),
				"duration", duration.String(),
			)
		}

		return err
	})
	appLogger.Info("Setting up routes...")

	setupRoutes(
		app,
		cfg,
		contactHandler,
		productHandler,
		chatLabelHandler,
		chatHandler,
		aiConfigHandler,
		aiAgentHandler,
		platformHandler,
		humanAgentHandler,
		broadcastHandler,
		quickReplyHandler,
		webhookHandler,
		vectorHandler,
	)

	// ==================== START SERVER ====================
	port := cfg.Server.Port
	appLogger.Info("Server starting with Vector RAG capabilities", "port", port)

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		appLogger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
