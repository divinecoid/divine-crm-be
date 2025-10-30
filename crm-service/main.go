package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var jwtSecret []byte

// ==================== MODELS ====================

// Contact & Lead Models
type Contact struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Code          string    `json:"code" gorm:"uniqueIndex"`
	Channel       string    `json:"channel"`
	ChannelID     string    `json:"channel_id"`
	Name          string    `json:"name"`
	Temperature   string    `json:"temperature"`
	FirstContact  time.Time `json:"first_contact"`
	LastContact   time.Time `json:"last_contact"`
	LastAgent     string    `json:"last_agent"`
	LastAgentType string    `json:"last_agent_type"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Lead struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Code          string    `json:"code" gorm:"uniqueIndex"`
	Channel       string    `json:"channel"`
	ChannelID     string    `json:"channel_id"`
	Temperature   string    `json:"temperature"`
	FirstContact  time.Time `json:"first_contact"`
	LastContact   time.Time `json:"last_contact"`
	LastAgent     string    `json:"last_agent"`
	LastAgentType string    `json:"last_agent_type"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Product struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Code       string    `json:"code" gorm:"uniqueIndex"`
	Name       string    `json:"name"`
	Price      float64   `json:"price"`
	Stock      int       `json:"stock"`
	UploadedBy string    `json:"uploaded_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Chat Models
type ChatMessage struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	ContactID     uint      `json:"contact_id"`
	ContactName   string    `json:"contact_name"`
	Message       string    `json:"message" gorm:"type:text"`
	Response      string    `json:"response" gorm:"type:text"`
	Status        string    `json:"status" gorm:"default:'Unassigned'"`
	AssignedTo    string    `json:"assigned_to"`
	AssignedAgent string    `json:"assigned_agent"`
	Channel       string    `json:"channel"`
	Labels        string    `json:"labels" gorm:"type:text"`
	TokensUsed    int       `json:"tokens_used" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ChatLabel struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Label       string    `json:"label"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type QuickReply struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Trigger   string    `json:"trigger" gorm:"uniqueIndex"`
	Response  string    `json:"response" gorm:"type:text"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BroadcastTemplate struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Content   string    `json:"content" gorm:"type:text"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AI Models
type AIConfiguration struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	AIEngine  string    `json:"ai_engine" gorm:"uniqueIndex"`
	Token     string    `json:"token"`
	Endpoint  string    `json:"endpoint"`
	Model     string    `json:"model"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AIAgent struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex"`
	AIEngine    string    `json:"ai_engine"`
	BasicPrompt string    `json:"basic_prompt" gorm:"type:text"`
	Active      bool      `json:"active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TokenBalance struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	AIEngine        string    `json:"ai_engine" gorm:"uniqueIndex"`
	TotalTokens     int       `json:"total_tokens"`
	UsedTokens      int       `json:"used_tokens" gorm:"default:0"`
	RemainingTokens int       `json:"remaining_tokens"`
	LastReset       time.Time `json:"last_reset"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Webhook Models
type ConnectedPlatform struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Platform      string    `json:"platform" gorm:"uniqueIndex"`
	Token         string    `json:"token"`
	ClientID      string    `json:"client_id"`
	ClientSecret  string    `json:"client_secret"`
	WebhookURL    string    `json:"webhook_url"`
	PhoneNumberID string    `json:"phone_number_id"`
	Active        bool      `json:"active" gorm:"default:true"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// AI Request Models
type AIRequest struct {
	Model       string      `json:"model"`
	Messages    []AIMessage `json:"messages"`
	Temperature float64     `json:"temperature,omitempty"`
	MaxTokens   int         `json:"max_tokens,omitempty"`
}

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

// Webhook Models
type WhatsAppWebhook struct {
	Object string `json:"object"`
	Entry  []struct {
		Changes []struct {
			Value struct {
				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					} `json:"profile"`
					WaID string `json:"wa_id"`
				} `json:"contacts"`
				Messages []struct {
					From string `json:"from"`
					Text struct {
						Body string `json:"body"`
					} `json:"text"`
				} `json:"messages"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// ==================== INITIALIZATION ====================

