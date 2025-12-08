package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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

// ==================== JWT CLAIMS (for token verification only) ====================

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// ==================== OTHER MODELS ====================

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
	Notes         string    `json:"notes" gorm:"type:text"`
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
	ID          uint      `json:"id" gorm:"primaryKey"`
	Code        string    `json:"code" gorm:"uniqueIndex"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Description string    `json:"description" gorm:"type:text"`
	UploadedBy  string    `json:"uploaded_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

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
	Type      string    `json:"type"` // text, image, document
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Broadcast represents a sent broadcast message
type Broadcast struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	TemplateID     uint      `json:"template_id"`
	TemplateName   string    `json:"template_name"`
	Message        string    `json:"message" gorm:"type:text"`
	Channel        string    `json:"channel"`       // WhatsApp, Telegram, All
	TargetFilter   string    `json:"target_filter"` // all, hot, warm, cold
	TotalSent      int       `json:"total_sent"`
	TotalDelivered int       `json:"total_delivered"`
	TotalFailed    int       `json:"total_failed"`
	Status         string    `json:"status" gorm:"default:pending"` // pending, sending, completed, failed
	SentBy         string    `json:"sent_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

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
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"uniqueIndex"`
	AIEngine     string    `json:"ai_engine"`
	Platform     string    `json:"platform" gorm:"default:'All'"` // All, WhatsApp, Telegram, Instagram
	BasicPrompt  string    `json:"basic_prompt" gorm:"type:text"`
	IntroMessage string    `json:"intro_message" gorm:"type:text"` // Introduction message for new contacts
	Active       bool      `json:"active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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

type ConnectedPlatform struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Platform      string    `json:"platform" gorm:"uniqueIndex"`
	PlatformID    string    `json:"platform_id"`
	Token         string    `json:"token"`
	ClientID      string    `json:"client_id"`
	ClientSecret  string    `json:"client_secret"`
	WebhookURL    string    `json:"webhook_url"`
	PhoneNumberID string    `json:"phone_number_id"`
	Active        bool      `json:"active" gorm:"default:true"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// AI Request/Response
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

// ==================== INITIALIZATION ====================

func initDB() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_USER", "postgres"),
			getEnv("DB_PASSWORD", "postgres"),
			getEnv("DB_NAME", "divine_crm"),
			getEnv("DB_SSLMODE", "disable"),
		)
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Note: User model is managed by auth-service
	db.AutoMigrate(
		&Contact{}, &Lead{}, &Product{},
		&ChatMessage{}, &ChatLabel{}, &QuickReply{}, &BroadcastTemplate{}, &Broadcast{},
		&AIConfiguration{}, &AIAgent{}, &TokenBalance{},
		&ConnectedPlatform{},
	)

	log.Println("✅ CRM Service: Database connected and migrated")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ==================== JWT VERIFICATION (tokens from auth-service) ====================

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

// ==================== MIDDLEWARE ====================

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
	c.Locals("username", claims.Username)
	c.Locals("email", claims.Email)
	c.Locals("role", claims.Role)

	return c.Next()
}

func requireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role").(string)
		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}
		return c.Status(403).JSON(fiber.Map{"success": false, "message": "Access denied"})
	}
}

// ==================== MAIN ====================

func main() {
	godotenv.Load()
	initDB()

	// JWT Secret - must match auth-service
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("divine-crm-jwt-secret-key-2024")
	}

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

	// NOTE: Auth routes (/auth/*, /users/*) are handled by auth-service on port 3001

	// Masterdata Routes (Protected - requires valid JWT from auth-service)
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

	// Chat Routes
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

	// AI Routes
	ai := api.Group("/ai", authMiddleware)
	ai.Post("/process", processWithAI)
	ai.Post("/suggest-reply", suggestReply)
	ai.Get("/token-balance", getTokenBalances)
	ai.Get("/debug-prompt", debugAIPrompt) // Debug endpoint

	// Analytics
	analytics := api.Group("/analytics", authMiddleware)
	analytics.Get("/overview", getAnalyticsOverview)

	// Broadcast Routes
	broadcast := api.Group("/broadcasts", authMiddleware)
	broadcast.Get("/", getBroadcasts)
	broadcast.Get("/:id", getBroadcastByID)
	broadcast.Post("/send", sendBroadcast)

	// Webhooks (Public - no auth required)
	webhooks := api.Group("/webhooks")
	webhooks.Get("/whatsapp", verifyWhatsAppWebhook)
	webhooks.Post("/whatsapp", handleWhatsAppWebhook)
	webhooks.Get("/telegram", verifyTelegramWebhook)
	webhooks.Post("/telegram", handleTelegramWebhook)

	port := getEnv("SERVER_PORT", "3002")

	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("📦 Divine CRM Service running on port %s", port)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("⚠️  Auth is handled by auth-service (port 3001)")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("")
	log.Println("📡 Endpoints:")
	log.Println("   /api/v1/masterdata/*     - All masterdata CRUD")
	log.Println("   /api/v1/chats/*          - Chat operations")
	log.Println("   /api/v1/ai/*             - AI processing")
	log.Println("   /api/v1/broadcasts/*     - Broadcast messaging")
	log.Println("   /api/v1/analytics/*      - Analytics")
	log.Println("   /api/v1/webhooks/*       - Webhooks (public)")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	log.Fatal(app.Listen(":" + port))
}

// ==================== CRUD HANDLERS ====================

func getContacts(c *fiber.Ctx) error {
	var contacts []Contact
	db.Order("created_at desc").Find(&contacts)
	return c.JSON(fiber.Map{"success": true, "data": contacts})
}

func getContactByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var contact Contact
	if err := db.First(&contact, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
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
	db.Create(contact)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": contact})
}

