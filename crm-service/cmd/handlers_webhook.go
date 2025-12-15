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

	var latestChat ChatMessage
	humanHandled := false
	if db.Where("contact_id = ?", contact.ID).Order("created_at DESC").First(&latestChat).Error == nil {
		if latestChat.AssignedTo == "Human" && time.Since(latestChat.CreatedAt) < 24*time.Hour {
			humanHandled = true
			log.Printf("👤 Contact %s is being handled by human agent %s - skipping AI", contact.Name, latestChat.AssignedAgent)
		}
	}

	if humanHandled {
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

	if isNewContact && agent.IntroMessage != "" {
		introMsg := strings.ReplaceAll(agent.IntroMessage, "{name}", senderName)
		introMsg = strings.ReplaceAll(introMsg, "{agent_name}", agent.Name)
		sendWhatsAppMessage(senderPhone, introMsg)
		log.Printf("👋 Sent intro message to new contact %s", senderName)

		introChatMsg := ChatMessage{
			ContactID: contact.ID, ContactName: contact.Name, Message: "(New Contact)",
			Response: introMsg, Status: "Unassigned", AssignedTo: "AI Bot", AssignedAgent: agent.Name,
			Channel: "WhatsApp", TokensUsed: 0,
		}
		db.Create(&introChatMsg)
	}

	systemPrompt := buildAIPromptWithProducts(agent.BasicPrompt)
	aiResponse, tokens, _ := callOpenAI(&config, systemPrompt, messageText)

	chatMsg := ChatMessage{
		ContactID: contact.ID, ContactName: contact.Name, Message: messageText,
		Response: aiResponse, Status: "Unassigned", AssignedTo: "AI Bot", AssignedAgent: agent.Name,
		Channel: "WhatsApp", TokensUsed: tokens,
	}
	db.Create(&chatMsg)

	updateContactTemperatureWithAI(&contact)
	log.Printf("🌡️ Contact %s temperature: %s %s", contact.Name, contact.Temperature, getTemperatureEmoji(contact.Temperature))

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

	var latestChat ChatMessage
	humanHandled := false
	if db.Where("contact_id = ?", contact.ID).Order("created_at DESC").First(&latestChat).Error == nil {
		if latestChat.AssignedTo == "Human" && time.Since(latestChat.CreatedAt) < 24*time.Hour {
			humanHandled = true
			log.Printf("👤 Contact %s is being handled by human agent %s - skipping AI", contact.Name, latestChat.AssignedAgent)
		}
	}

	if humanHandled {
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

	if isNewContact && agent.IntroMessage != "" {
		introMsg := strings.ReplaceAll(agent.IntroMessage, "{name}", senderName)
		introMsg = strings.ReplaceAll(introMsg, "{agent_name}", agent.Name)
		sendTelegramMessage(senderID, introMsg)
		log.Printf("👋 Sent intro message to new contact %s", senderName)

		introChatMsg := ChatMessage{
			ContactID: contact.ID, ContactName: contact.Name, Message: "(New Contact)",
			Response: introMsg, Status: "Unassigned", AssignedTo: "AI Bot", AssignedAgent: agent.Name,
			Channel: "Telegram", TokensUsed: 0,
		}
		db.Create(&introChatMsg)
	}

	systemPrompt := buildAIPromptWithProducts(agent.BasicPrompt)
	aiResponse, tokens, _ := callOpenAI(&config, systemPrompt, messageText)

	chatMsg := ChatMessage{
		ContactID: contact.ID, ContactName: contact.Name, Message: messageText,
		Response: aiResponse, Status: "Unassigned", AssignedTo: "AI Bot", AssignedAgent: agent.Name,
		Channel: "Telegram", TokensUsed: tokens,
	}
	db.Create(&chatMsg)

	updateContactTemperatureWithAI(&contact)
	log.Printf("🌡️ Contact %s temperature: %s %s", contact.Name, contact.Temperature, getTemperatureEmoji(contact.Temperature))

	contact.LastAgent = agent.Name
	contact.LastAgentType = "AI"
	db.Save(&contact)

	sendTelegramMessage(senderID, aiResponse)
	return c.JSON(fiber.Map{"status": "ok"})
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