func initDB() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=divine_crm port=5432 sslmode=disable"
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate all models
	db.AutoMigrate(
		&Contact{}, &Lead{}, &Product{},
		&ChatMessage{}, &ChatLabel{}, &QuickReply{}, &BroadcastTemplate{},
		&AIConfiguration{}, &AIAgent{}, &TokenBalance{},
		&ConnectedPlatform{},
	)

	log.Println("✅ CRM Service: Database connected and migrated")
}

func main() {
	godotenv.Load()
	initDB()

	// JWT Secret
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("divine-crm-secret-key-change-in-production")
	}

	app := fiber.New(fiber.Config{
		AppName: "CRM Service",
	})

	// Middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "crm-service",
			"time":    time.Now(),
		})
	})

	api := app.Group("/api/v1")

	// Contact & Lead & Product Routes (Protected)
	masterdata := api.Group("/masterdata", authMiddleware)

	// Contacts
	masterdata.Get("/contacts", getContacts)
	masterdata.Get("/contacts/:id", getContactByID)
	masterdata.Post("/contacts", createContact)
	masterdata.Put("/contacts/:id", updateContact)
	masterdata.Delete("/contacts/:id", deleteContact)

	// Leads
	masterdata.Get("/leads", getLeads)
	masterdata.Post("/leads", createLead)
	masterdata.Put("/leads/:id", updateLead)
	masterdata.Delete("/leads/:id", deleteLead)

	// Products
	masterdata.Get("/products", getProducts)
	masterdata.Post("/products", createProduct)
	masterdata.Put("/products/:id", updateProduct)
	masterdata.Delete("/products/:id", deleteProduct)

	// Chat Routes
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

	// AI Routes
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

	// Chat Messages
	api.Get("/chats", authMiddleware, getChatMessages)
	api.Get("/chats/:id", authMiddleware, getChatMessageByID)
	api.Post("/chats", authMiddleware, createChatMessage)
	api.Put("/chats/:id", authMiddleware, updateChatMessage)
	api.Post("/chats/:id/assign", authMiddleware, assignChat)
	api.Post("/chats/:id/resolve", authMiddleware, resolveChat)

	// AI Processing
	ai := api.Group("/ai", authMiddleware)
	ai.Post("/process", processWithAI)
	ai.Post("/suggest-reply", suggestReply)
	ai.Get("/token-balance", getTokenBalances)
	ai.Post("/token-balance/reset", resetTokenBalance)

	// Analytics
	analytics := api.Group("/analytics", authMiddleware)
	analytics.Get("/overview", getAnalyticsOverview)
	analytics.Get("/agent-performance", getAgentPerformance)

	// Webhooks (Public - no auth)
	webhooks := api.Group("/webhooks")
	webhooks.Get("/whatsapp", verifyWhatsAppWebhook)
	webhooks.Post("/whatsapp", handleWhatsAppWebhook)
	webhooks.Post("/telegram", handleTelegramWebhook)
	webhooks.Get("/instagram", verifyInstagramWebhook)
	webhooks.Post("/instagram", handleInstagramWebhook)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
	}

	log.Printf("🚀 CRM Service starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

// ==================== CONTACT HANDLERS ====================

func getContacts(c *fiber.Ctx) error {
	var contacts []Contact
	db.Order("created_at desc").Find(&contacts)
	return c.JSON(fiber.Map{"success": true, "data": contacts})
}

func getContactByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var contact Contact
	if err := db.First(&contact, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Contact not found"})
	}
	return c.JSON(fiber.Map{"success": true, "data": contact})
}

func createContact(c *fiber.Ctx) error {
	contact := new(Contact)
	if err := c.BodyParser(contact); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if contact.Code == "" {
		contact.Code = generateCode("C", &Contact{})
	}
	if err := db.Create(contact).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "data": contact})
}

func updateContact(c *fiber.Ctx) error {
	id := c.Params("id")
	contact := new(Contact)
	if err := db.First(contact, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Contact not found"})
	}
	if err := c.BodyParser(contact); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	db.Save(contact)
	return c.JSON(fiber.Map{"success": true, "data": contact})
}

