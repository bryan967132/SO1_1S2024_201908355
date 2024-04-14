package main

import (
	controller "P1/Controller"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	ctrl := controller.NewController()
	app.Get("/", ctrl.Running)
	app.Get("/cpuram", ctrl.Cpuram)
	app.Post("/inscpuram", ctrl.InsRAMCPU)
	app.Get("/history", ctrl.History)
	app.Get("/pids", ctrl.Pids)
	app.Get("/proc/:pid", ctrl.Proc)
	app.Get("/thereisproc", ctrl.ThereIsProc)
	app.Get("/start", ctrl.Start)
	app.Get("/stop/:pid", ctrl.Stop)
	app.Get("/resume/:pid", ctrl.Resume)
	app.Get("/kill/:pid", ctrl.Kill)
	app.Listen(":8000")
}
