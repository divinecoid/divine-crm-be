package services

import (
	"divine-crm/internal/utils"
	"fmt"
	"log"
)

type WebhookService struct {
	chatService     *ChatService
	whatsappService *WhatsAppService
	logger          *utils.Logger
}

func NewWebhookService(
	chatService *ChatService,
	whatsappService *WhatsAppService,
	logger *utils.Logger,
) *WebhookService {
	return &WebhookService{
		chatService:     chatService,
		whatsappService: whatsappService,
		logger:          logger,
	}
}

// WhatsApp webhook payload structures
type WhatsAppWebhookPayload struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberID      string `json:"phone_number_id"`
				} `json:"metadata"`
				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					} `json:"profile"`
					WaID string `json:"wa_id"`
				} `json:"contacts"`
				Messages []struct {
					From      string `json:"from"`
					ID        string `json:"id"`
					Timestamp string `json:"timestamp"`
					Text      struct {
						Body string `json:"body"`
					} `json:"text"`
					Type string `json:"type"`
				} `json:"messages"`
			} `json:"value"`
			Field string `json:"field"`
		} `json:"changes"`
	} `json:"entry"`
}

func (s *WebhookService) ProcessWhatsAppWebhook(payload *WhatsAppWebhookPayload) error {
	// ✅ Extensive nil checks
	log.Printf("🔍 WebhookService: Starting ProcessWhatsAppWebhook")

	if s == nil {
		log.Printf("❌ WebhookService is nil!")
		return fmt.Errorf("webhook service is nil")
	}

	if s.logger == nil {
		log.Printf("⚠️  Logger is nil, using standard log")
	} else {
		s.logger.Info("Processing WhatsApp webhook")
	}

	if s.chatService == nil {
		log.Printf("❌ ChatService is nil!")
		return fmt.Errorf("chat service is nil")
	}

	if s.whatsappService == nil {
		log.Printf("❌ WhatsAppService is nil!")
		return fmt.Errorf("whatsapp service is nil")
	}

	// Validate payload
	if payload == nil {
		log.Printf("❌ Payload is nil!")
		return fmt.Errorf("payload is nil")
	}

	log.Printf("📊 Payload entries: %d", len(payload.Entry))

	if len(payload.Entry) == 0 {
		log.Printf("ℹ️  No entries in webhook")
		return nil
	}

	if len(payload.Entry[0].Changes) == 0 {
		log.Printf("ℹ️  No changes in webhook")
		return nil
	}

	value := payload.Entry[0].Changes[0].Value
	log.Printf("📊 Messages count: %d", len(value.Messages))

	if len(value.Messages) == 0 {
		log.Printf("ℹ️  No messages in webhook (might be status update)")
		return nil
	}

	message := value.Messages[0]

	if message.Text.Body == "" {
		log.Printf("ℹ️  Message has no text body")
		return nil
	}

	senderPhone := message.From
	messageText := message.Text.Body
	senderName := ""

	if len(value.Contacts) > 0 {
		senderName = value.Contacts[0].Profile.Name
	}

	log.Printf("📨 Message details: from=%s, phone=%s, text=%s", senderName, senderPhone, messageText)

	if s.logger != nil {
		s.logger.Info("WhatsApp message received",
			"from", senderName,
			"phone", senderPhone,
			"message", messageText,
		)
	}

	// Process with chat service
	log.Printf("🤖 Calling chatService.ProcessIncomingMessage...")

	outMsg, aiResponse, err := s.chatService.ProcessIncomingMessage(
		"WhatsApp",
		senderPhone,
		senderName,
		messageText,
	)

	if err != nil {
		log.Printf("❌ ChatService error: %v", err)
		if s.logger != nil {
			s.logger.Error("Failed to process message", "error", err)
		}
		return err
	}

	log.Printf("✅ ChatService returned: msg_id=%v, response=%s", outMsg.ID, aiResponse)

	// Send response via WhatsApp
	log.Printf("📤 Sending response to WhatsApp...")

	if err := s.whatsappService.SendMessage(senderPhone, aiResponse); err != nil {
		log.Printf("❌ WhatsApp send error: %v", err)
		if s.logger != nil {
			s.logger.Error("Failed to send WhatsApp message", "error", err)
		}
		return err
	}

	log.Printf("✅ Message sent successfully")
	if s.logger != nil {
		s.logger.Info("Message processed and sent successfully")
	}

	return nil
}
