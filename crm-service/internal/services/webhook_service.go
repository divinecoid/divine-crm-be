package services

import (
	"divine-crm/internal/utils"
	"fmt"
	"log"
)

type WebhookService struct {
	chatService      *ChatService
	whatsappService  *WhatsAppService
	instagramService *InstagramService
	telegramService  *TelegramService
	logger           *utils.Logger
}

func NewWebhookService(
	chatService *ChatService,
	whatsappService *WhatsAppService,
	instagramService *InstagramService,
	telegramService *TelegramService,
	logger *utils.Logger,
) *WebhookService {
	return &WebhookService{
		chatService:      chatService,
		whatsappService:  whatsappService,
		instagramService: instagramService,
		telegramService:  telegramService,
		logger:           logger,
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

// Instagram webhook payload
type InstagramWebhookPayload struct {
	Object string `json:"object"`
	Entry  []struct {
		ID        string `json:"id"`
		Time      int64  `json:"time"`
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Recipient struct {
				ID string `json:"id"`
			} `json:"recipient"`
			Timestamp int64 `json:"timestamp"`
			Message   struct {
				Mid  string `json:"mid"`
				Text string `json:"text"`
			} `json:"message"`
		} `json:"messaging"`
	} `json:"entry"`
}

func (s *WebhookService) ProcessInstagramWebhook(payload *InstagramWebhookPayload) error {
	log.Println("🔍 Processing Instagram webhook")

	if payload == nil || len(payload.Entry) == 0 {
		log.Println("⚠️  Empty Instagram payload")
		return fmt.Errorf("empty payload")
	}

	for _, entry := range payload.Entry {
		for _, messaging := range entry.Messaging {
			// Skip if no message text
			if messaging.Message.Text == "" {
				log.Println("ℹ️  No text in Instagram message, skipping")
				continue
			}

			senderID := messaging.Sender.ID
			messageText := messaging.Message.Text

			log.Printf("💬 Instagram DM from: %s", senderID)
			log.Printf("📝 Message: %s", messageText)

			// Process message and get AI response
			_, aiResponse, err := s.chatService.ProcessIncomingMessage(
				"Instagram",
				senderID,
				"Instagram User",
				messageText,
			)

			if err != nil {
				log.Printf("❌ Failed to process Instagram message: %v", err)
				// Send error message
				s.instagramService.SendMessage(
					senderID,
					"Maaf, terjadi kesalahan. Silakan coba lagi. 🙏",
				)
				continue
			}

			log.Printf("✅ AI Response: %s", aiResponse)

			// Send AI response back to user
			log.Printf("📤 Sending Instagram reply to %s...", senderID)
			if err := s.instagramService.SendMessage(senderID, aiResponse); err != nil {
				log.Printf("❌ Failed to send Instagram reply: %v", err)
				continue
			}

			log.Println("✅ Instagram reply sent successfully")
		}
	}

	return nil
}

// Telegram webhook payload
type TelegramWebhookPayload struct {
	UpdateID int64 `json:"update_id"`
	Message  struct {
		MessageID int64 `json:"message_id"`
		From      struct {
			ID        int64  `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID   int64  `json:"id"`
			Type string `json:"type"`
		} `json:"chat"`
		Date int64  `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

// Process Telegram webhook
func (s *WebhookService) ProcessTelegramWebhook(payload *TelegramWebhookPayload) error {
	if payload == nil || payload.Message.Text == "" {
		return nil
	}

	chatID := fmt.Sprintf("%d", payload.Message.Chat.ID)
	senderName := payload.Message.From.FirstName
	if payload.Message.From.LastName != "" {
		senderName += " " + payload.Message.From.LastName
	}
	messageText := payload.Message.Text

	log.Printf("📨 Telegram message from %s: %s", senderName, messageText)

	// Process message
	_, aiResponse, err := s.chatService.ProcessIncomingMessage(
		"Telegram",
		chatID,
		senderName,
		messageText,
	)

	if err != nil {
		log.Printf("❌ Error processing Telegram message: %v", err)
		return err
	}

	// Send response
	if err := s.telegramService.SendMessage(chatID, aiResponse); err != nil {
		log.Printf("❌ Error sending Telegram message: %v", err)
		return err
	}

	return nil
}
