// api-gateway/main.go
package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

type Service struct {
	Name string
	URL  string
}

var services = map[string]Service{
	"contact": {
		Name: "Contact Service",
		URL:  getEnv("CONTACT_SERVICE_URL", "http://localhost:3001"),
	},
	"chat": {
		Name: "Chat Service",
		URL:  getEnv("CHAT_SERVICE_URL", "http://localhost:3002"),
	},
	"ai": {
		Name: "AI Service",
		URL:  getEnv("AI_SERVICE_URL", "http://localhost:3003"),
	},
	"webhook": {
		Name: "Webhook Service",
		URL:  getEnv("WEBHOOK_SERVICE_URL", "http://localhost:3004"),
	},
}

func main() {
	godotenv.Load()

	app := fiber.New(fiber.Config{
		AppName:      "Divine CRM API Gateway",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	})

	// Middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Rate limiter
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "Too many requests. Please try again later.",
			})
		},
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "ok",
			"gateway":  "running",
			"time":     time.Now(),
			"services": getServicesStatus(),
		})
	})

	// API Routes
	api := app.Group("/api/v1")

	// Contact Service Routes
	api.All("/contacts*", proxyToService("contact", "/api/v1/contacts"))
	api.All("/leads*", proxyToService("contact", "/api/v1/leads"))
	api.All("/products*", proxyToService("contact", "/api/v1/products"))
	api.All("/masterdata/contacts*", proxyToService("contact", "/api/v1/masterdata/contacts"))
	api.All("/masterdata/leads*", proxyToService("contact", "/api/v1/masterdata/leads"))
	api.All("/masterdata/products*", proxyToService("contact", "/api/v1/masterdata/products"))

	// Chat Service Routes
	api.All("/chats*", proxyToService("chat", "/api/v1/chats"))
	api.All("/masterdata/chat-labels*", proxyToService("chat", "/api/v1/masterdata/chat-labels"))
	api.All("/masterdata/quick-replies*", proxyToService("chat", "/api/v1/masterdata/quick-replies"))
	api.All("/masterdata/broadcast-templates*", proxyToService("chat", "/api/v1/masterdata/broadcast-templates"))

	// AI Service Routes
	api.All("/ai/*", proxyToService("ai", "/api/v1/ai"))
	api.All("/masterdata/ai-configurations*", proxyToService("ai", "/api/v1/masterdata/ai-configurations"))
	api.All("/masterdata/ai-agents*", proxyToService("ai", "/api/v1/masterdata/ai-agents"))
	api.All("/analytics/*", proxyToService("ai", "/api/v1/analytics"))

	// Webhook Service Routes
	api.All("/webhooks/*", proxyToService("webhook", "/api/v1/webhooks"))
	api.All("/masterdata/connected-platforms*", proxyToService("webhook", "/api/v1/masterdata/connected-platforms"))

	port := getEnv("PORT", "8080")
	log.Printf("🚀 API Gateway starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func proxyToService(serviceName, stripPrefix string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		service, exists := services[serviceName]
		if !exists {
			return c.Status(503).JSON(fiber.Map{
				"error": "Service not available",
			})
		}

		// Build target URL
		targetURL := service.URL + c.OriginalURL()

		// Create request
		req, err := http.NewRequest(c.Method(), targetURL, bytes.NewReader(c.Body()))
		if err != nil {
			log.Printf("Error creating request to %s: %v", service.Name, err)
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to create request",
			})
		}

		// Copy headers
		c.Request().Header.VisitAll(func(key, value []byte) {
			req.Header.Set(string(key), string(value))
		})

		// Forward request
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error forwarding request to %s: %v", service.Name, err)
			return c.Status(503).JSON(fiber.Map{
				"error":   "Service unavailable",
				"service": service.Name,
			})
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response from %s: %v", service.Name, err)
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to read response",
			})
		}

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Set(key, value)
			}
		}

		// Send response
		c.Status(resp.StatusCode)
		return c.Send(body)
	}
}

func getServicesStatus() map[string]string {
	status := make(map[string]string)
	client := &http.Client{Timeout: 2 * time.Second}

	for key, service := range services {
		resp, err := client.Get(service.URL + "/health")
		if err != nil || resp.StatusCode != 200 {
			status[key] = "down"
		} else {
			status[key] = "up"
		}
		if resp != nil {
			resp.Body.Close()
		}
	}

	return status
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