func updateContact(c *fiber.Ctx) error {
	id := c.Params("id")
	contact := new(Contact)
	if err := db.First(contact, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(contact)
	db.Save(contact)
	return c.JSON(fiber.Map{"success": true, "data": contact})
}

func deleteContact(c *fiber.Ctx) error {
	db.Delete(&Contact{}, c.Params("id"))
	return c.JSON(fiber.Map{"success": true, "message": "Deleted"})
}

func getLeads(c *fiber.Ctx) error {
	var leads []Lead
	db.Order("created_at desc").Find(&leads)
	return c.JSON(fiber.Map{"success": true, "data": leads})
}

func createLead(c *fiber.Ctx) error {
	lead := new(Lead)
	c.BodyParser(lead)
	if lead.Code == "" {
		lead.Code = generateCode("L", &Lead{})
	}
	db.Create(lead)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": lead})
}

func updateLead(c *fiber.Ctx) error {
	lead := new(Lead)
	if err := db.First(lead, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(lead)
	db.Save(lead)
	return c.JSON(fiber.Map{"success": true, "data": lead})
}

func deleteLead(c *fiber.Ctx) error {
	db.Delete(&Lead{}, c.Params("id"))
	return c.JSON(fiber.Map{"success": true, "message": "Deleted"})
}

func getProducts(c *fiber.Ctx) error {
	var products []Product
	db.Order("created_at desc").Find(&products)
	return c.JSON(fiber.Map{"success": true, "data": products})
}

func createProduct(c *fiber.Ctx) error {
	product := new(Product)
	c.BodyParser(product)
	if product.Code == "" {
		product.Code = generateCode("P", &Product{})
	}
	db.Create(product)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": product})
}

func updateProduct(c *fiber.Ctx) error {
	product := new(Product)
	if err := db.First(product, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(product)
	db.Save(product)
	return c.JSON(fiber.Map{"success": true, "data": product})
}

func deleteProduct(c *fiber.Ctx) error {
	db.Delete(&Product{}, c.Params("id"))
	return c.JSON(fiber.Map{"success": true, "message": "Deleted"})
}

func getChatLabels(c *fiber.Ctx) error {
	var labels []ChatLabel
	db.Find(&labels)
	return c.JSON(fiber.Map{"success": true, "data": labels})
}

func createChatLabel(c *fiber.Ctx) error {
	label := new(ChatLabel)
	c.BodyParser(label)
	db.Create(label)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": label})
}

func updateChatLabel(c *fiber.Ctx) error {
	label := new(ChatLabel)
	if err := db.First(label, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(label)
	db.Save(label)
	return c.JSON(fiber.Map{"success": true, "data": label})
}

func deleteChatLabel(c *fiber.Ctx) error {
	db.Delete(&ChatLabel{}, c.Params("id"))
	return c.JSON(fiber.Map{"success": true, "message": "Deleted"})
}

func getQuickReplies(c *fiber.Ctx) error {
	var replies []QuickReply
	db.Find(&replies)
	return c.JSON(fiber.Map{"success": true, "data": replies})
}

func createQuickReply(c *fiber.Ctx) error {
	reply := new(QuickReply)
	c.BodyParser(reply)
	db.Create(reply)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": reply})
}

func updateQuickReply(c *fiber.Ctx) error {
	reply := new(QuickReply)
	if err := db.First(reply, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(reply)
	db.Save(reply)
	return c.JSON(fiber.Map{"success": true, "data": reply})
}

func deleteQuickReply(c *fiber.Ctx) error {
	db.Delete(&QuickReply{}, c.Params("id"))
	return c.JSON(fiber.Map{"success": true, "message": "Deleted"})
}

func getBroadcastTemplates(c *fiber.Ctx) error {
	var templates []BroadcastTemplate
	db.Find(&templates)
	return c.JSON(fiber.Map{"success": true, "data": templates})
}

func createBroadcastTemplate(c *fiber.Ctx) error {
	template := new(BroadcastTemplate)
	c.BodyParser(template)
	db.Create(template)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": template})
}

func updateBroadcastTemplate(c *fiber.Ctx) error {
	template := new(BroadcastTemplate)
	if err := db.First(template, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(template)
	db.Save(template)
	return c.JSON(fiber.Map{"success": true, "data": template})
}

func deleteBroadcastTemplate(c *fiber.Ctx) error {
	db.Delete(&BroadcastTemplate{}, c.Params("id"))
	return c.JSON(fiber.Map{"success": true, "message": "Deleted"})
}

func getAIConfigurations(c *fiber.Ctx) error {
	var configs []AIConfiguration
	db.Find(&configs)
	return c.JSON(fiber.Map{"success": true, "data": configs})
}

func createAIConfiguration(c *fiber.Ctx) error {
	config := new(AIConfiguration)
	c.BodyParser(config)
	db.Create(config)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": config})
}

func updateAIConfiguration(c *fiber.Ctx) error {
	config := new(AIConfiguration)
	if err := db.First(config, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(config)
	db.Save(config)
	return c.JSON(fiber.Map{"success": true, "data": config})
}

func deleteAIConfiguration(c *fiber.Ctx) error {
	db.Delete(&AIConfiguration{}, c.Params("id"))
	return c.JSON(fiber.Map{"success": true, "message": "Deleted"})
}

func getAIAgents(c *fiber.Ctx) error {
	var agents []AIAgent
	db.Find(&agents)
	return c.JSON(fiber.Map{"success": true, "data": agents})
}

func createAIAgent(c *fiber.Ctx) error {
	agent := new(AIAgent)
	c.BodyParser(agent)
	db.Create(agent)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": agent})
}

func updateAIAgent(c *fiber.Ctx) error {
	agent := new(AIAgent)
	if err := db.First(agent, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(agent)
	db.Save(agent)
	return c.JSON(fiber.Map{"success": true, "data": agent})
}

func deleteAIAgent(c *fiber.Ctx) error {
	db.Delete(&AIAgent{}, c.Params("id"))
	return c.JSON(fiber.Map{"success": true, "message": "Deleted"})
}

func getConnectedPlatforms(c *fiber.Ctx) error {
	var platforms []ConnectedPlatform
	db.Find(&platforms)
	return c.JSON(fiber.Map{"success": true, "data": platforms})
}

func createConnectedPlatform(c *fiber.Ctx) error {
	platform := new(ConnectedPlatform)
	c.BodyParser(platform)
	db.Create(platform)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": platform})
}

func updateConnectedPlatform(c *fiber.Ctx) error {
	platform := new(ConnectedPlatform)
	if err := db.First(platform, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(platform)
	db.Save(platform)
	return c.JSON(fiber.Map{"success": true, "data": platform})
}

func deleteConnectedPlatform(c *fiber.Ctx) error {
	db.Delete(&ConnectedPlatform{}, c.Params("id"))
	return c.JSON(fiber.Map{"success": true, "message": "Deleted"})
}

func getChatMessages(c *fiber.Ctx) error {
	var messages []ChatMessage
	status := c.Query("status")

	query := db.Order("created_at desc")
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}
	query.Find(&messages)
	return c.JSON(fiber.Map{"success": true, "data": messages})
}

func getChatMessageByID(c *fiber.Ctx) error {
	var message ChatMessage
	if err := db.First(&message, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	return c.JSON(fiber.Map{"success": true, "data": message})
}

func createChatMessage(c *fiber.Ctx) error {
	message := new(ChatMessage)
	c.BodyParser(message)
	db.Create(message)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": message})
}

func updateChatMessage(c *fiber.Ctx) error {
	message := new(ChatMessage)
	if err := db.First(message, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(message)
	db.Save(message)
	return c.JSON(fiber.Map{"success": true, "data": message})
}

func assignChat(c *fiber.Ctx) error {
	var req struct {
		AssignedTo    string `json:"assigned_to"`
		AssignedAgent string `json:"assigned_agent"`
	}
	c.BodyParser(&req)
	var message ChatMessage
	if err := db.First(&message, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	message.AssignedTo = req.AssignedTo
	message.AssignedAgent = req.AssignedAgent
	message.Status = "Assigned"
	db.Save(&message)
	return c.JSON(fiber.Map{"success": true, "data": message})
}

func resolveChat(c *fiber.Ctx) error {
	var message ChatMessage
	if err := db.First(&message, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	message.Status = "Resolved"
	db.Save(&message)
	return c.JSON(fiber.Map{"success": true, "data": message})
}

// Get chat statistics
func getChatStats(c *fiber.Ctx) error {
	var unassigned, pending, assigned, resolved, total int64
	db.Model(&ChatMessage{}).Where("status = ?", "Unassigned").Count(&unassigned)
	db.Model(&ChatMessage{}).Where("status = ?", "Pending").Count(&pending)
	db.Model(&ChatMessage{}).Where("status = ?", "Assigned").Count(&assigned)
	db.Model(&ChatMessage{}).Where("status = ?", "Resolved").Count(&resolved)
	db.Model(&ChatMessage{}).Count(&total)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"unassigned": unassigned,
			"pending":    pending,
			"assigned":   assigned,
			"resolved":   resolved,
			"total":      total,
		},
	})
}

// Get conversations grouped by contact
func getChatConversations(c *fiber.Ctx) error {
	status := c.Query("status")

	type ConversationResult struct {
		ContactID   uint   `json:"contact_id"`
		ContactName string `json:"contact_name"`
		Channel     string `json:"channel"`
		LastMessage string `json:"last_message"`
		LastStatus  string `json:"last_status"`
		UnreadCount int    `json:"unread_count"`
		LastUpdated string `json:"last_updated"`
	}

	var results []ConversationResult

	query := `
		SELECT 
			contact_id,
			contact_name,
			channel,
			message as last_message,
			status as last_status,
			0 as unread_count,
			updated_at as last_updated
		FROM chat_messages
		WHERE id IN (
			SELECT MAX(id) FROM chat_messages GROUP BY contact_id
		)
	`

	if status != "" && status != "all" {
		query += " AND status = ?"
		db.Raw(query, status).Scan(&results)
	} else {
		db.Raw(query).Scan(&results)
	}

	return c.JSON(fiber.Map{"success": true, "data": results})
}

// Get chats by contact ID
func getChatsByContact(c *fiber.Ctx) error {
	contactId := c.Params("contactId")
	var messages []ChatMessage
	db.Where("contact_id = ?", contactId).Order("created_at asc").Find(&messages)
	return c.JSON(fiber.Map{"success": true, "data": messages})
}

// Takeover chat - human agent takes over from AI
func takeoverChat(c *fiber.Ctx) error {
	var req struct {
		AgentName string `json:"agent_name"`
	}
	c.BodyParser(&req)

	var message ChatMessage
	if err := db.First(&message, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}

	message.Status = "Assigned"
	message.AssignedTo = "Human"
	message.AssignedAgent = req.AgentName
	db.Save(&message)

	// Update contact's last agent
	var contact Contact
	if db.First(&contact, message.ContactID).Error == nil {
		contact.LastAgent = req.AgentName
		contact.LastAgentType = "Human"
		db.Save(&contact)
	}

	log.Printf("👤 Chat takeover: Agent %s took over chat #%d", req.AgentName, message.ID)
	return c.JSON(fiber.Map{"success": true, "data": message, "message": "Takeover successful"})
}

// Back to AI - return chat to AI handling
func backToAIChat(c *fiber.Ctx) error {
	var message ChatMessage
	if err := db.First(&message, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}

	message.Status = "Unassigned"
	message.AssignedTo = "AI Bot"
	message.AssignedAgent = "AI Assistant"
	db.Save(&message)

	// Update contact's last agent
	var contact Contact
	if db.First(&contact, message.ContactID).Error == nil {
		contact.LastAgent = "AI"
		contact.LastAgentType = "Bot"
		db.Save(&contact)
	}

	log.Printf("🤖 Chat returned to AI: chat #%d", message.ID)
	return c.JSON(fiber.Map{"success": true, "data": message, "message": "Returned to AI"})
}

// Set chat to pending status
func setPendingChat(c *fiber.Ctx) error {
	var req struct {
		FollowupDays int `json:"followup_days"`
	}
	c.BodyParser(&req)

	if req.FollowupDays == 0 {
		req.FollowupDays = 1
	}

	var message ChatMessage
	if err := db.First(&message, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}

	message.Status = "Pending"
	db.Save(&message)

	log.Printf("⏰ Chat set to pending: chat #%d, follow-up in %d days", message.ID, req.FollowupDays)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    message,
		"message": fmt.Sprintf("Set to pending, follow-up in %d days", req.FollowupDays),
	})
}

// Send reply to customer
func sendReplyChat(c *fiber.Ctx) error {
	var req struct {
		Message   string `json:"message"`
		Channel   string `json:"channel"`
		ChannelID string `json:"channel_id"`
	}
	c.BodyParser(&req)

	// Get original chat
	var originalMsg ChatMessage
	if err := db.First(&originalMsg, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Chat not found"})
	}

	// Get current user from JWT
	user := c.Locals("user")
	agentName := "Agent"
	if claims, ok := user.(*JWTClaims); ok {
		agentName = claims.FullName
	}

	// Create new message record
	newMsg := ChatMessage{
		ContactID:     originalMsg.ContactID,
		ContactName:   originalMsg.ContactName,
		Message:       "", // Customer's message would be empty for agent reply
		Response:      req.Message,
		Channel:       req.Channel,
		Status:        "Assigned",
		AssignedTo:    "Human",
		AssignedAgent: agentName,
		TokensUsed:    0,
	}
	db.Create(&newMsg)

	// Send message via appropriate channel
	var sendErr error
	switch req.Channel {
	case "WhatsApp":
		sendErr = sendWhatsAppMessage(req.ChannelID, req.Message)
	case "Telegram":
		sendErr = sendTelegramMessage(req.ChannelID, req.Message)
	case "Instagram":
		// Instagram send would go here
		log.Printf("📸 Instagram reply to %s: %s", req.ChannelID, req.Message)
	default:
		log.Printf("📨 Reply to %s via %s: %s", req.ChannelID, req.Channel, req.Message)
	}

	if sendErr != nil {
		log.Printf("❌ Failed to send message: %v", sendErr)
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to send message"})
	}

	// Update contact
	var contact Contact
	if db.First(&contact, originalMsg.ContactID).Error == nil {
		contact.LastAgent = agentName
		contact.LastAgentType = "Human"
		contact.LastContact = time.Now()
		db.Save(&contact)
	}

	log.Printf("📤 Reply sent by %s to %s via %s", agentName, req.ChannelID, req.Channel)
	return c.JSON(fiber.Map{"success": true, "data": newMsg, "message": "Reply sent"})
}

func processWithAI(c *fiber.Ctx) error {
	var req struct {
		Message string `json:"message"`
		AgentID uint   `json:"agent_id"`
	}
	c.BodyParser(&req)

	var agent AIAgent
	if err := db.Where("active = ?", true).First(&agent).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "No AI Agent"})
	}

	var config AIConfiguration
	if err := db.Where("ai_engine = ? AND active = ?", agent.AIEngine, true).First(&config).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "No AI Config"})
	}

	response, tokens, err := callOpenAI(&config, agent.BasicPrompt, req.Message)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"response":    response,
			"tokens_used": tokens,
			"agent":       agent.Name,
		},
	})
}

