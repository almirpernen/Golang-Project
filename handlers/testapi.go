package handlers

import "github.com/gofiber/fiber/v2"

func TestApi(c *fiber.Ctx) error {
	return c.JSON("Hello world")
}
