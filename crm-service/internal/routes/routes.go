package routes

import (
	"divine-crm/internal/config"
	"divine-crm/internal/handlers"
	"divine-crm/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// Handlers holds all HTTP handlers
type Handlers struct {
	Health  *handlers.HealthHandler
	Contact *handlers.ContactHandler
	Lead    *handlers.LeadHandler
	Product *handlers.ProductHandler
	Chat    *handlers.ChatHandler
	AI      *handlers.AIHandler
	Webhook *handlers.WebhookHandler
}

// SetupRoutes sets up all application routes
func SetupRoutes(app *fiber.App, h Handlers, cfg *config.Config) {
	// Health check (no auth)
	app.Get("/health", h.Health.Check)

	// API v1
	api := app.Group("/api/v1")

	// Webhook routes (no auth)
	webhooks := api.Group("/webhooks")
	webhooks.Get("/whatsapp", h.Webhook.VerifyWhatsApp)
	webhooks.Post("/whatsapp", h.Webhook.HandleWhatsApp)
	webhooks.Get("/telegram", h.Webhook.VerifyTelegram)
	webhooks.Post("/telegram", h.Webhook.HandleTelegram)
	webhooks.Get("/instagram", h.Webhook.VerifyInstagram)
	webhooks.Post("/instagram", h.Webhook.HandleInstagram)

	// Protected routes (require auth)
	masterdata := api.Group("/masterdata", middleware.AuthMiddleware(cfg))

	// Contact routes
	masterdata.Get("/contacts", h.Contact.GetAll)
	masterdata.Get("/contacts/:id", h.Contact.GetByID)
	masterdata.Post("/contacts", h.Contact.Create)
	masterdata.Put("/contacts/:id", h.Contact.Update)
	masterdata.Delete("/contacts/:id", h.Contact.Delete)

	// Lead routes
	masterdata.Get("/leads", h.Lead.GetAll)
	masterdata.Get("/leads/:id", h.Lead.GetByID)
	masterdata.Post("/leads", h.Lead.Create)
	masterdata.Put("/leads/:id", h.Lead.Update)
	masterdata.Delete("/leads/:id", h.Lead.Delete)

	// Product routes
	masterdata.Get("/products", h.Product.GetAll)
	masterdata.Get("/products/:id", h.Product.GetByID)
	masterdata.Post("/products", h.Product.Create)
	masterdata.Put("/products/:id", h.Product.Update)
	masterdata.Delete("/products/:id", h.Product.Delete)

	// Chat routes
	chats := api.Group("/chats", middleware.AuthMiddleware(cfg))
	chats.Get("/", h.Chat.GetAll)
	chats.Get("/:id", h.Chat.GetByID)
	chats.Post("/", h.Chat.Create)
	chats.Put("/:id", h.Chat.Update)
	chats.Post("/:id/assign", h.Chat.Assign)
	chats.Post("/:id/resolve", h.Chat.Resolve)

	// AI routes
	ai := api.Group("/ai", middleware.AuthMiddleware(cfg))
	ai.Post("/process", h.AI.Process)
	ai.Get("/token-balance", h.AI.GetTokenBalances)
	ai.Get("/configurations", h.AI.GetConfigurations)
	ai.Get("/agents", h.AI.GetAgents)
}
