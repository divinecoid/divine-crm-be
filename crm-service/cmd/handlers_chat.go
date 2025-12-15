package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func getChatMessages(c *fiber.Ctx) error {
	var messages []ChatMessage
	status := c.Query("status")

	query := db.Order("created_at desc")
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}
	query.Find(&messages)
	return c.JSON(fiber.Map{"success": true, "data": messages})
}

func getChatMessageByID(c *fiber.Ctx) error {
	var message ChatMessage
	if err := db.First(&message, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	return c.JSON(fiber.Map{"success": true, "data": message})
}

func createChatMessage(c *fiber.Ctx) error {
	message := new(ChatMessage)
	c.BodyParser(message)
	db.Create(message)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": message})
}

func updateChatMessage(c *fiber.Ctx) error {
	message := new(ChatMessage)
	if err := db.First(message, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(message)
	db.Save(message)
	return c.JSON(fiber.Map{"success": true, "data": message})
}

func assignChat(c *fiber.Ctx) error {
	var req struct {
		AssignedTo    string `json:"assigned_to"`
		AssignedAgent string `json:"assigned_agent"`
	}
	c.BodyParser(&req)
	var message ChatMessage
	if err := db.First(&message, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	message.AssignedTo = req.AssignedTo
	message.AssignedAgent = req.AssignedAgent
	message.Status = "Assigned"
	db.Save(&message)
	return c.JSON(fiber.Map{"success": true, "data": message})
}

func resolveChat(c *fiber.Ctx) error {
	var message ChatMessage
	if err := db.First(&message, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	message.Status = "Resolved"
	db.Save(&message)
	return c.JSON(fiber.Map{"success": true, "data": message})
}

func getChatStats(c *fiber.Ctx) error {
	var unassigned, pending, assigned, resolved, total int64
	db.Model(&ChatMessage{}).Where("status = ?", "Unassigned").Count(&unassigned)
	db.Model(&ChatMessage{}).Where("status = ?", "Pending").Count(&pending)
	db.Model(&ChatMessage{}).Where("status = ?", "Assigned").Count(&assigned)
	db.Model(&ChatMessage{}).Where("status = ?", "Resolved").Count(&resolved)
	db.Model(&ChatMessage{}).Count(&total)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"unassigned": unassigned,
			"pending":    pending,
			"assigned":   assigned,
			"resolved":   resolved,
			"total":      total,
		},
	})
}

func getChatConversations(c *fiber.Ctx) error {
	status := c.Query("status")

	type ConversationResult struct {
		ContactID   uint   `json:"contact_id"`
		ContactName string `json:"contact_name"`
		Channel     string `json:"channel"`
		LastMessage string `json:"last_message"`
		LastStatus  string `json:"last_status"`
		UnreadCount int    `json:"unread_count"`
		LastUpdated string `json:"last_updated"`
	}

	var results []ConversationResult

	query := `
		SELECT 
			contact_id,
			contact_name,
			channel,
			message as last_message,
			status as last_status,
			0 as unread_count,
			updated_at as last_updated
		FROM chat_messages
		WHERE id IN (
			SELECT MAX(id) FROM chat_messages GROUP BY contact_id
		)
	`

	if status != "" && status != "all" {
		query += " AND status = ?"
		db.Raw(query, status).Scan(&results)
	} else {
		db.Raw(query).Scan(&results)
	}

	return c.JSON(fiber.Map{"success": true, "data": results})
}

func getChatsByContact(c *fiber.Ctx) error {
	contactId := c.Params("contactId")
	var messages []ChatMessage
	db.Where("contact_id = ?", contactId).Order("created_at asc").Find(&messages)
	return c.JSON(fiber.Map{"success": true, "data": messages})
}

func takeoverChat(c *fiber.Ctx) error {
	var req struct {
		AgentName string `json:"agent_name"`
	}
	c.BodyParser(&req)

	var message ChatMessage
	if err := db.First(&message, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}

	message.Status = "Assigned"
	message.AssignedTo = "Human"
	message.AssignedAgent = req.AgentName
	db.Save(&message)

	var contact Contact
	if db.First(&contact, message.ContactID).Error == nil {
		contact.LastAgent = req.AgentName
		contact.LastAgentType = "Human"
		db.Save(&contact)
	}

	log.Printf("👤 Chat takeover: Agent %s took over chat #%d", req.AgentName, message.ID)
	return c.JSON(fiber.Map{"success": true, "data": message, "message": "Takeover successful"})
}

func backToAIChat(c *fiber.Ctx) error {
	var message ChatMessage
	if err := db.First(&message, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}

	message.Status = "Unassigned"
	message.AssignedTo = "AI Bot"
	message.AssignedAgent = "AI Assistant"
	db.Save(&message)

	var contact Contact
	if db.First(&contact, message.ContactID).Error == nil {
		contact.LastAgent = "AI"
		contact.LastAgentType = "Bot"
		db.Save(&contact)
	}

	log.Printf("🤖 Chat returned to AI: chat #%d", message.ID)
	return c.JSON(fiber.Map{"success": true, "data": message, "message": "Returned to AI"})
}

func setPendingChat(c *fiber.Ctx) error {
	var req struct {
		FollowupDays int `json:"followup_days"`
	}
	c.BodyParser(&req)

	if req.FollowupDays == 0 {
		req.FollowupDays = 1
	}

	var message ChatMessage
	if err := db.First(&message, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}

	message.Status = "Pending"
	db.Save(&message)

	log.Printf("⏰ Chat set to pending: chat #%d, follow-up in %d days", message.ID, req.FollowupDays)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    message,
		"message": fmt.Sprintf("Set to pending, follow-up in %d days", req.FollowupDays),
	})
}

