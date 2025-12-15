package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func getBroadcasts(c *fiber.Ctx) error {
	var broadcasts []Broadcast
	db.Order("created_at DESC").Find(&broadcasts)
	return c.JSON(fiber.Map{"success": true, "data": broadcasts})
}

func getBroadcastByID(c *fiber.Ctx) error {
	var broadcast Broadcast
	if err := db.First(&broadcast, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	return c.JSON(fiber.Map{"success": true, "data": broadcast})
}

func sendBroadcast(c *fiber.Ctx) error {
	var req struct {
		TemplateID   uint   `json:"template_id"`
		Message      string `json:"message"`
		Channel      string `json:"channel"`
		TargetFilter string `json:"target_filter"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid request"})
	}

	sentBy := "System"
	if claims, ok := c.Locals("user").(*JWTClaims); ok {
		sentBy = claims.FullName
	}

	templateName := "Custom"
	if req.TemplateID > 0 {
		var template BroadcastTemplate
		if db.First(&template, req.TemplateID).Error == nil {
			templateName = template.Name
			if req.Message == "" {
				req.Message = template.Content
			}
		}
	}

	if req.Message == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Message is required"})
	}

	var contacts []Contact
	query := db

	switch req.TargetFilter {
	case "hot":
		query = query.Where("temperature = ?", TempHot)
	case "warm":
		query = query.Where("temperature = ?", TempWarm)
	case "cold":
		query = query.Where("temperature = ?", TempCold)
	}

	if req.Channel != "All" && req.Channel != "" {
		query = query.Where("channel = ?", req.Channel)
	}

	query.Find(&contacts)

	if len(contacts) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "No contacts match the filter"})
	}

	broadcast := Broadcast{
		TemplateID:     req.TemplateID,
		TemplateName:   templateName,
		Message:        req.Message,
		Channel:        req.Channel,
		TargetFilter:   req.TargetFilter,
		TotalSent:      0,
		TotalDelivered: 0,
		TotalFailed:    0,
		Status:         "sending",
		SentBy:         sentBy,
	}
	db.Create(&broadcast)

	go func() {
		sent := 0
		delivered := 0
		failed := 0

		for _, contact := range contacts {
			var err error
			switch contact.Channel {
			case "WhatsApp":
				err = sendWhatsAppMessage(contact.ChannelID, req.Message)
			case "Telegram":
				err = sendTelegramMessage(contact.ChannelID, req.Message)
			case "Instagram":
				// Truncate message for Instagram 1000 char limit
				instagramMsg := req.Message
				if len(instagramMsg) > 950 {
					instagramMsg = instagramMsg[:950] + "..."
				}
				err = sendInstagramMessage(contact.ChannelID, instagramMsg)
			default:
				msgPreview := req.Message
				if len(msgPreview) > 50 {
					msgPreview = msgPreview[:50]
				}
				log.Printf("📨 Broadcast to %s via %s: %s", contact.Name, contact.Channel, msgPreview)
			}

			sent++
			if err != nil {
				failed++
				log.Printf("❌ Broadcast failed to %s: %v", contact.Name, err)
			} else {
				delivered++
				log.Printf("✅ Broadcast sent to %s via %s", contact.Name, contact.Channel)
			}

			// Rate limiting - Instagram has stricter limits
			if contact.Channel == "Instagram" {
				time.Sleep(500 * time.Millisecond) // 500ms for Instagram
			} else {
				time.Sleep(100 * time.Millisecond) // 100ms for others
			}
		}

		broadcast.TotalSent = sent
		broadcast.TotalDelivered = delivered
		broadcast.TotalFailed = failed
		broadcast.Status = "completed"
		db.Save(&broadcast)

		log.Printf("📢 Broadcast completed: %d sent, %d delivered, %d failed", sent, delivered, failed)
	}()

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Broadcasting to %d contacts", len(contacts)),
		"data":    broadcast,
	})
}
