package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {}

func (c Controller) Running(ctx *fiber.Ctx) error {
	return ctx.SendString("Server is running!!!")
}

func (c Controller) Data(ctx *fiber.Ctx) error {
	fmt.Println(`{
	"carnet": "201908355",
	"nombre": "Danny Hugo Bryan Tejaxún Pichiyá"
}`)
	return ctx.JSON(fiber.Map {
		"carnet": "201908355",
		"nombre": "Danny Hugo Bryan Tejaxún Pichiyá",
	})
}