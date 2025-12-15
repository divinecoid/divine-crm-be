package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ==================== INSTAGRAM WEBHOOK HANDLERS ====================

func verifyInstagramWebhook(c *fiber.Ctx) error {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	verifyToken := os.Getenv("INSTAGRAM_VERIFY_TOKEN")
	if verifyToken == "" {
		verifyToken = "divine_webhook_verify_123"
	}

	log.Printf("📸 Instagram Webhook Verification - Mode: %s, Token: %s", mode, token)

	if mode == "subscribe" && token == verifyToken {
		log.Println("✅ Instagram Webhook verified successfully")
		return c.SendString(challenge)
	}

	log.Println("❌ Instagram Webhook verification failed")
	return c.Status(403).SendString("Forbidden")
}

func handleInstagramWebhook(c *fiber.Ctx) error {
	body := c.Body()
	log.Printf("📸 Instagram Webhook received: %s", string(body))

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("❌ Failed to parse Instagram webhook: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	// Check object type
	object, _ := payload["object"].(string)
	if object != "instagram" {
		log.Printf("⚠️ Unknown object type: %s", object)
		return c.SendStatus(200)
	}

	// Get entry array
	entries, ok := payload["entry"].([]interface{})
	if !ok || len(entries) == 0 {
		return c.SendStatus(200)
	}

	for _, entry := range entries {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}

		// Handle messaging (DM)
		messaging, ok := entryMap["messaging"].([]interface{})
		if ok {
			for _, msg := range messaging {
				if msgMap, ok := msg.(map[string]interface{}); ok {
					go processInstagramMessage(msgMap)
				}
			}
		}

		// Handle changes (comments, mentions, etc)
		changes, ok := entryMap["changes"].([]interface{})
		if ok {
			for _, change := range changes {
				if changeMap, ok := change.(map[string]interface{}); ok {
					go processInstagramChange(changeMap)
				}
			}
		}
	}

	// Always return 200 OK to acknowledge receipt
	return c.SendStatus(200)
}

