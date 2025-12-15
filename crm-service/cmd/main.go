package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	initDB()
	initJWT()

	app := fiber.New(fiber.Config{AppName: "Divine CRM Service"})

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path}\n",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH, OPTIONS",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "crm-service"})
	})

	api := app.Group("/api/v1")

	masterdata := api.Group("/masterdata", authMiddleware)
	masterdata.Get("/contacts", getContacts)
	masterdata.Get("/contacts/:id", getContactByID)
	masterdata.Post("/contacts", createContact)
	masterdata.Put("/contacts/:id", updateContact)
	masterdata.Delete("/contacts/:id", deleteContact)

	masterdata.Get("/leads", getLeads)
	masterdata.Post("/leads", createLead)
	masterdata.Put("/leads/:id", updateLead)
	masterdata.Delete("/leads/:id", deleteLead)

	masterdata.Get("/products", getProducts)
	masterdata.Post("/products", createProduct)
	masterdata.Put("/products/:id", updateProduct)
	masterdata.Delete("/products/:id", deleteProduct)

	masterdata.Get("/chat-labels", getChatLabels)
	masterdata.Post("/chat-labels", createChatLabel)
	masterdata.Put("/chat-labels/:id", updateChatLabel)
	masterdata.Delete("/chat-labels/:id", deleteChatLabel)

	masterdata.Get("/quick-replies", getQuickReplies)
	masterdata.Post("/quick-replies", createQuickReply)
	masterdata.Put("/quick-replies/:id", updateQuickReply)
	masterdata.Delete("/quick-replies/:id", deleteQuickReply)

	masterdata.Get("/broadcast-templates", getBroadcastTemplates)
	masterdata.Post("/broadcast-templates", createBroadcastTemplate)
	masterdata.Put("/broadcast-templates/:id", updateBroadcastTemplate)
	masterdata.Delete("/broadcast-templates/:id", deleteBroadcastTemplate)

	masterdata.Get("/ai-configurations", getAIConfigurations)
	masterdata.Post("/ai-configurations", createAIConfiguration)
	masterdata.Put("/ai-configurations/:id", updateAIConfiguration)
	masterdata.Delete("/ai-configurations/:id", deleteAIConfiguration)

	masterdata.Get("/ai-agents", getAIAgents)
	masterdata.Post("/ai-agents", createAIAgent)
	masterdata.Put("/ai-agents/:id", updateAIAgent)
	masterdata.Delete("/ai-agents/:id", deleteAIAgent)

	masterdata.Get("/connected-platforms", getConnectedPlatforms)
	masterdata.Post("/connected-platforms", createConnectedPlatform)
	masterdata.Put("/connected-platforms/:id", updateConnectedPlatform)
	masterdata.Delete("/connected-platforms/:id", deleteConnectedPlatform)

	api.Get("/chats", authMiddleware, getChatMessages)
	api.Get("/chats/stats", authMiddleware, getChatStats)
	api.Get("/chats/conversations", authMiddleware, getChatConversations)
	api.Get("/chats/contact/:contactId", authMiddleware, getChatsByContact)
	api.Get("/chats/:id", authMiddleware, getChatMessageByID)
	api.Post("/chats", authMiddleware, createChatMessage)
	api.Put("/chats/:id", authMiddleware, updateChatMessage)
	api.Post("/chats/:id/assign", authMiddleware, assignChat)
	api.Post("/chats/:id/takeover", authMiddleware, takeoverChat)
	api.Post("/chats/:id/back-to-ai", authMiddleware, backToAIChat)
	api.Post("/chats/:id/pending", authMiddleware, setPendingChat)
	api.Post("/chats/:id/resolve", authMiddleware, resolveChat)
	api.Post("/chats/:id/reply", authMiddleware, sendReplyChat)
	api.Post("/chats/:id/labels", authMiddleware, addChatLabel)
	api.Delete("/chats/:id/labels/:labelId", authMiddleware, removeChatLabel)

	ai := api.Group("/ai", authMiddleware)
	ai.Post("/process", processWithAI)
	ai.Post("/suggest-reply", suggestReply)
	ai.Get("/token-balance", getTokenBalances)
	ai.Get("/debug-prompt", debugAIPrompt)

	analytics := api.Group("/analytics", authMiddleware)
	analytics.Get("/overview", getAnalyticsOverview)

	broadcast := api.Group("/broadcasts", authMiddleware)
	broadcast.Get("/", getBroadcasts)
	broadcast.Get("/:id", getBroadcastByID)
	broadcast.Post("/send", sendBroadcast)

	webhooks := api.Group("/webhooks")
	webhooks.Get("/whatsapp", verifyWhatsAppWebhook)
	webhooks.Post("/whatsapp", handleWhatsAppWebhook)
	webhooks.Get("/telegram", verifyTelegramWebhook)
	webhooks.Post("/telegram", handleTelegramWebhook)
	webhooks.Get("/instagram", verifyInstagramWebhook)
	webhooks.Post("/instagram", handleInstagramWebhook)

	port := getEnv("SERVER_PORT", "3002")

	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("📦 Divine CRM Service running on port %s", port)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	log.Fatal(app.Listen(":" + port))
}
