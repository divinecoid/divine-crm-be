package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func processWithAI(c *fiber.Ctx) error {
	var req struct {
		Message string `json:"message"`
		AgentID uint   `json:"agent_id"`
	}
	c.BodyParser(&req)

	var agent AIAgent
	if err := db.Where("active = ?", true).First(&agent).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "No AI Agent"})
	}

	var config AIConfiguration
	if err := db.Where("ai_engine = ? AND active = ?", agent.AIEngine, true).First(&config).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "No AI Config"})
	}

	response, tokens, err := callOpenAI(&config, agent.BasicPrompt, req.Message)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"response":    response,
			"tokens_used": tokens,
			"agent":       agent.Name,
		},
	})
}

func suggestReply(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"success": true, "suggestions": []string{}})
}

func getTokenBalances(c *fiber.Ctx) error {
	var balances []TokenBalance
	db.Find(&balances)
	return c.JSON(fiber.Map{"success": true, "data": balances})
}

func debugAIPrompt(c *fiber.Ctx) error {
	var agent AIAgent
	if db.Where("active = ?", true).First(&agent).Error != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "No active AI agent"})
	}

	var products []Product
	db.Find(&products)

	var productsWithStock []Product
	db.Where("stock > 0").Find(&productsWithStock)

	fullPrompt := buildAIPromptWithProducts(agent.BasicPrompt)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"agent_name":          agent.Name,
			"agent_basic_prompt":  agent.BasicPrompt,
			"total_products":      len(products),
			"products_with_stock": len(productsWithStock),
			"products":            products,
			"full_ai_prompt":      fullPrompt,
			"prompt_length":       len(fullPrompt),
		},
	})
}

func callOpenAI(config *AIConfiguration, systemPrompt, userMessage string) (string, int, error) {
	aiReq := AIRequest{
		Model:       config.Model,
		Messages:    []AIMessage{{Role: "system", Content: systemPrompt}, {Role: "user", Content: userMessage}},
		Temperature: 0.7,
		MaxTokens:   500,
	}

	reqBody, _ := json.Marshal(aiReq)
	req, _ := http.NewRequest("POST", config.Endpoint, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", 0, fmt.Errorf("AI error: %s", string(body))
	}

	var aiResp OpenAIResponse
	json.Unmarshal(body, &aiResp)
	if len(aiResp.Choices) == 0 {
		return "", 0, fmt.Errorf("no response")
	}

	return aiResp.Choices[0].Message.Content, aiResp.Usage.TotalTokens, nil
}

func buildAIPromptWithProducts(basePrompt string) string {
	var products []Product
	db.Where("stock > 0").Find(&products)

	log.Printf("📦 Building AI prompt with %d products from database", len(products))

	if len(products) == 0 {
		log.Printf("⚠️ No products found in database!")
		return basePrompt + "\n\nIMPORTANT: No products are currently available in the catalog. If customer asks about products, inform them that no products are available at the moment."
	}

	productCatalog := "\n\n" + strings.Repeat("=", 50) + "\n"
	productCatalog += "⚠️ CRITICAL: OFFICIAL PRODUCT CATALOG - USE ONLY THESE PRODUCTS ⚠️\n"
	productCatalog += strings.Repeat("=", 50) + "\n\n"
	productCatalog += "YOU MUST ONLY MENTION PRODUCTS FROM THIS LIST. DO NOT INVENT OR HALLUCINATE ANY OTHER PRODUCTS.\n"
	productCatalog += "If a product is not in this list, it does not exist.\n\n"

	for _, p := range products {
		productCatalog += fmt.Sprintf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		productCatalog += fmt.Sprintf("📦 PRODUCT: %s\n", p.Name)
		productCatalog += fmt.Sprintf("   🔖 Code: %s\n", p.Code)
		productCatalog += fmt.Sprintf("   💰 Price: Rp %s\n", formatPrice(p.Price))
		productCatalog += fmt.Sprintf("   📊 Stock: %d units available\n", p.Stock)
		if p.Description != "" {
			productCatalog += fmt.Sprintf("   📝 Description: %s\n", p.Description)
		}
		log.Printf("   - Added product: %s (Rp %s, Stock: %d)", p.Name, formatPrice(p.Price), p.Stock)
	}

	productCatalog += fmt.Sprintf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
	productCatalog += strings.Repeat("=", 50) + "\n"
	productCatalog += "END OF OFFICIAL PRODUCT CATALOG\n"
	productCatalog += strings.Repeat("=", 50) + "\n\n"

	productCatalog += "🚨 STRICT RULES FOR PRODUCT INQUIRIES:\n"
	productCatalog += "1. ONLY mention products listed above - NO EXCEPTIONS\n"
	productCatalog += "2. Use the EXACT prices shown above\n"
	productCatalog += "3. Use the EXACT stock numbers shown above\n"
	productCatalog += "4. If customer asks about a product NOT in this list, say 'Maaf, produk tersebut tidak tersedia'\n"
	productCatalog += "5. NEVER invent products like 'Divine CRM Basic/Pro/Enterprise' - those do NOT exist\n"
	productCatalog += "6. When listing products, ONLY list the products above\n\n"

	return basePrompt + productCatalog
}