func suggestReply(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"success": true, "suggestions": []string{}})
}

func getTokenBalances(c *fiber.Ctx) error {
	var balances []TokenBalance
	db.Find(&balances)
	return c.JSON(fiber.Map{"success": true, "data": balances})
}

func getAnalyticsOverview(c *fiber.Ctx) error {
	var contacts, messages, products, resolved, pending int64
	db.Model(&Contact{}).Count(&contacts)
	db.Model(&ChatMessage{}).Count(&messages)
	db.Model(&Product{}).Count(&products)
	db.Model(&ChatMessage{}).Where("status = ?", "Resolved").Count(&resolved)
	db.Model(&ChatMessage{}).Where("status IN ?", []string{"Pending", "Unassigned"}).Count(&pending)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"total_contacts": contacts,
			"total_messages": messages,
			"total_products": products,
			"resolved_chats": resolved,
			"pending_chats":  pending,
		},
	})
}

// Debug endpoint to see AI prompt with products
func debugAIPrompt(c *fiber.Ctx) error {
	// Get active AI agent
	var agent AIAgent
	if db.Where("active = ?", true).First(&agent).Error != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "No active AI agent"})
	}

	// Get all products
	var products []Product
	db.Find(&products)

	var productsWithStock []Product
	db.Where("stock > 0").Find(&productsWithStock)

	// Build the full prompt
	fullPrompt := buildAIPromptWithProducts(agent.BasicPrompt)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"agent_name":          agent.Name,
			"agent_basic_prompt":  agent.BasicPrompt,
			"total_products":      len(products),
			"products_with_stock": len(productsWithStock),
			"products":            products,
			"full_ai_prompt":      fullPrompt,
			"prompt_length":       len(fullPrompt),
		},
	})
}