func processInstagramMessage(msg map[string]interface{}) {
	// Extract sender info
	sender, _ := msg["sender"].(map[string]interface{})
	senderID, _ := sender["id"].(string)

	// Extract recipient (your page)
	recipient, _ := msg["recipient"].(map[string]interface{})
	recipientID, _ := recipient["id"].(string)

	// Skip if message is from your own page (echo)
	if senderID == recipientID {
		log.Printf("📸 Skipping echo message from own page")
		return
	}

	// Extract message content
	message, ok := msg["message"].(map[string]interface{})
	if !ok {
		log.Printf("📸 No message content found")
		return
	}

	// Check if it's an echo (sent by page)
	if isEcho, ok := message["is_echo"].(bool); ok && isEcho {
		log.Printf("📸 Skipping echo message")
		return
	}

	messageText, _ := message["text"].(string)
	messageID, _ := message["mid"].(string)

	if messageText == "" {
		log.Printf("📸 Empty message text, might be attachment")
		// Check for attachments
		if attachments, ok := message["attachments"].([]interface{}); ok && len(attachments) > 0 {
			messageText = "[Attachment received]"
		} else {
			return
		}
	}

	log.Printf("📸 Instagram DM from %s: %s", senderID, messageText)

	// Get sender name from Instagram API (optional, can be empty)
	senderName := fmt.Sprintf("IG_%s", senderID[:8])

	// Find or create contact
	var contact Contact
	isNewContact := false
	if db.Where("channel = ? AND channel_id = ?", "Instagram", senderID).First(&contact).Error != nil {
		contact = Contact{
			Code:         generateCode("C", &Contact{}),
			Channel:      "Instagram",
			ChannelID:    senderID,
			Name:         senderName,
			Temperature:  TempCold,
			FirstContact: time.Now(),
			LastContact:  time.Now(),
		}
		db.Create(&contact)
		isNewContact = true
		log.Printf("📸 New Instagram contact: %s (ID: %s)", senderName, senderID)
	} else {
		contact.LastContact = time.Now()
		db.Save(&contact)
	}

	// Check if being handled by human agent
	var latestChat ChatMessage
	humanHandled := false
	if db.Where("contact_id = ?", contact.ID).Order("created_at DESC").First(&latestChat).Error == nil {
		if latestChat.AssignedTo == "Human" && time.Since(latestChat.CreatedAt) < 24*time.Hour {
			humanHandled = true
			log.Printf("👤 Instagram contact %s is being handled by human agent %s - skipping AI", contact.Name, latestChat.AssignedAgent)
		}
	}

	// If human is handling, just save the message
	if humanHandled {
		chatMsg := ChatMessage{
			ContactID:     contact.ID,
			ContactName:   contact.Name,
			Message:       messageText,
			Response:      "",
			Status:        "Assigned",
			AssignedTo:    "Human",
			AssignedAgent: latestChat.AssignedAgent,
			Channel:       "Instagram",
			TokensUsed:    0,
		}
		db.Create(&chatMsg)
		log.Printf("📥 Instagram message from %s saved for human agent %s", contact.Name, latestChat.AssignedAgent)
		return
	}

	// Get AI Agent for Instagram
	var agent AIAgent
	if db.Where("platform = ? AND active = ?", "Instagram", true).First(&agent).Error != nil {
		if db.Where("platform = ? AND active = ?", "instagram", true).First(&agent).Error != nil {
			if db.Where("platform = ? AND active = ?", "All", true).First(&agent).Error != nil {
				if db.Where("active = ?", true).First(&agent).Error != nil {
					log.Printf("❌ No active AI agent found for Instagram")
					return
				}
			}
		}
	}
	log.Printf("🤖 Instagram using AI Agent: %s (Engine: %s)", agent.Name, agent.AIEngine)

	// Get AI Configuration
	var config AIConfiguration
	if db.Where("ai_engine = ? AND active = ?", agent.AIEngine, true).First(&config).Error != nil {
		log.Printf("❌ No AI configuration found for engine: %s", agent.AIEngine)
		return
	}

	// Send intro message for new contact
	if isNewContact && agent.IntroMessage != "" {
		introMsg := strings.ReplaceAll(agent.IntroMessage, "{name}", senderName)
		introMsg = strings.ReplaceAll(introMsg, "{agent_name}", agent.Name)

		// Truncate for Instagram limit
		introMsg = truncateForInstagram(introMsg)

		err := sendInstagramMessage(senderID, introMsg)
		if err != nil {
			log.Printf("❌ Failed to send Instagram intro: %v", err)
		} else {
			log.Printf("👋 Sent intro message to new Instagram contact %s", senderName)
		}

		introChatMsg := ChatMessage{
			ContactID:     contact.ID,
			ContactName:   contact.Name,
			Message:       "(New Contact)",
			Response:      introMsg,
			Status:        "Unassigned",
			AssignedTo:    "AI Bot",
			AssignedAgent: agent.Name,
			Channel:       "Instagram",
			TokensUsed:    0,
		}
		db.Create(&introChatMsg)
	}

	// Build AI prompt with products
	systemPrompt := buildAIPromptWithProducts(agent.BasicPrompt)
	if systemPrompt == "" {
		systemPrompt = agent.BasicPrompt
	}

	// Add Instagram-specific instruction for shorter responses
	systemPrompt += "\n\n🔴 PENTING: Ini adalah Instagram DM. Batasi respons maksimal 800 karakter. Jawab singkat, padat, dan jelas. Jika perlu info lebih detail, tawarkan untuk melanjutkan percakapan."

	// Call AI
	aiResponse, tokens, err := callOpenAI(&config, systemPrompt, messageText)
	if err != nil {
		log.Printf("❌ AI call failed: %v", err)
		aiResponse = "Maaf, terjadi kendala teknis. Silakan coba lagi atau hubungi admin kami."
		tokens = 0
	}

	// Save chat message
	chatMsg := ChatMessage{
		ContactID:     contact.ID,
		ContactName:   contact.Name,
		Message:       messageText,
		Response:      aiResponse,
		Status:        "Unassigned",
		AssignedTo:    "AI Bot",
		AssignedAgent: agent.Name,
		Channel:       "Instagram",
		TokensUsed:    tokens,
	}
	db.Create(&chatMsg)

	// Update contact temperature
	updateContactTemperatureWithAI(&contact)
	log.Printf("🌡️ Instagram contact %s temperature: %s %s", contact.Name, contact.Temperature, getTemperatureEmoji(contact.Temperature))

	// Update contact's last agent
	contact.LastAgent = agent.Name
	contact.LastAgentType = "AI"
	db.Save(&contact)

	// Send reply via Instagram (with 1000 char limit)
	instagramResponse := truncateForInstagram(aiResponse)
	err = sendInstagramMessage(senderID, instagramResponse)
	if err != nil {
		log.Printf("❌ Failed to send Instagram reply: %v", err)
	} else {
		log.Printf("📤 Instagram reply sent to %s: %s", senderID, truncateString(instagramResponse, 50))
	}

	// Log message ID for reference
	if messageID != "" {
		log.Printf("📸 Processed Instagram message ID: %s", messageID)
	}
}