func analyzeTemperatureWithAI(contactID uint) string {
	var messages []ChatMessage
	db.Where("contact_id = ?", contactID).Order("created_at desc").Limit(3).Find(&messages)

	if len(messages) == 0 {
		return TempCold
	}

	var conversationText string
	for i := len(messages) - 1; i >= 0; i-- {
		conversationText += fmt.Sprintf("Customer: %s\n", messages[i].Message)
		if messages[i].Response != "" {
			conversationText += fmt.Sprintf("Response: %s\n", messages[i].Response)
		}
	}

	var config AIConfiguration
	if db.Where("active = ?", true).First(&config).Error != nil {
		log.Println("⚠️ No AI config for temperature analysis, using fallback")
		return analyzeTemperatureFallback(messages[0].Message)
	}

	systemPrompt := `You are a sales lead analyzer. Analyze the customer conversation and determine their buying intent temperature.

Respond with ONLY one word: Hot, Warm, or Cold

Definitions:
- Hot: Customer shows strong buying intent (wants to buy, order, checkout, asks for payment details, says "deal", "fix order", etc.)
- Warm: Customer shows interest but not ready to buy (asks about price, discount, product details, stock, shipping, compares options)
- Cold: Customer is just browsing or greeting (says hello, asks general questions, not interested, cancels, says "later", "think about it")

Analyze the entire conversation context, focusing on the LATEST messages to determine current intent.`

	userMessage := fmt.Sprintf("Analyze this conversation and determine the customer's buying temperature:\n\n%s\n\nRespond with only: Hot, Warm, or Cold", conversationText)

	response, _, err := callOpenAI(&config, systemPrompt, userMessage)
	if err != nil {
		log.Printf("⚠️ AI temperature analysis failed: %v, using fallback", err)
		return analyzeTemperatureFallback(messages[0].Message)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	switch {
	case strings.Contains(response, "hot"):
		return TempHot
	case strings.Contains(response, "warm"):
		return TempWarm
	default:
		return TempCold
	}
}

func analyzeTemperatureFallback(message string) string {
	lowerMsg := strings.ToLower(message)

	hotKeywords := []string{"beli", "order", "pesan", "checkout", "bayar", "transfer", "deal", "fix", "jadi", "kirim", "alamat", "buy", "purchase"}
	for _, kw := range hotKeywords {
		if strings.Contains(lowerMsg, kw) {
			return TempHot
		}
	}

	warmKeywords := []string{"harga", "berapa", "price", "diskon", "promo", "stock", "stok", "ready", "ongkir", "cod", "tertarik", "interested"}
	for _, kw := range warmKeywords {
		if strings.Contains(lowerMsg, kw) {
			return TempWarm
		}
	}

	return TempCold
}

func updateContactTemperatureWithAI(contact *Contact) {
	newTemp := analyzeTemperatureWithAI(contact.ID)
	oldTemp := contact.Temperature

	shouldUpdate := false

	switch oldTemp {
	case TempCold, "":
		if newTemp == TempWarm || newTemp == TempHot {
			shouldUpdate = true
		}
	case TempWarm:
		if newTemp == TempHot {
			shouldUpdate = true
		}
	case TempHot:
		shouldUpdate = false
	default:
		shouldUpdate = true
	}

	if shouldUpdate && newTemp != oldTemp {
		log.Printf("🔥 AI Temperature Upgrade: %s -> %s (contact: %s)", oldTemp, newTemp, contact.Name)
		contact.Temperature = newTemp
		db.Save(contact)
	} else {
		log.Printf("🌡️ AI Temperature: %s stays at %s (contact: %s)", contact.Name, contact.Temperature, contact.Name)
	}
}

func getTemperatureEmoji(temp string) string {
	switch temp {
	case TempHot:
		return "🔥"
	case TempWarm:
		return "☀️"
	case TempCold:
		return "❄️"
	default:
		return "❓"
	}
}