// ==================== BROADCAST FUNCTIONS ====================

// Get all broadcasts
func getBroadcasts(c *fiber.Ctx) error {
	var broadcasts []Broadcast
	db.Order("created_at DESC").Find(&broadcasts)
	return c.JSON(fiber.Map{"success": true, "data": broadcasts})
}

// Get broadcast by ID
func getBroadcastByID(c *fiber.Ctx) error {
	var broadcast Broadcast
	if err := db.First(&broadcast, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	return c.JSON(fiber.Map{"success": true, "data": broadcast})
}

// Send broadcast to contacts
func sendBroadcast(c *fiber.Ctx) error {
	var req struct {
		TemplateID   uint   `json:"template_id"`
		Message      string `json:"message"`       // Custom message or from template
		Channel      string `json:"channel"`       // WhatsApp, Telegram, All
		TargetFilter string `json:"target_filter"` // all, hot, warm, cold
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid request"})
	}

	// Get current user
	user := c.Locals("user")
	sentBy := "System"
	if claims, ok := user.(*JWTClaims); ok {
		sentBy = claims.FullName
	}

	// Get template if provided
	templateName := "Custom"
	if req.TemplateID > 0 {
		var template BroadcastTemplate
		if db.First(&template, req.TemplateID).Error == nil {
			templateName = template.Name
			if req.Message == "" {
				req.Message = template.Content
			}
		}
	}

	if req.Message == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Message is required"})
	}

	// Get contacts based on filter
	var contacts []Contact
	query := db

	// Filter by temperature
	switch req.TargetFilter {
	case "hot":
		query = query.Where("temperature = ?", TempHot)
	case "warm":
		query = query.Where("temperature = ?", TempWarm)
	case "cold":
		query = query.Where("temperature = ?", TempCold)
		// "all" doesn't need filter
	}

	// Filter by channel
	if req.Channel != "All" && req.Channel != "" {
		query = query.Where("channel = ?", req.Channel)
	}

	query.Find(&contacts)

	if len(contacts) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "No contacts match the filter"})
	}

	// Create broadcast record
	broadcast := Broadcast{
		TemplateID:     req.TemplateID,
		TemplateName:   templateName,
		Message:        req.Message,
		Channel:        req.Channel,
		TargetFilter:   req.TargetFilter,
		TotalSent:      0,
		TotalDelivered: 0,
		TotalFailed:    0,
		Status:         "sending",
		SentBy:         sentBy,
	}
	db.Create(&broadcast)

	// Send messages in goroutine
	go func() {
		sent := 0
		delivered := 0
		failed := 0

		for _, contact := range contacts {
			var err error
			switch contact.Channel {
			case "WhatsApp":
				err = sendWhatsAppMessage(contact.ChannelID, req.Message)
			case "Telegram":
				err = sendTelegramMessage(contact.ChannelID, req.Message)
			default:
				msgPreview := req.Message
				if len(msgPreview) > 50 {
					msgPreview = msgPreview[:50]
				}
				log.Printf("📨 Broadcast to %s via %s: %s", contact.Name, contact.Channel, msgPreview)
			}

			sent++
			if err != nil {
				failed++
				log.Printf("❌ Broadcast failed to %s: %v", contact.Name, err)
			} else {
				delivered++
				log.Printf("✅ Broadcast sent to %s via %s", contact.Name, contact.Channel)
			}

			// Small delay to avoid rate limiting
			time.Sleep(100 * time.Millisecond)
		}

		// Update broadcast record
		broadcast.TotalSent = sent
		broadcast.TotalDelivered = delivered
		broadcast.TotalFailed = failed
		broadcast.Status = "completed"
		db.Save(&broadcast)

		log.Printf("📢 Broadcast completed: %d sent, %d delivered, %d failed", sent, delivered, failed)
	}()

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Broadcasting to %d contacts", len(contacts)),
		"data":    broadcast,
	})
}

