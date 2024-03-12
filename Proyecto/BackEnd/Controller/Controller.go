package controller

import "github.com/gofiber/fiber/v2"

type Controller struct{}

func (c *Controller) Running(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status": "Server is running!!!",
	})
}
