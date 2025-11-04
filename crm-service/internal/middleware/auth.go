package middleware

import (
	"divine-crm/internal/config"
	"divine-crm/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.UnauthorizedResponse(c, "No authorization header")
		}

		// Extract token from "Bearer <token>"
		tokenString := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// Verify token
		claims, err := utils.VerifyToken(tokenString, cfg.JWT.Secret)
		if err != nil {
			return utils.UnauthorizedResponse(c, "Invalid token")
		}

		// Store user info in context
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}