func verifyWhatsAppWebhook(c *fiber.Ctx) error {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")
	verifyToken := os.Getenv("WHATSAPP_VERIFY_TOKEN")
	if verifyToken == "" {
		verifyToken = "divine-crm-webhook-2024"
	}
	if mode == "subscribe" && token == verifyToken {
		return c.SendString(challenge)
	}
	return c.Status(403).SendString("Forbidden")
}

func handleWhatsAppWebhook(c *fiber.Ctx) error {
	var webhook WhatsAppWebhook
	c.BodyParser(&webhook)

	if len(webhook.Entry) == 0 || len(webhook.Entry[0].Changes) == 0 {
		return c.JSON(fiber.Map{"status": "ok"})
	}

	value := webhook.Entry[0].Changes[0].Value
	if len(value.Messages) == 0 {
		return c.JSON(fiber.Map{"status": "ok"})
	}

	message := value.Messages[0]
	senderPhone := message.From
	messageText := message.Text.Body
	senderName := ""
	if len(value.Contacts) > 0 {
		senderName = value.Contacts[0].Profile.Name
	}

	var contact Contact
	isNewContact := false
	if db.Where("channel = ? AND channel_id = ?", "WhatsApp", senderPhone).First(&contact).Error != nil {
		// New contact - set initial temperature as Cold
		contact = Contact{
			Code: generateCode("C", &Contact{}), Channel: "WhatsApp", ChannelID: senderPhone,
			Name: senderName, Temperature: TempCold, FirstContact: time.Now(), LastContact: time.Now(),
		}
		db.Create(&contact)
		isNewContact = true
		log.Printf("📱 New WhatsApp contact: %s", senderName)
	} else {
		contact.LastContact = time.Now()
		db.Save(&contact)
	}

	// Check if contact has active human agent assigned - check LAST message status
	var latestChat ChatMessage
	humanHandled := false
	if db.Where("contact_id = ?", contact.ID).Order("created_at DESC").First(&latestChat).Error == nil {
		// Only consider human handled if the LAST message was assigned to human (within 24 hours)
		if latestChat.AssignedTo == "Human" && time.Since(latestChat.CreatedAt) < 24*time.Hour {
			humanHandled = true
			log.Printf("👤 Contact %s is being handled by human agent %s - skipping AI", contact.Name, latestChat.AssignedAgent)
		}
	}

	if humanHandled {
		// Just save the message for human agent to see, don't auto-reply
		chatMsg := ChatMessage{
			ContactID: contact.ID, ContactName: contact.Name, Message: messageText,
			Response: "", Status: "Assigned", AssignedTo: "Human", AssignedAgent: latestChat.AssignedAgent,
			Channel: "WhatsApp", TokensUsed: 0,
		}
		db.Create(&chatMsg)
		log.Printf("📥 New message from %s saved for human agent %s", contact.Name, latestChat.AssignedAgent)
		return c.JSON(fiber.Map{"status": "ok", "handled_by": "human"})
	}

	var agent AIAgent
	// Try to find agent specific for WhatsApp first, then fallback to "All"
	if db.Where("platform = ? AND active = ?", "WhatsApp", true).First(&agent).Error != nil {
		if db.Where("platform = ? AND active = ?", "All", true).First(&agent).Error != nil {
			if db.Where("active = ?", true).First(&agent).Error != nil {
				return c.JSON(fiber.Map{"status": "ok"})
			}
		}
	}
	log.Printf("🤖 WhatsApp using AI Agent: %s (Engine: %s)", agent.Name, agent.AIEngine)

	var config AIConfiguration
	if db.Where("ai_engine = ? AND active = ?", agent.AIEngine, true).First(&config).Error != nil {
		return c.JSON(fiber.Map{"status": "ok"})
	}

	// Send intro message for new contacts
	if isNewContact && agent.IntroMessage != "" {
		introMsg := strings.ReplaceAll(agent.IntroMessage, "{name}", senderName)
		introMsg = strings.ReplaceAll(introMsg, "{agent_name}", agent.Name)
		sendWhatsAppMessage(senderPhone, introMsg)
		log.Printf("👋 Sent intro message to new contact %s", senderName)

		// Save intro as first chat message
		introChatMsg := ChatMessage{
			ContactID: contact.ID, ContactName: contact.Name, Message: "(New Contact)",
			Response: introMsg, Status: "Unassigned", AssignedTo: "AI Bot", AssignedAgent: agent.Name,
			Channel: "WhatsApp", TokensUsed: 0,
		}
		db.Create(&introChatMsg)
	}

	// Build AI prompt with products
	systemPrompt := buildAIPromptWithProducts(agent.BasicPrompt)
	aiResponse, tokens, _ := callOpenAI(&config, systemPrompt, messageText)

	chatMsg := ChatMessage{
		ContactID: contact.ID, ContactName: contact.Name, Message: messageText,
		Response: aiResponse, Status: "Unassigned", AssignedTo: "AI Bot", AssignedAgent: agent.Name,
		Channel: "WhatsApp", TokensUsed: tokens,
	}
	db.Create(&chatMsg)

	// Update temperature with AI analysis (after creating chat message)
	updateContactTemperatureWithAI(&contact)
	log.Printf("🌡️ Contact %s temperature: %s %s", contact.Name, contact.Temperature, getTemperatureEmoji(contact.Temperature))

	// Update contact's last agent
	contact.LastAgent = agent.Name
	contact.LastAgentType = "AI"
	db.Save(&contact)

	sendWhatsAppMessage(senderPhone, aiResponse)
	return c.JSON(fiber.Map{"status": "ok"})
}