func deleteContact(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.Delete(&Contact{}, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Contact not found"})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Contact deleted"})
}

// ==================== LEAD HANDLERS ====================

func getLeads(c *fiber.Ctx) error {
	var leads []Lead
	db.Order("created_at desc").Find(&leads)
	return c.JSON(fiber.Map{"success": true, "data": leads})
}

func createLead(c *fiber.Ctx) error {
	lead := new(Lead)
	if err := c.BodyParser(lead); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if lead.Code == "" {
		lead.Code = generateCode("L", &Lead{})
	}
	if err := db.Create(lead).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "data": lead})
}

func updateLead(c *fiber.Ctx) error {
	id := c.Params("id")
	lead := new(Lead)
	if err := db.First(lead, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Lead not found"})
	}
	if err := c.BodyParser(lead); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	db.Save(lead)
	return c.JSON(fiber.Map{"success": true, "data": lead})
}

func deleteLead(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.Delete(&Lead{}, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Lead not found"})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Lead deleted"})
}

// ==================== PRODUCT HANDLERS ====================

func getProducts(c *fiber.Ctx) error {
	var products []Product
	db.Order("created_at desc").Find(&products)
	return c.JSON(fiber.Map{"success": true, "data": products})
}

func createProduct(c *fiber.Ctx) error {
	product := new(Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if product.Code == "" {
		product.Code = generateCode("P", &Product{})
	}
	if err := db.Create(product).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "data": product})
}

func updateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	product := new(Product)
	if err := db.First(product, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Product not found"})
	}
	if err := c.BodyParser(product); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	db.Save(product)
	return c.JSON(fiber.Map{"success": true, "data": product})
}

func deleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.Delete(&Product{}, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Product not found"})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Product deleted"})
}

// ==================== CHAT HANDLERS ====================

func getChatMessages(c *fiber.Ctx) error {
	var messages []ChatMessage
	db.Order("created_at desc").Find(&messages)
	return c.JSON(fiber.Map{"success": true, "data": messages})
}

func getChatMessageByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var message ChatMessage
	if err := db.First(&message, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Chat not found"})
	}
	return c.JSON(fiber.Map{"success": true, "data": message})
}

func createChatMessage(c *fiber.Ctx) error {
	message := new(ChatMessage)
	if err := c.BodyParser(message); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if err := db.Create(message).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "data": message})
}

func updateChatMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	message := new(ChatMessage)
	if err := db.First(message, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Chat not found"})
	}
	if err := c.BodyParser(message); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	db.Save(message)
	return c.JSON(fiber.Map{"success": true, "data": message})
}

