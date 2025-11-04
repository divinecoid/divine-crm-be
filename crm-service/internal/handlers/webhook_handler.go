package handlers

import (
	"divine-crm/internal/config"
	"divine-crm/internal/services"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

type WebhookHandler struct {
	service *services.WebhookService
	config  *config.Config
}

func NewWebhookHandler(service *services.WebhookService, cfg *config.Config) *WebhookHandler {
	return &WebhookHandler{
		service: service,
		config:  cfg,
	}
}

func (h *WebhookHandler) VerifyWhatsApp(c *fiber.Ctx) error {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == h.config.WhatsApp.VerifyToken {
		return c.Status(fiber.StatusOK).SendString(challenge)
	}

	return c.Status(fiber.StatusForbidden).SendString("Forbidden")
}

func (h *WebhookHandler) HandleWhatsApp(c *fiber.Ctx) error {
	// ✅ Add panic recovery with detailed stack trace
	defer func() {
		if r := recover(); r != nil {
			log.Printf("🚨 PANIC in HandleWhatsApp!")
			log.Printf("Error: %v", r)
			log.Printf("Stack trace:\n%s", debug.Stack())

			c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   fmt.Sprintf("panic: %v", r),
				"details": "Check server logs for stack trace",
			})
		}
	}()

	// Log request
	log.Printf("📥 Received webhook request")
	log.Printf("Body: %s", string(c.Body()))

	// ✅ Check if handler is nil
	if h == nil {
		log.Printf("❌ Handler is nil!")
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "handler is nil",
		})
	}

	// ✅ Check if service is nil
	if h.service == nil {
		log.Printf("❌ Service is nil!")
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "webhook service is nil",
		})
	}

	var payload services.WhatsAppWebhookPayload

	if err := c.BodyParser(&payload); err != nil {
		log.Printf("❌ Failed to parse payload: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid payload",
		})
	}

	log.Printf("✅ Payload parsed successfully")

	if err := h.service.ProcessWhatsAppWebhook(&payload); err != nil {
		log.Printf("❌ Failed to process webhook: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	log.Printf("✅ Webhook processed successfully")
	return c.JSON(fiber.Map{"success": true})
}

// Other handlers...
func (h *WebhookHandler) VerifyTelegram(c *fiber.Ctx) error {
	return c.SendString("Telegram webhook verified")
}

func (h *WebhookHandler) HandleTelegram(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"success": true})
}

func (h *WebhookHandler) VerifyInstagram(c *fiber.Ctx) error {
	return c.SendString("Instagram webhook verified")
}

func (h *WebhookHandler) HandleInstagram(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"success": true})
}
