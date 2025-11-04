package services

import (
	"bytes"
	"divine-crm/internal/repository"
	"divine-crm/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type WhatsAppService struct {
	platformRepo *repository.PlatformRepository
	logger       *utils.Logger
}

func NewWhatsAppService(platformRepo *repository.PlatformRepository, logger *utils.Logger) *WhatsAppService {
	return &WhatsAppService{
		platformRepo: platformRepo,
		logger:       logger,
	}
}

// WhatsApp API request structure
type WhatsAppMessageRequest struct {
	MessagingProduct string `json:"messaging_product"`
	RecipientType    string `json:"recipient_type"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Text             struct {
		PreviewURL bool   `json:"preview_url"`
		Body       string `json:"body"`
	} `json:"text"`
}

// SendMessage sends a WhatsApp message
func (s *WhatsAppService) SendMessage(phone, message string) error {
	// ✅ Safe nil checks
	if s == nil {
		log.Printf("❌ WhatsAppService is nil")
		return fmt.Errorf("whatsapp service is nil")
	}

	if s.logger != nil {
		s.logger.Info("Attempting to send WhatsApp message", "phone", phone)
	} else {
		log.Printf("📤 Attempting to send WhatsApp message to %s", phone)
	}

	if s.platformRepo == nil {
		log.Printf("⚠️  Platform repository is nil, skipping WhatsApp send")
		return nil // Don't fail
	}

	// Get WhatsApp platform configuration
	platform, err := s.platformRepo.GetByPlatform("WhatsApp")

	// ✅ Handle all error cases safely
	if err != nil || platform == nil {
		if s.logger != nil {
			s.logger.Warn("WhatsApp platform not configured", "error", err)
		}
		log.Printf("⚠️  WhatsApp not configured in database")
		log.Printf("💬 Message would have been sent: %s", message)
		log.Printf("📝 To setup WhatsApp, add config to connected_platforms table")
		return nil // Don't fail - just skip sending
	}

	// ✅ Extra safety check
	if platform == nil {
		log.Printf("⚠️  Platform object is nil after query")
		return nil
	}

	if !platform.Active {
		if s.logger != nil {
			s.logger.Warn("WhatsApp platform is not active")
		}
		log.Printf("⚠️  WhatsApp platform exists but is not active")
		return nil
	}

	// Validate required fields
	if platform.PhoneNumberID == "" || platform.Token == "" {
		log.Printf("⚠️  WhatsApp config incomplete (missing phone_number_id or token)")
		return nil
	}

	// Prepare request
	req := WhatsAppMessageRequest{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               phone,
		Type:             "text",
	}
	req.Text.Body = message
	req.Text.PreviewURL = false

	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Build API URL
	apiURL := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/messages", platform.PhoneNumberID)

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+platform.Token)

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("WhatsApp API error [%d]: %s", resp.StatusCode, string(body))
	}

	if s.logger != nil {
		s.logger.Info("WhatsApp message sent successfully", "phone", phone)
	}
	log.Printf("✅ WhatsApp message sent successfully to %s", phone)

	return nil
}