func sendReplyChat(c *fiber.Ctx) error {
	var req struct {
		Message   string `json:"message"`
		Channel   string `json:"channel"`
		ChannelID string `json:"channel_id"`
	}
	c.BodyParser(&req)

	var originalMsg ChatMessage
	if err := db.First(&originalMsg, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Chat not found"})
	}

	agentName := "Agent"
	if claims, ok := c.Locals("user").(*JWTClaims); ok {
		agentName = claims.FullName
	}

	newMsg := ChatMessage{
		ContactID:     originalMsg.ContactID,
		ContactName:   originalMsg.ContactName,
		Message:       "",
		Response:      req.Message,
		Channel:       req.Channel,
		Status:        "Assigned",
		AssignedTo:    "Human",
		AssignedAgent: agentName,
		TokensUsed:    0,
	}
	db.Create(&newMsg)

	var sendErr error
	switch req.Channel {
	case "WhatsApp":
		sendErr = sendWhatsAppMessage(req.ChannelID, req.Message)
	case "Telegram":
		sendErr = sendTelegramMessage(req.ChannelID, req.Message)
	case "Instagram":
		// Truncate message for Instagram 1000 char limit
		instagramMsg := req.Message
		if len(instagramMsg) > 950 {
			instagramMsg = instagramMsg[:950] + "..."
		}
		sendErr = sendInstagramMessage(req.ChannelID, instagramMsg)
	default:
		log.Printf("📨 Reply to %s via %s: %s", req.ChannelID, req.Channel, req.Message)
	}

	if sendErr != nil {
		log.Printf("❌ Failed to send message: %v", sendErr)
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to send message"})
	}

	var contact Contact
	if db.First(&contact, originalMsg.ContactID).Error == nil {
		contact.LastAgent = agentName
		contact.LastAgentType = "Human"
		contact.LastContact = time.Now()
		db.Save(&contact)
	}

	log.Printf("📤 Reply sent by %s to %s via %s", agentName, req.ChannelID, req.Channel)
	return c.JSON(fiber.Map{"success": true, "data": newMsg, "message": "Reply sent"})
}

func addChatLabel(c *fiber.Ctx) error {
	chatID := c.Params("id")

	var req struct {
		LabelID uint `json:"label_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid request"})
	}

	var chatMsg ChatMessage
	if err := db.First(&chatMsg, chatID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Chat not found"})
	}

	var label ChatLabel
	if err := db.First(&label, req.LabelID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Label not found"})
	}

	var labels []uint
	if chatMsg.Labels != "" {
		json.Unmarshal([]byte(chatMsg.Labels), &labels)
	}

	for _, l := range labels {
		if l == req.LabelID {
			return c.JSON(fiber.Map{"success": true, "message": "Label already added", "data": chatMsg})
		}
	}

	labels = append(labels, req.LabelID)
	labelsJSON, _ := json.Marshal(labels)
	chatMsg.Labels = string(labelsJSON)
	db.Save(&chatMsg)

	log.Printf("🏷️ Label '%s' added to chat #%d", label.Label, chatMsg.ID)

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Label '%s' added", label.Label),
		"data":    chatMsg,
	})
}

func removeChatLabel(c *fiber.Ctx) error {
	chatID := c.Params("id")
	labelID := c.Params("labelId")

	var chatMsg ChatMessage
	if err := db.First(&chatMsg, chatID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Chat not found"})
	}

	var labels []uint
	if chatMsg.Labels != "" {
		json.Unmarshal([]byte(chatMsg.Labels), &labels)
	}

	var labelIDUint uint
	fmt.Sscanf(labelID, "%d", &labelIDUint)

	var newLabels []uint
	removed := false
	for _, l := range labels {
		if l != labelIDUint {
			newLabels = append(newLabels, l)
		} else {
			removed = true
		}
	}

	if !removed {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Label not found on this chat"})
	}

	if len(newLabels) > 0 {
		labelsJSON, _ := json.Marshal(newLabels)
		chatMsg.Labels = string(labelsJSON)
	} else {
		chatMsg.Labels = ""
	}
	db.Save(&chatMsg)

	log.Printf("🏷️ Label removed from chat #%d", chatMsg.ID)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Label removed",
		"data":    chatMsg,
	})
}
