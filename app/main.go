package main

import (
	"app/database"
	"app/handlers"
	"app/utilities"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello,World!")
	})
	handlers.UseRoute(app)
	//create mongodb instance
	utilities.InitDocker()
	utilities.InitProcfs()
	database.Connect()
	app.Listen(":5000")
}
