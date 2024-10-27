package middleware

import (
	"github.com/gofiber/fiber/v2"
	"im/internal/response"
)

func ErrorMiddlewareResponse(c *fiber.Ctx) error {
	if err := c.Next(); err != nil {
		return c.Status(fiber.StatusOK).
			JSON(response.Fail(500, err.Error()))
	}
	return nil
}