func assignChat(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		AssignedTo    string `json:"assigned_to"`
		AssignedAgent string `json:"assigned_agent"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	var message ChatMessage
	if err := db.First(&message, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Chat not found"})
	}
	message.AssignedTo = req.AssignedTo
	message.AssignedAgent = req.AssignedAgent
	message.Status = "Assigned"
	db.Save(&message)
	return c.JSON(fiber.Map{"success": true, "data": message})
}

func resolveChat(c *fiber.Ctx) error {
	id := c.Params("id")
	var message ChatMessage
	if err := db.First(&message, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Chat not found"})
	}
	message.Status = "Resolved"
	db.Save(&message)
	return c.JSON(fiber.Map{"success": true, "data": message})
}

// Implement remaining handlers (ChatLabel, QuickReply, BroadcastTemplate, AI, Webhooks)
// Similar pattern as above...

func getChatLabels(c *fiber.Ctx) error {
	var labels []ChatLabel
	db.Find(&labels)
	return c.JSON(fiber.Map{"success": true, "data": labels})
}

func createChatLabel(c *fiber.Ctx) error {
	label := new(ChatLabel)
	if err := c.BodyParser(label); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	db.Create(label)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": label})
}

func updateChatLabel(c *fiber.Ctx) error {
	id := c.Params("id")
	label := new(ChatLabel)
	if err := db.First(label, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Label not found"})
	}
	if err := c.BodyParser(label); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	db.Save(label)
	return c.JSON(fiber.Map{"success": true, "data": label})
}

func deleteChatLabel(c *fiber.Ctx) error {
	id := c.Params("id")
	db.Delete(&ChatLabel{}, id)
	return c.JSON(fiber.Map{"success": true, "message": "Label deleted"})
}

func getQuickReplies(c *fiber.Ctx) error {
	var replies []QuickReply
	db.Find(&replies)
	return c.JSON(fiber.Map{"success": true, "data": replies})
}

func createQuickReply(c *fiber.Ctx) error {
	reply := new(QuickReply)
	if err := c.BodyParser(reply); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	db.Create(reply)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": reply})
}

func updateQuickReply(c *fiber.Ctx) error {
	id := c.Params("id")
	reply := new(QuickReply)
	if err := db.First(reply, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Reply not found"})
	}
	if err := c.BodyParser(reply); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	db.Save(reply)
	return c.JSON(fiber.Map{"success": true, "data": reply})
}

func deleteQuickReply(c *fiber.Ctx) error {
	id := c.Params("id")
	db.Delete(&QuickReply{}, id)
	return c.JSON(fiber.Map{"success": true, "message": "Reply deleted"})
}

func getBroadcastTemplates(c *fiber.Ctx) error {
	var templates []BroadcastTemplate
	db.Find(&templates)
	return c.JSON(fiber.Map{"success": true, "data": templates})
}

func createBroadcastTemplate(c *fiber.Ctx) error {
	template := new(BroadcastTemplate)
	if err := c.BodyParser(template); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	db.Create(template)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": template})
}

func updateBroadcastTemplate(c *fiber.Ctx) error {
	id := c.Params("id")
	template := new(BroadcastTemplate)
	if err := db.First(template, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Template not found"})
	}
	if err := c.BodyParser(template); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	db.Save(template)
	return c.JSON(fiber.Map{"success": true, "data": template})
}

func deleteBroadcastTemplate(c *fiber.Ctx) error {
	id := c.Params("id")
	db.Delete(&BroadcastTemplate{}, id)
	return c.JSON(fiber.Map{"success": true, "message": "Template deleted"})
}

// AI Handlers
func getAIConfigurations(c *fiber.Ctx) error {
	var configs []AIConfiguration
	db.Find(&configs)
	return c.JSON(fiber.Map{"success": true, "data": configs})
}

func createAIConfiguration(c *fiber.Ctx) error {
	config := new(AIConfiguration)
	if err := c.BodyParser(config); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	db.Create(config)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": config})
}

func updateAIConfiguration(c *fiber.Ctx) error {
	id := c.Params("id")
	config := new(AIConfiguration)
	if err := db.First(config, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Config not found"})
	}
	if err := c.BodyParser(config); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	db.Save(config)
	return c.JSON(fiber.Map{"success": true, "data": config})
}

func deleteAIConfiguration(c *fiber.Ctx) error {
	id := c.Params("id")
	db.Delete(&AIConfiguration{}, id)
	return c.JSON(fiber.Map{"success": true, "message": "Config deleted"})
}

func getAIAgents(c *fiber.Ctx) error {
	var agents []AIAgent
	db.Find(&agents)
	return c.JSON(fiber.Map{"success": true, "data": agents})
}

func createAIAgent(c *fiber.Ctx) error {
	agent := new(AIAgent)
	if err := c.BodyParser(agent); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	db.Create(agent)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": agent})
}

func updateAIAgent(c *fiber.Ctx) error {
	id := c.Params("id")
	agent := new(AIAgent)
	if err := db.First(agent, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Agent not found"})
	}
	if err := c.BodyParser(agent); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	db.Save(agent)
	return c.JSON(fiber.Map{"success": true, "data": agent})
}

func deleteAIAgent(c *fiber.Ctx) error {
	id := c.Params("id")
	db.Delete(&AIAgent{}, id)
	return c.JSON(fiber.Map{"success": true, "message": "Agent deleted"})
}

func getConnectedPlatforms(c *fiber.Ctx) error {
	var platforms []ConnectedPlatform
	db.Find(&platforms)
	return c.JSON(fiber.Map{"success": true, "data": platforms})
}

func createConnectedPlatform(c *fiber.Ctx) error {
	platform := new(ConnectedPlatform)
	if err := c.BodyParser(platform); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	db.Create(platform)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": platform})
}

func updateConnectedPlatform(c *fiber.Ctx) error {
	id := c.Params("id")
	platform := new(ConnectedPlatform)
	if err := db.First(platform, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Platform not found"})
	}
	if err := c.BodyParser(platform); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	db.Save(platform)
	return c.JSON(fiber.Map{"success": true, "data": platform})
}

func deleteConnectedPlatform(c *fiber.Ctx) error {
	id := c.Params("id")
	db.Delete(&ConnectedPlatform{}, id)
	return c.JSON(fiber.Map{"success": true, "message": "Platform deleted"})
}

// AI Processing
func processWithAI(c *fiber.Ctx) error {
	var req struct {
		Message string `json:"message"`
		AgentID uint   `json:"agent_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	var agent AIAgent
	if err := db.First(&agent, req.AgentID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "AI Agent not found"})
	}

	var config AIConfiguration
	if err := db.Where("ai_engine = ? AND active = ?", agent.AIEngine, true).First(&config).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "AI Configuration not found"})
	}

	response, tokensUsed, err := callOpenAI(&config, agent.BasicPrompt, req.Message)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "AI processing failed: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"response":    response,
			"tokens_used": tokensUsed,
			"agent":       agent.Name,
		},
	})
}

