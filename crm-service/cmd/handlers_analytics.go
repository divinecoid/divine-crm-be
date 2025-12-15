package main

import (
	"github.com/gofiber/fiber/v2"
)

func getAnalyticsOverview(c *fiber.Ctx) error {
	var contacts, messages, products, resolved, pending int64
	db.Model(&Contact{}).Count(&contacts)
	db.Model(&ChatMessage{}).Count(&messages)
	db.Model(&Product{}).Count(&products)
	db.Model(&ChatMessage{}).Where("status = ?", "Resolved").Count(&resolved)
	db.Model(&ChatMessage{}).Where("status IN ?", []string{"Pending", "Unassigned"}).Count(&pending)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"total_contacts": contacts,
			"total_messages": messages,
			"total_products": products,
			"resolved_chats": resolved,
			"pending_chats":  pending,
		},
	})
}
