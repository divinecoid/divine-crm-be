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
	"strconv"
)

type TelegramService struct {
	repo   *repository.PlatformRepository
	logger *utils.Logger
	client *http.Client
}

func NewTelegramService(repo *repository.PlatformRepository, logger *utils.Logger) *TelegramService {
	return &TelegramService{
		repo:   repo,
		logger: logger,
		client: &http.Client{},
	}
}

// SendMessage sends a message via Telegram Bot API
// Accepts both int64 and string chatID
func (s *TelegramService) SendMessage(chatID interface{}, message string) error {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if botToken == "" {
		return fmt.Errorf("Telegram bot token not configured")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	// Handle both int64 and string chatID
	var finalChatID int64
	switch v := chatID.(type) {
	case int64:
		finalChatID = v
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid chat ID: %s", v)
		}
		finalChatID = parsed
	default:
		return fmt.Errorf("unsupported chat ID type: %T", chatID)
	}

	s.logger.Info("📤 Sending Telegram message",
		"chat_id", finalChatID,
		"message_length", len(message),
	)

	reqBody := map[string]interface{}{
		"chat_id": finalChatID,
		"text":    message,
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error("Failed to send Telegram message", "error", err)
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("Telegram API error",
			"status", resp.StatusCode,
			"response", string(body),
		)
		return fmt.Errorf("Telegram API error: %s", string(body))
	}

	s.logger.Info("✅ Telegram message sent successfully")
	return nil
}