func processInstagramChange(change map[string]interface{}) {
	field, _ := change["field"].(string)
	value, _ := change["value"].(map[string]interface{})

	log.Printf("📸 Instagram Change - Field: %s", field)

	switch field {
	case "comments":
		// Handle new comments
		commentID, _ := value["id"].(string)
		text, _ := value["text"].(string)
		log.Printf("💬 New Instagram comment: %s - %s", commentID, truncateString(text, 50))

	case "mentions":
		// Handle mentions
		log.Printf("🏷️ New Instagram mention received")

	case "story_insights":
		log.Printf("📊 Instagram story insights received")

	default:
		log.Printf("📸 Unhandled Instagram change type: %s", field)
	}
}

// sendInstagramMessage sends a message via Instagram Messaging API
func sendInstagramMessage(recipientID, message string) error {
	// Try to get from connected_platforms first
	var platform ConnectedPlatform
	if db.Where("platform = ? AND active = ?", "Instagram", true).First(&platform).Error != nil {
		// Fall back to environment variable
		accessToken := os.Getenv("INSTAGRAM_ACCESS_TOKEN")
		pageID := os.Getenv("INSTAGRAM_PAGE_ID")

		if accessToken == "" {
			return fmt.Errorf("Instagram not configured - no access token")
		}

		return sendInstagramMessageDirect(pageID, accessToken, recipientID, message)
	}

	return sendInstagramMessageDirect(platform.PlatformID, platform.Token, recipientID, message)
}

func sendInstagramMessageDirect(igAccountID, accessToken, recipientID, message string) error {
	// Instagram Messaging API - uses graph.instagram.com
	url := fmt.Sprintf("https://graph.instagram.com/v18.0/%s/messages", igAccountID)

	payload := map[string]interface{}{
		"recipient": map[string]string{
			"id": recipientID,
		},
		"message": map[string]string{
			"text": message,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	log.Printf("📸 Sending Instagram message to %s via IG account %s", recipientID, igAccountID)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)

	if resp.StatusCode != 200 {
		log.Printf("❌ Instagram API error (status %d): %v", resp.StatusCode, respBody)
		return fmt.Errorf("Instagram API error: status %d - %v", resp.StatusCode, respBody)
	}

	log.Printf("✅ Instagram message sent successfully to %s", recipientID)
	return nil
}

// Helper function to truncate string for logging
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// truncateForInstagram limits message to Instagram's 1000 character limit
func truncateForInstagram(message string) string {
	maxLen := 950 // Leave some buffer for safety
	if len(message) <= maxLen {
		return message
	}

	// Try to cut at last sentence or word
	truncated := message[:maxLen]

	// Find last sentence ending
	lastPeriod := strings.LastIndex(truncated, ". ")
	lastQuestion := strings.LastIndex(truncated, "? ")
	lastExclaim := strings.LastIndex(truncated, "! ")

	// Get the latest sentence ending
	cutPoint := lastPeriod
	if lastQuestion > cutPoint {
		cutPoint = lastQuestion
	}
	if lastExclaim > cutPoint {
		cutPoint = lastExclaim
	}

	// If found a sentence ending, cut there
	if cutPoint > maxLen/2 {
		return truncated[:cutPoint+1] + "\n\n(Pesan terpotong, ketik 'lanjut' untuk info lebih lanjut)"
	}

	// Otherwise find last space
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > maxLen/2 {
		return truncated[:lastSpace] + "...\n\n(Ketik 'lanjut' untuk info lebih lanjut)"
	}

	return truncated + "..."
}
