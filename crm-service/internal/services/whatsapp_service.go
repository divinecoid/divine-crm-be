package services

import (
	"bytes"
	"divine-crm/internal/repository"
	"divine-crm/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type WhatsAppService struct {
	repo   *repository.PlatformRepository
	logger *utils.Logger
	client *http.Client
}

func NewWhatsAppService(repo *repository.PlatformRepository, logger *utils.Logger) *WhatsAppService {
	return &WhatsAppService{
		repo:   repo,
		logger: logger,
		client: &http.Client{},
	}
}

// SendMessage sends a message via WhatsApp Business API
func (s *WhatsAppService) SendMessage(to, message string) error {
	// Get WhatsApp credentials from env
	accessToken := os.Getenv("WHATSAPP_ACCESS_TOKEN")
	phoneNumberID := os.Getenv("WHATSAPP_PHONE_NUMBER_ID")
	apiVersion := os.Getenv("WHATSAPP_API_VERSION")

	if accessToken == "" || phoneNumberID == "" {
		return fmt.Errorf("WhatsApp credentials not configured")
	}

	if apiVersion == "" {
		apiVersion = "v18.0"
	}

	// Build API URL
	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/messages", apiVersion, phoneNumberID)

	s.logger.Info("📤 Sending WhatsApp message",
		"to", to,
		"message_length", len(message),
	)

	// Prepare request body
	reqBody := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "text",
		"text": map[string]string{
			"preview_url": "false",
			"body":        message,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		s.logger.Error("Failed to marshal request", "error", err)
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		s.logger.Error("Failed to create request", "error", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Execute request
	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error("Failed to send request", "error", err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response", "error", err)
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		s.logger.Error("WhatsApp API error",
			"status", resp.StatusCode,
			"response", string(body),
		)
		return fmt.Errorf("WhatsApp API error (status %d): %s", resp.StatusCode, string(body))
	}

	s.logger.Info("✅ WhatsApp message sent successfully",
		"status", resp.StatusCode,
		"response", string(body),
	)

	return nil
}

// SendTemplateMessage sends a template message
func (s *WhatsAppService) SendTemplateMessage(to, templateName string, params map[string]interface{}) error {
	accessToken := os.Getenv("WHATSAPP_ACCESS_TOKEN")
	phoneNumberID := os.Getenv("WHATSAPP_PHONE_NUMBER_ID")
	apiVersion := os.Getenv("WHATSAPP_API_VERSION")

	if accessToken == "" || phoneNumberID == "" {
		return fmt.Errorf("WhatsApp credentials not configured")
	}

	if apiVersion == "" {
		apiVersion = "v18.0"
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/messages", apiVersion, phoneNumberID)

	reqBody := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "template",
		"template": map[string]interface{}{
			"name":     templateName,
			"language": map[string]string{"code": "id"},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("WhatsApp API error: %s", string(body))
	}

	return nil
}
