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

type InstagramService struct {
	repo   *repository.PlatformRepository
	logger *utils.Logger
	client *http.Client
}

func NewInstagramService(repo *repository.PlatformRepository, logger *utils.Logger) *InstagramService {
	return &InstagramService{
		repo:   repo,
		logger: logger,
		client: &http.Client{},
	}
}

// SendMessage sends a message via Instagram Messaging API
func (s *InstagramService) SendMessage(recipientID, message string) error {
	accessToken := os.Getenv("INSTAGRAM_ACCESS_TOKEN")
	pageID := os.Getenv("INSTAGRAM_PAGE_ID")

	if accessToken == "" {
		return fmt.Errorf("INSTAGRAM_ACCESS_TOKEN not configured")
	}

	if pageID == "" {
		return fmt.Errorf("INSTAGRAM_PAGE_ID not configured")
	}

	// Instagram uses different API endpoint than WhatsApp
	url := fmt.Sprintf("https://graph.facebook.com/v18.0/me/messages?access_token=%s", accessToken)

	s.logger.Info("📤 Sending Instagram message",
		"recipient_id", recipientID,
		"message_length", len(message),
	)

	// Prepare request body
	reqBody := map[string]interface{}{
		"recipient": map[string]string{
			"id": recipientID,
		},
		"message": map[string]string{
			"text": message,
		},
		"messaging_type": "RESPONSE", // Important for Instagram
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error("Failed to send Instagram message", "error", err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		s.logger.Error("Instagram API error",
			"status", resp.StatusCode,
			"response", string(body),
		)
		return fmt.Errorf("Instagram API error (status %d): %s", resp.StatusCode, string(body))
	}

	s.logger.Info("✅ Instagram message sent successfully",
		"response", string(body),
	)

	return nil
}
