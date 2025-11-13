package main

import (
	"divine-crm/internal/config"
	"divine-crm/internal/handlers"
	"divine-crm/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(
	app *fiber.App,
	cfg *config.Config,
	contactHandler *handlers.ContactHandler,
	productHandler *handlers.ProductHandler,
	chatLabelHandler *handlers.ChatLabelHandler,
	chatHandler *handlers.ChatHandler,
	aiConfigHandler *handlers.AIConfigHandler,
	aiAgentHandler *handlers.AIAgentHandler,
	platformHandler *handlers.PlatformHandler,
	humanAgentHandler *handlers.HumanAgentHandler,
	broadcastHandler *handlers.BroadcastHandler,
	quickReplyHandler *handlers.QuickReplyHandler,
	webhookHandler *handlers.WebhookHandler,
	vectorHandler *handlers.VectorHandler,
) {
	// ==================== ROOT ====================
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "Divine CRM API",
			"version": "1.0.0",
			"status":  "running",
			"features": []string{
				"Multi-platform messaging (WhatsApp, Instagram, Telegram)",
				"AI-powered responses with RAG",
				"Vector embeddings for semantic search",
				"Contact & lead management",
				"Product catalog",
				"Broadcast messaging",
				"Analytics & reporting",
			},
		})
	})

	// ==================== API V1 ====================
	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "divine-crm",
			"version": "1.0.0",
		})
	})

	// ==================== PUBLIC ROUTES ====================

	// Webhooks (no auth)
	webhooks := api.Group("/webhooks")
	webhooks.Get("/whatsapp", webhookHandler.VerifyWhatsApp)
	webhooks.Post("/whatsapp", webhookHandler.HandleWhatsApp)
	webhooks.Get("/instagram", webhookHandler.VerifyInstagram)
	webhooks.Post("/instagram", webhookHandler.HandleInstagram)
	webhooks.Get("/telegram", webhookHandler.VerifyTelegram)
	webhooks.Post("/telegram", webhookHandler.HandleTelegram)

	// Auth (no auth for login)
	auth := api.Group("/auth")
	auth.Post("/login", humanAgentHandler.Login)

	// ==================== PROTECTED ROUTES ====================

	// Create protected group with auth middleware
	protected := api.Group("", middleware.AuthMiddleware(cfg))

	// Contacts
	contacts := protected.Group("/contacts")
	contacts.Get("/", contactHandler.GetAll)
	contacts.Get("/stats", contactHandler.GetStats)
	contacts.Get("/search", contactHandler.Search)
	contacts.Get("/status", contactHandler.GetByStatus)
	contacts.Get("/temperature", contactHandler.GetByTemperature)
	contacts.Get("/:id", contactHandler.GetByID)
	contacts.Post("/", contactHandler.Create)
	contacts.Put("/:id", contactHandler.Update)
	contacts.Delete("/:id", contactHandler.Delete)
	contacts.Patch("/:id/temperature", contactHandler.UpdateTemperature)
	contacts.Patch("/:id/status", contactHandler.UpdateStatus)

	// Products
	products := protected.Group("/products")
	products.Get("/", productHandler.GetAll)
	products.Get("/active", productHandler.GetActive)
	products.Get("/search", productHandler.Search)
	products.Get("/:id", productHandler.GetByID)
	products.Post("/", productHandler.Create)
	products.Put("/:id", productHandler.Update)
	products.Delete("/:id", productHandler.Delete)

	// Chat Labels
	labels := protected.Group("/chat-labels")
	labels.Get("/", chatLabelHandler.GetAll)
	labels.Get("/:id", chatLabelHandler.GetByID)
	labels.Post("/", chatLabelHandler.Create)
	labels.Put("/:id", chatLabelHandler.Update)
	labels.Delete("/:id", chatLabelHandler.Delete)

	// Chats
	chats := protected.Group("/chats")
	chats.Get("/", chatHandler.GetAll)
	chats.Get("/stats", chatHandler.GetStats)
	chats.Get("/unassigned", chatHandler.GetUnassigned)
	chats.Get("/assigned", chatHandler.GetAssigned)
	chats.Get("/resolved", chatHandler.GetResolved)
	chats.Get("/:id", chatHandler.GetByID)
	chats.Post("/", chatHandler.Create)
	chats.Put("/:id", chatHandler.Update)
	chats.Patch("/:id/assign", chatHandler.Assign)
	chats.Patch("/:id/resolve", chatHandler.Resolve)
	chats.Patch("/:id/takeover", chatHandler.TakeOver)
	chats.Patch("/:id/to-ai", chatHandler.BackToAI)

	// AI Configurations
	aiConfigs := protected.Group("/ai/configurations")
	aiConfigs.Get("/", aiConfigHandler.GetAll)
	aiConfigs.Get("/active", aiConfigHandler.GetActive)
	aiConfigs.Get("/:id", aiConfigHandler.GetByID)
	aiConfigs.Post("/", aiConfigHandler.Create)
	aiConfigs.Put("/:id", aiConfigHandler.Update)
	aiConfigs.Delete("/:id", aiConfigHandler.Delete)

	// AI Agents
	aiAgents := protected.Group("/ai/agents")
	aiAgents.Get("/", aiAgentHandler.GetAll)
	aiAgents.Get("/active", aiAgentHandler.GetActive)
	aiAgents.Get("/:id", aiAgentHandler.GetByID)
	aiAgents.Post("/", aiAgentHandler.Create)
	aiAgents.Put("/:id", aiAgentHandler.Update)
	aiAgents.Delete("/:id", aiAgentHandler.Delete)
	aiAgents.Patch("/:id/activate", aiAgentHandler.Activate)

	// Platforms
	platforms := protected.Group("/platforms")
	platforms.Get("/", platformHandler.GetAll)
	platforms.Get("/active", platformHandler.GetActive)
	platforms.Get("/:id", platformHandler.GetByID)
	platforms.Post("/", platformHandler.Create)
	platforms.Put("/:id", platformHandler.Update)
	platforms.Delete("/:id", platformHandler.Delete)

	// Human Agents
	agents := protected.Group("/agents")
	agents.Get("/", humanAgentHandler.GetAll)
	agents.Get("/active", humanAgentHandler.GetActive)
	agents.Get("/:id", humanAgentHandler.GetByID)
	agents.Post("/", humanAgentHandler.Create)
	agents.Put("/:id", humanAgentHandler.Update)
	agents.Delete("/:id", humanAgentHandler.Delete)
	agents.Post("/:id/revoke", humanAgentHandler.RevokeAccess)
	agents.Post("/:id/reset-password", humanAgentHandler.ResetPassword)

	// Broadcast
	broadcast := protected.Group("/broadcast")
	broadcast.Get("/templates", broadcastHandler.GetAllTemplates)
	broadcast.Get("/templates/:id", broadcastHandler.GetTemplateByID)
	broadcast.Post("/templates", broadcastHandler.CreateTemplate)
	broadcast.Put("/templates/:id", broadcastHandler.UpdateTemplate)
	broadcast.Delete("/templates/:id", broadcastHandler.DeleteTemplate)
	broadcast.Post("/send", broadcastHandler.SendBroadcast)
	broadcast.Get("/history", broadcastHandler.GetHistory)

	// Quick Replies
	quickReplies := protected.Group("/quick-replies")
	quickReplies.Get("/", quickReplyHandler.GetAll)
	quickReplies.Post("/", quickReplyHandler.Create)
	quickReplies.Put("/:id", quickReplyHandler.Update)
	quickReplies.Delete("/:id", quickReplyHandler.Delete)

	// ==================== VECTOR EMBEDDINGS & RAG ====================
	vectors := protected.Group("/vectors")

	// Knowledge Base
	vectors.Post("/knowledge", vectorHandler.AddKnowledge)
	vectors.Get("/knowledge", vectorHandler.GetAllKnowledge)
	vectors.Get("/knowledge/search", vectorHandler.SearchKnowledge)

	// FAQ
	vectors.Post("/faq", vectorHandler.AddFAQ)
	vectors.Get("/faq", vectorHandler.GetAllFAQ)
	vectors.Get("/faq/search", vectorHandler.SearchFAQ)

	// Product Semantic Search
	vectors.Post("/products/embedding", vectorHandler.AddProductEmbedding)
	vectors.Get("/products/search", vectorHandler.SearchProducts)

	// Chat History
	vectors.Get("/chat-history/:contactId", vectorHandler.GetSimilarConversations)
}
