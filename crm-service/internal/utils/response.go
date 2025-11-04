package utils

import (
	"github.com/gofiber/fiber/v2"
)

// StandardResponse is the standard API response format
type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse returns a success response
func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(StandardResponse{
		Success: true,
		Data:    data,
	})
}

// SuccessMessageResponse returns a success response with message
func SuccessMessageResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse returns an error response
func ErrorResponse(c *fiber.Ctx, statusCode int, message string, err error) error {
	response := StandardResponse{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	return c.Status(statusCode).JSON(response)
}

// ValidationErrorResponse returns a validation error response
func ValidationErrorResponse(c *fiber.Ctx, errors interface{}) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"message": "Validation failed",
		"errors":  errors,
	})
}

// NotFoundResponse returns a not found response
func NotFoundResponse(c *fiber.Ctx, resource string) error {
	return c.Status(fiber.StatusNotFound).JSON(StandardResponse{
		Success: false,
		Message: resource + " not found",
	})
}

// UnauthorizedResponse returns an unauthorized response
func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(StandardResponse{
		Success: false,
		Message: message,
	})
}
