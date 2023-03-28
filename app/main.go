package main //

import (
	"app/database"
	"app/handlers"
	"app/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello,   a World!")
	})
	handlers.UseRoute(app)
	//create mongodb instance
	utilities.InitDocker()
	utilities.InitProcfs()
	utilities.InitRedis()
	database.Connect()
	app.Listen(":5000")
}
