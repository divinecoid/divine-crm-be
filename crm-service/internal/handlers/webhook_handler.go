package handlers

import (
	"divine-crm/internal/config"
	"divine-crm/internal/services"
	"log"
	"log/slog"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

type WebhookHandler struct {
	service *services.WebhookService
	config  *config.Config
	logger  *slog.Logger
}

func NewWebhookHandler(service *services.WebhookService, cfg *config.Config) *WebhookHandler {
	return &WebhookHandler{
		service: service,
		config:  cfg,
		logger:  slog.Default(),
	}
}

// ==================== WHATSAPP ====================

func (h *WebhookHandler) VerifyWhatsApp(c *fiber.Ctx) error {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	h.logger.Info("WhatsApp verification request",
		"mode", mode,
		"token_received", token,
		"token_expected", h.config.WhatsApp.VerifyToken,
	)

	if mode == "subscribe" && token == h.config.WhatsApp.VerifyToken {
		h.logger.Info("✅ WhatsApp webhook verified")
		c.Set("Content-Type", "text/plain")
		return c.SendString(challenge)
	}

	h.logger.Error("❌ WhatsApp verification failed")
	return c.Status(fiber.StatusForbidden).SendString("Verification failed")
}

func (h *WebhookHandler) HandleWhatsApp(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("🚨 PANIC in HandleWhatsApp: %v", r)
			log.Printf("Stack: %s", debug.Stack())
		}
	}()

	log.Println("📥 WhatsApp webhook received")

	var payload services.WhatsAppWebhookPayload
	if err := c.BodyParser(&payload); err != nil {
		log.Printf("❌ Failed to parse WhatsApp payload: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := h.service.ProcessWhatsAppWebhook(&payload); err != nil {
		log.Printf("❌ Failed to process WhatsApp webhook: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

// ==================== INSTAGRAM ====================

func (h *WebhookHandler) VerifyInstagram(c *fiber.Ctx) error {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	h.logger.Info("Instagram verification request",
		"mode", mode,
		"token", token,
	)

	if mode == "subscribe" && token == h.config.Instagram.VerifyToken {
		h.logger.Info("✅ Instagram webhook verified")
		c.Set("Content-Type", "text/plain")
		return c.SendString(challenge)
	}

	h.logger.Error("❌ Instagram verification failed")
	return c.Status(fiber.StatusForbidden).SendString("Verification failed")
}

func (h *WebhookHandler) HandleInstagram(c *fiber.Ctx) error {
	log.Println("📥 Instagram webhook received")
	log.Printf("Body: %s", string(c.Body()))

	var payload services.InstagramWebhookPayload
	if err := c.BodyParser(&payload); err != nil {
		log.Printf("❌ Failed to parse Instagram payload: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := h.service.ProcessInstagramWebhook(&payload); err != nil {
		log.Printf("❌ Failed to process Instagram webhook: %v", err)
		// Return 200 anyway to avoid retries
		return c.JSON(fiber.Map{"success": false, "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

// ==================== TELEGRAM ====================

func (h *WebhookHandler) VerifyTelegram(c *fiber.Ctx) error {
	// Telegram doesn't use hub verification
	return c.JSON(fiber.Map{"status": "ok"})
}

func (h *WebhookHandler) HandleTelegram(c *fiber.Ctx) error {
	log.Println("📥 Telegram webhook received")
	log.Printf("Body: %s", string(c.Body()))

	var payload services.TelegramWebhookPayload
	if err := c.BodyParser(&payload); err != nil {
		log.Printf("❌ Failed to parse Telegram payload: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := h.service.ProcessTelegramWebhook(&payload); err != nil {
		log.Printf("❌ Failed to process Telegram webhook: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"ok": true})
}
