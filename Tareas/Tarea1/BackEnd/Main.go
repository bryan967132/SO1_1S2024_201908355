package main

import (
	controller "T1/Controller"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	endpoints := controller.Controller{}

	app.Get("/", endpoints.Running)
	app.Get("/data", endpoints.Data)
	app.Listen(":8000")
}