func verifyTelegramWebhook(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

func handleTelegramWebhook(c *fiber.Ctx) error {
	var update map[string]interface{}
	c.BodyParser(&update)

	message, ok := update["message"].(map[string]interface{})
	if !ok {
		return c.JSON(fiber.Map{"status": "ok"})
	}

	messageText, _ := message["text"].(string)
	from, _ := message["from"].(map[string]interface{})
	senderID := fmt.Sprintf("%.0f", from["id"].(float64))
	senderName, _ := from["first_name"].(string)

	var contact Contact
	isNewContact := false
	if db.Where("channel = ? AND channel_id = ?", "Telegram", senderID).First(&contact).Error != nil {
		// New contact - set initial temperature as Cold
		contact = Contact{
			Code: generateCode("C", &Contact{}), Channel: "Telegram", ChannelID: senderID,
			Name: senderName, Temperature: TempCold, FirstContact: time.Now(), LastContact: time.Now(),
		}
		db.Create(&contact)
		isNewContact = true
		log.Printf("📱 New Telegram contact: %s", senderName)
	} else {
		contact.LastContact = time.Now()
		db.Save(&contact)
	}

	// Check if contact has active human agent assigned - check LAST message status
	var latestChat ChatMessage
	humanHandled := false
	if db.Where("contact_id = ?", contact.ID).Order("created_at DESC").First(&latestChat).Error == nil {
		// Only consider human handled if the LAST message was assigned to human (within 24 hours)
		if latestChat.AssignedTo == "Human" && time.Since(latestChat.CreatedAt) < 24*time.Hour {
			humanHandled = true
			log.Printf("👤 Contact %s is being handled by human agent %s - skipping AI", contact.Name, latestChat.AssignedAgent)
		}
	}

	if humanHandled {
		// Just save the message for human agent to see, don't auto-reply
		chatMsg := ChatMessage{
			ContactID: contact.ID, ContactName: contact.Name, Message: messageText,
			Response: "", Status: "Assigned", AssignedTo: "Human", AssignedAgent: latestChat.AssignedAgent,
			Channel: "Telegram", TokensUsed: 0,
		}
		db.Create(&chatMsg)
		log.Printf("📥 New message from %s saved for human agent %s", contact.Name, latestChat.AssignedAgent)
		return c.JSON(fiber.Map{"status": "ok", "handled_by": "human"})
	}

	var agent AIAgent
	// Try to find agent specific for Telegram first, then fallback to "All"
	if db.Where("platform = ? AND active = ?", "Telegram", true).First(&agent).Error != nil {
		if db.Where("platform = ? AND active = ?", "All", true).First(&agent).Error != nil {
			if db.Where("active = ?", true).First(&agent).Error != nil {
				return c.JSON(fiber.Map{"status": "ok"})
			}
		}
	}
	log.Printf("🤖 Telegram using AI Agent: %s (Engine: %s)", agent.Name, agent.AIEngine)

	var config AIConfiguration
	if db.Where("ai_engine = ? AND active = ?", agent.AIEngine, true).First(&config).Error != nil {
		return c.JSON(fiber.Map{"status": "ok"})
	}

	// Send intro message for new contacts
	if isNewContact && agent.IntroMessage != "" {
		introMsg := strings.ReplaceAll(agent.IntroMessage, "{name}", senderName)
		introMsg = strings.ReplaceAll(introMsg, "{agent_name}", agent.Name)
		sendTelegramMessage(senderID, introMsg)
		log.Printf("👋 Sent intro message to new contact %s", senderName)

		// Save intro as first chat message
		introChatMsg := ChatMessage{
			ContactID: contact.ID, ContactName: contact.Name, Message: "(New Contact)",
			Response: introMsg, Status: "Unassigned", AssignedTo: "AI Bot", AssignedAgent: agent.Name,
			Channel: "Telegram", TokensUsed: 0,
		}
		db.Create(&introChatMsg)
	}

	// Build AI prompt with products
	systemPrompt := buildAIPromptWithProducts(agent.BasicPrompt)
	aiResponse, tokens, _ := callOpenAI(&config, systemPrompt, messageText)

	chatMsg := ChatMessage{
		ContactID: contact.ID, ContactName: contact.Name, Message: messageText,
		Response: aiResponse, Status: "Unassigned", AssignedTo: "AI Bot", AssignedAgent: agent.Name,
		Channel: "Telegram", TokensUsed: tokens,
	}
	db.Create(&chatMsg)

	// Update temperature with AI analysis (after creating chat message)
	updateContactTemperatureWithAI(&contact)
	log.Printf("🌡️ Contact %s temperature: %s %s", contact.Name, contact.Temperature, getTemperatureEmoji(contact.Temperature))

	// Update contact's last agent
	contact.LastAgent = agent.Name
	contact.LastAgentType = "AI"
	db.Save(&contact)

	sendTelegramMessage(senderID, aiResponse)
	return c.JSON(fiber.Map{"status": "ok"})
}

// ==================== AI-BASED TEMPERATURE DETECTION ====================

// Temperature constants
const (
	TempHot  = "Hot"  // Ready to buy
	TempWarm = "Warm" // Interested but not ready
	TempCold = "Cold" // Just browsing/asking
)

// analyzeTemperatureWithAI uses AI to analyze last 3 messages and determine temperature
func analyzeTemperatureWithAI(contactID uint) string {
	// Get last 3 messages for this contact
	var messages []ChatMessage
	db.Where("contact_id = ?", contactID).Order("created_at desc").Limit(3).Find(&messages)

	if len(messages) == 0 {
		return TempCold
	}

	// Build conversation context
	var conversationText string
	for i := len(messages) - 1; i >= 0; i-- {
		conversationText += fmt.Sprintf("Customer: %s\n", messages[i].Message)
		if messages[i].Response != "" {
			conversationText += fmt.Sprintf("Response: %s\n", messages[i].Response)
		}
	}

	// Get active AI configuration
	var config AIConfiguration
	if db.Where("active = ?", true).First(&config).Error != nil {
		log.Println("⚠️ No AI config for temperature analysis, using fallback")
		return analyzeTemperatureFallback(messages[0].Message)
	}

	// AI prompt for temperature analysis
	systemPrompt := `You are a sales lead analyzer. Analyze the customer conversation and determine their buying intent temperature.

Respond with ONLY one word: Hot, Warm, or Cold

Definitions:
- Hot: Customer shows strong buying intent (wants to buy, order, checkout, asks for payment details, says "deal", "fix order", etc.)
- Warm: Customer shows interest but not ready to buy (asks about price, discount, product details, stock, shipping, compares options)
- Cold: Customer is just browsing or greeting (says hello, asks general questions, not interested, cancels, says "later", "think about it")

Analyze the entire conversation context, focusing on the LATEST messages to determine current intent.`

	userMessage := fmt.Sprintf("Analyze this conversation and determine the customer's buying temperature:\n\n%s\n\nRespond with only: Hot, Warm, or Cold", conversationText)

	response, _, err := callOpenAI(&config, systemPrompt, userMessage)
	if err != nil {
		log.Printf("⚠️ AI temperature analysis failed: %v, using fallback", err)
		return analyzeTemperatureFallback(messages[0].Message)
	}

	// Parse response
	response = strings.TrimSpace(strings.ToLower(response))
	switch {
	case strings.Contains(response, "hot"):
		return TempHot
	case strings.Contains(response, "warm"):
		return TempWarm
	default:
		return TempCold
	}
}

// analyzeTemperatureFallback - simple keyword fallback when AI unavailable
func analyzeTemperatureFallback(message string) string {
	lowerMsg := strings.ToLower(message)

	// Hot keywords
	hotKeywords := []string{"beli", "order", "pesan", "checkout", "bayar", "transfer", "deal", "fix", "jadi", "kirim", "alamat", "buy", "purchase"}
	for _, kw := range hotKeywords {
		if strings.Contains(lowerMsg, kw) {
			return TempHot
		}
	}

	// Warm keywords
	warmKeywords := []string{"harga", "berapa", "price", "diskon", "promo", "stock", "stok", "ready", "ongkir", "cod", "tertarik", "interested"}
	for _, kw := range warmKeywords {
		if strings.Contains(lowerMsg, kw) {
			return TempWarm
		}
	}

	return TempCold
}

// updateContactTemperatureWithAI updates contact temperature using AI analysis
func updateContactTemperatureWithAI(contact *Contact) {
	newTemp := analyzeTemperatureWithAI(contact.ID)
	oldTemp := contact.Temperature

	// Temperature can only go up (Cold -> Warm -> Hot), never down
	shouldUpdate := false

	switch oldTemp {
	case TempCold, "":
		if newTemp == TempWarm || newTemp == TempHot {
			shouldUpdate = true
		}
	case TempWarm:
		if newTemp == TempHot {
			shouldUpdate = true
		}
	case TempHot:
		shouldUpdate = false // Already at highest
	default:
		shouldUpdate = true
	}

	if shouldUpdate && newTemp != oldTemp {
		log.Printf("🔥 AI Temperature Upgrade: %s -> %s (contact: %s)", oldTemp, newTemp, contact.Name)
		contact.Temperature = newTemp
		db.Save(contact)
	} else {
		log.Printf("🌡️ AI Temperature: %s stays at %s (contact: %s)", contact.Name, contact.Temperature, contact.Name)
	}
}

// getTemperatureEmoji returns emoji for temperature
func getTemperatureEmoji(temp string) string {
	switch temp {
	case TempHot:
		return "🔥"
	case TempWarm:
		return "☀️"
	case TempCold:
		return "❄️"
	default:
		return "❓"
	}
}

// ==================== AI WITH PRODUCTS ====================

// buildAIPromptWithProducts builds AI system prompt with product catalog
func buildAIPromptWithProducts(basePrompt string) string {
	// Get all products from database
	var products []Product
	db.Where("stock > 0").Find(&products)

	log.Printf("📦 Building AI prompt with %d products from database", len(products))

	if len(products) == 0 {
		log.Printf("⚠️ No products found in database!")
		return basePrompt + "\n\nIMPORTANT: No products are currently available in the catalog. If customer asks about products, inform them that no products are available at the moment."
	}

	// Build product catalog with VERY explicit instructions
	productCatalog := "\n\n" + strings.Repeat("=", 50) + "\n"
	productCatalog += "⚠️ CRITICAL: OFFICIAL PRODUCT CATALOG - USE ONLY THESE PRODUCTS ⚠️\n"
	productCatalog += strings.Repeat("=", 50) + "\n\n"
	productCatalog += "YOU MUST ONLY MENTION PRODUCTS FROM THIS LIST. DO NOT INVENT OR HALLUCINATE ANY OTHER PRODUCTS.\n"
	productCatalog += "If a product is not in this list, it does not exist.\n\n"

	for _, p := range products {
		productCatalog += fmt.Sprintf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		productCatalog += fmt.Sprintf("📦 PRODUCT: %s\n", p.Name)
		productCatalog += fmt.Sprintf("   🔖 Code: %s\n", p.Code)
		productCatalog += fmt.Sprintf("   💰 Price: Rp %s\n", formatPrice(p.Price))
		productCatalog += fmt.Sprintf("   📊 Stock: %d units available\n", p.Stock)
		if p.Description != "" {
			productCatalog += fmt.Sprintf("   📝 Description: %s\n", p.Description)
		}
		log.Printf("   - Added product: %s (Rp %s, Stock: %d)", p.Name, formatPrice(p.Price), p.Stock)
	}

	productCatalog += fmt.Sprintf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
	productCatalog += strings.Repeat("=", 50) + "\n"
	productCatalog += "END OF OFFICIAL PRODUCT CATALOG\n"
	productCatalog += strings.Repeat("=", 50) + "\n\n"

	productCatalog += "🚨 STRICT RULES FOR PRODUCT INQUIRIES:\n"
	productCatalog += "1. ONLY mention products listed above - NO EXCEPTIONS\n"
	productCatalog += "2. Use the EXACT prices shown above\n"
	productCatalog += "3. Use the EXACT stock numbers shown above\n"
	productCatalog += "4. If customer asks about a product NOT in this list, say 'Maaf, produk tersebut tidak tersedia'\n"
	productCatalog += "5. NEVER invent products like 'Divine CRM Basic/Pro/Enterprise' - those do NOT exist\n"
	productCatalog += "6. When listing products, ONLY list the products above\n\n"

	return basePrompt + productCatalog
}

// formatPrice formats price with thousand separators
func formatPrice(price float64) string {
	// Simple formatting for Indonesian Rupiah
	priceStr := fmt.Sprintf("%.0f", price)
	n := len(priceStr)
	if n <= 3 {
		return priceStr
	}

	var result strings.Builder
	for i, c := range priceStr {
		if i > 0 && (n-i)%3 == 0 {
			result.WriteRune('.')
		}
		result.WriteRune(c)
	}
	return result.String()
}

// ==================== CHAT LABEL HANDLERS ====================

// addChatLabel adds a label to a chat message
func addChatLabel(c *fiber.Ctx) error {
	chatID := c.Params("id")

	var req struct {
		LabelID uint `json:"label_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid request"})
	}

	// Get chat message
	var chatMsg ChatMessage
	if err := db.First(&chatMsg, chatID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Chat not found"})
	}

	// Get label
	var label ChatLabel
	if err := db.First(&label, req.LabelID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Label not found"})
	}

	// Parse existing labels
	var labels []uint
	if chatMsg.Labels != "" {
		json.Unmarshal([]byte(chatMsg.Labels), &labels)
	}

	// Check if label already exists
	for _, l := range labels {
		if l == req.LabelID {
			return c.JSON(fiber.Map{"success": true, "message": "Label already added", "data": chatMsg})
		}
	}

	// Add new label
	labels = append(labels, req.LabelID)
	labelsJSON, _ := json.Marshal(labels)
	chatMsg.Labels = string(labelsJSON)
	db.Save(&chatMsg)

	log.Printf("🏷️ Label '%s' added to chat #%d", label.Label, chatMsg.ID)

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Label '%s' added", label.Label),
		"data":    chatMsg,
	})
}

// removeChatLabel removes a label from a chat message
func removeChatLabel(c *fiber.Ctx) error {
	chatID := c.Params("id")
	labelID := c.Params("labelId")

	// Get chat message
	var chatMsg ChatMessage
	if err := db.First(&chatMsg, chatID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Chat not found"})
	}

	// Parse existing labels
	var labels []uint
	if chatMsg.Labels != "" {
		json.Unmarshal([]byte(chatMsg.Labels), &labels)
	}

	// Remove label
	var labelIDUint uint
	fmt.Sscanf(labelID, "%d", &labelIDUint)

	var newLabels []uint
	removed := false
	for _, l := range labels {
		if l != labelIDUint {
			newLabels = append(newLabels, l)
		} else {
			removed = true
		}
	}

	if !removed {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Label not found on this chat"})
	}

	// Update labels
	if len(newLabels) > 0 {
		labelsJSON, _ := json.Marshal(newLabels)
		chatMsg.Labels = string(labelsJSON)
	} else {
		chatMsg.Labels = ""
	}
	db.Save(&chatMsg)

	log.Printf("🏷️ Label removed from chat #%d", chatMsg.ID)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Label removed",
		"data":    chatMsg,
	})
}

// getChatLabelsForMessage returns labels for a specific chat message
func getChatLabelsForMessage(chatMsg *ChatMessage) []ChatLabel {
	var labels []ChatLabel

	if chatMsg.Labels == "" {
		return labels
	}

	var labelIDs []uint
	json.Unmarshal([]byte(chatMsg.Labels), &labelIDs)

	if len(labelIDs) > 0 {
		db.Where("id IN ?", labelIDs).Find(&labels)
	}

	return labels
}

// Utilities
func generateCode(prefix string, model interface{}) string {
	var count int64
	db.Model(model).Count(&count)
	return fmt.Sprintf("%s%04d", prefix, count+1)
}

func callOpenAI(config *AIConfiguration, systemPrompt, userMessage string) (string, int, error) {
	aiReq := AIRequest{
		Model:       config.Model,
		Messages:    []AIMessage{{Role: "system", Content: systemPrompt}, {Role: "user", Content: userMessage}},
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
		return "", 0, fmt.Errorf("AI error: %s", string(body))
	}

	var aiResp OpenAIResponse
	json.Unmarshal(body, &aiResp)
	if len(aiResp.Choices) == 0 {
		return "", 0, fmt.Errorf("no response")
	}

	return aiResp.Choices[0].Message.Content, aiResp.Usage.TotalTokens, nil
}

func sendWhatsAppMessage(to, message string) error {
	var platform ConnectedPlatform
	if db.Where("platform = ? AND active = ?", "WhatsApp", true).First(&platform).Error != nil {
		return fmt.Errorf("WhatsApp not configured")
	}

	url := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/messages", platform.PhoneNumberID)
	payload := map[string]interface{}{
		"messaging_product": "whatsapp", "recipient_type": "individual", "to": to, "type": "text",
		"text": map[string]string{"body": message},
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+platform.Token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func sendTelegramMessage(chatID, message string) error {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		return fmt.Errorf("no bot token")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	payload := map[string]interface{}{"chat_id": chatID, "text": message}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
