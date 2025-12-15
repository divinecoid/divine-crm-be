package main

import (
	"github.com/gofiber/fiber/v2"
)

func getContacts(c *fiber.Ctx) error {
	var contacts []Contact
	db.Order("created_at desc").Find(&contacts)
	return c.JSON(fiber.Map{"success": true, "data": contacts})
}

func getContactByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var contact Contact
	if err := db.First(&contact, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	return c.JSON(fiber.Map{"success": true, "data": contact})
}

func createContact(c *fiber.Ctx) error {
	contact := new(Contact)
	if err := c.BodyParser(contact); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if contact.Code == "" {
		contact.Code = generateCode("C", &Contact{})
	}
	db.Create(contact)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": contact})
}

func updateContact(c *fiber.Ctx) error {
	id := c.Params("id")
	contact := new(Contact)
	if err := db.First(contact, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(contact)
	db.Save(contact)
	return c.JSON(fiber.Map{"success": true, "data": contact})
}

func deleteContact(c *fiber.Ctx) error {
	db.Delete(&Contact{}, c.Params("id"))
	return c.JSON(fiber.Map{"success": true, "message": "Deleted"})
}

func getLeads(c *fiber.Ctx) error {
	var leads []Lead
	db.Order("created_at desc").Find(&leads)
	return c.JSON(fiber.Map{"success": true, "data": leads})
}

func createLead(c *fiber.Ctx) error {
	lead := new(Lead)
	c.BodyParser(lead)
	if lead.Code == "" {
		lead.Code = generateCode("L", &Lead{})
	}
	db.Create(lead)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": lead})
}

func updateLead(c *fiber.Ctx) error {
	lead := new(Lead)
	if err := db.First(lead, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(lead)
	db.Save(lead)
	return c.JSON(fiber.Map{"success": true, "data": lead})
}

func deleteLead(c *fiber.Ctx) error {
	db.Delete(&Lead{}, c.Params("id"))
	return c.JSON(fiber.Map{"success": true, "message": "Deleted"})
}

func getProducts(c *fiber.Ctx) error {
	var products []Product
	db.Order("created_at desc").Find(&products)
	return c.JSON(fiber.Map{"success": true, "data": products})
}

func createProduct(c *fiber.Ctx) error {
	product := new(Product)
	c.BodyParser(product)
	if product.Code == "" {
		product.Code = generateCode("P", &Product{})
	}
	db.Create(product)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": product})
}

func updateProduct(c *fiber.Ctx) error {
	product := new(Product)
	if err := db.First(product, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(product)
	db.Save(product)
	return c.JSON(fiber.Map{"success": true, "data": product})
}

func deleteProduct(c *fiber.Ctx) error {
	db.Delete(&Product{}, c.Params("id"))
	return c.JSON(fiber.Map{"success": true, "message": "Deleted"})
}

func getChatLabels(c *fiber.Ctx) error {
	var labels []ChatLabel
	db.Find(&labels)
	return c.JSON(fiber.Map{"success": true, "data": labels})
}

func createChatLabel(c *fiber.Ctx) error {
	label := new(ChatLabel)
	c.BodyParser(label)
	db.Create(label)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": label})
}

func updateChatLabel(c *fiber.Ctx) error {
	label := new(ChatLabel)
	if err := db.First(label, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(label)
	db.Save(label)
	return c.JSON(fiber.Map{"success": true, "data": label})
}

func deleteChatLabel(c *fiber.Ctx) error {
	db.Delete(&ChatLabel{}, c.Params("id"))
	return c.JSON(fiber.Map{"success": true, "message": "Deleted"})
}

func getQuickReplies(c *fiber.Ctx) error {
	var replies []QuickReply
	db.Find(&replies)
	return c.JSON(fiber.Map{"success": true, "data": replies})
}

func createQuickReply(c *fiber.Ctx) error {
	reply := new(QuickReply)
	c.BodyParser(reply)
	db.Create(reply)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": reply})
}

func updateQuickReply(c *fiber.Ctx) error {
	reply := new(QuickReply)
	if err := db.First(reply, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(reply)
	db.Save(reply)
	return c.JSON(fiber.Map{"success": true, "data": reply})
}

func deleteQuickReply(c *fiber.Ctx) error {
	db.Delete(&QuickReply{}, c.Params("id"))
	return c.JSON(fiber.Map{"success": true, "message": "Deleted"})
}

func getBroadcastTemplates(c *fiber.Ctx) error {
	var templates []BroadcastTemplate
	db.Find(&templates)
	return c.JSON(fiber.Map{"success": true, "data": templates})
}

func createBroadcastTemplate(c *fiber.Ctx) error {
	template := new(BroadcastTemplate)
	c.BodyParser(template)
	db.Create(template)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": template})
}

func updateBroadcastTemplate(c *fiber.Ctx) error {
	template := new(BroadcastTemplate)
	if err := db.First(template, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(template)
	db.Save(template)
	return c.JSON(fiber.Map{"success": true, "data": template})
}

func deleteBroadcastTemplate(c *fiber.Ctx) error {
	db.Delete(&BroadcastTemplate{}, c.Params("id"))
	return c.JSON(fiber.Map{"success": true, "message": "Deleted"})
}

func getAIConfigurations(c *fiber.Ctx) error {
	var configs []AIConfiguration
	db.Find(&configs)
	return c.JSON(fiber.Map{"success": true, "data": configs})
}

func createAIConfiguration(c *fiber.Ctx) error {
	config := new(AIConfiguration)
	c.BodyParser(config)
	db.Create(config)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": config})
}

func updateAIConfiguration(c *fiber.Ctx) error {
	config := new(AIConfiguration)
	if err := db.First(config, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(config)
	db.Save(config)
	return c.JSON(fiber.Map{"success": true, "data": config})
}

func deleteAIConfiguration(c *fiber.Ctx) error {
	db.Delete(&AIConfiguration{}, c.Params("id"))
	return c.JSON(fiber.Map{"success": true, "message": "Deleted"})
}

func getAIAgents(c *fiber.Ctx) error {
	var agents []AIAgent
	db.Find(&agents)
	return c.JSON(fiber.Map{"success": true, "data": agents})
}

func createAIAgent(c *fiber.Ctx) error {
	agent := new(AIAgent)
	c.BodyParser(agent)
	db.Create(agent)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": agent})
}

func updateAIAgent(c *fiber.Ctx) error {
	agent := new(AIAgent)
	if err := db.First(agent, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(agent)
	db.Save(agent)
	return c.JSON(fiber.Map{"success": true, "data": agent})
}

func deleteAIAgent(c *fiber.Ctx) error {
	db.Delete(&AIAgent{}, c.Params("id"))
	return c.JSON(fiber.Map{"success": true, "message": "Deleted"})
}

func getConnectedPlatforms(c *fiber.Ctx) error {
	var platforms []ConnectedPlatform
	db.Find(&platforms)
	return c.JSON(fiber.Map{"success": true, "data": platforms})
}

func createConnectedPlatform(c *fiber.Ctx) error {
	platform := new(ConnectedPlatform)
	c.BodyParser(platform)
	db.Create(platform)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": platform})
}

func updateConnectedPlatform(c *fiber.Ctx) error {
	platform := new(ConnectedPlatform)
	if err := db.First(platform, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Not found"})
	}
	c.BodyParser(platform)
	db.Save(platform)
	return c.JSON(fiber.Map{"success": true, "data": platform})
}

func deleteConnectedPlatform(c *fiber.Ctx) error {
	db.Delete(&ConnectedPlatform{}, c.Params("id"))
	return c.JSON(fiber.Map{"success": true, "message": "Deleted"})
}
