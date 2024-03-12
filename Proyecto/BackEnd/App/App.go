package main

import (
	controller "P1/Controller"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	ctrl := &controller.Controller{}
	app.Get("/", ctrl.Running)
	app.Listen(":8000")
}
