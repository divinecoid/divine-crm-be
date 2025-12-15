package main

import (
	"github.com/gofiber/fiber/v2"
)

func authMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "No authorization header"})
	}

	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	claims, err := verifyToken(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Invalid token"})
	}

	c.Locals("userID", claims.UserID)
	c.Locals("username", claims.Username)
	c.Locals("email", claims.Email)
	c.Locals("role", claims.Role)
	c.Locals("user", claims)

	return c.Next()
}

func requireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role").(string)
		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}
		return c.Status(403).JSON(fiber.Map{"success": false, "message": "Access denied"})
	}
}