func suggestReply(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"success": true, "suggestions": []string{"Sample suggestion"}})
}

func getTokenBalances(c *fiber.Ctx) error {
	var balances []TokenBalance
	db.Find(&balances)
	return c.JSON(fiber.Map{"success": true, "data": balances})
}

func resetTokenBalance(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"success": true, "message": "Token balance reset"})
}

func getAnalyticsOverview(c *fiber.Ctx) error {
	var total int64
	db.Model(&ChatMessage{}).Count(&total)
	return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"total_chats": total}})
}

func getAgentPerformance(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"success": true, "data": fiber.Map{}})
}

// Webhook Handlers
func verifyWhatsAppWebhook(c *fiber.Ctx) error {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")
	if mode == "subscribe" && token == os.Getenv("WHATSAPP_VERIFY_TOKEN") {
		return c.SendString(challenge)
	}
	return c.Status(403).SendString("Forbidden")
}

func handleWhatsAppWebhook(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

func handleTelegramWebhook(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

func verifyInstagramWebhook(c *fiber.Ctx) error {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")
	if mode == "subscribe" && token == os.Getenv("INSTAGRAM_VERIFY_TOKEN") {
		return c.SendString(challenge)
	}
	return c.Status(403).SendString("Forbidden")
}

func handleInstagramWebhook(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

// Utilities
func generateCode(prefix string, model interface{}) string {
	var count int64
	db.Model(model).Count(&count)
	return fmt.Sprintf("%s%04d", prefix, count+1)
}

func callOpenAI(config *AIConfiguration, systemPrompt, userMessage string) (string, int, error) {
	aiReq := AIRequest{
		Model: config.Model,
		Messages: []AIMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userMessage},
		},
		Temperature: 0.7,
		MaxTokens:   500,
	}

	reqBody, _ := json.Marshal(aiReq)
	req, _ := http.NewRequest("POST", config.Endpoint, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", 0, fmt.Errorf("API error: %s", string(body))
	}

	var aiResp OpenAIResponse
	json.Unmarshal(body, &aiResp)
	if len(aiResp.Choices) == 0 {
		return "", 0, fmt.Errorf("no response from AI")
	}

	return aiResp.Choices[0].Message.Content, aiResp.Usage.TotalTokens, nil
}

// Middleware
func authMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "No authorization header"})
	}

	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	claims, err := verifyToken(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Invalid token"})
	}

	c.Locals("userID", claims.UserID)
	c.Locals("email", claims.Email)
	c.Locals("role", claims.Role)

	return c.Next()
}

func verifyToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
