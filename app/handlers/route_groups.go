package handlers

import (
	dashboardhandler "app/handlers/dashboardHandler"
	"app/handlers/imageHandler"
	"app/handlers/userHandler"
	"app/middlewares/auth"
	"app/utilities"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/gofiber/fiber/v2"
)

var api *fiber.App

func UseRoute(app *fiber.App) {

	api = app
	user_handler := app.Group("/api/user")
	userHandler.Register(user_handler)
	image_handler := app.Group("/api/image")
	imageHandler.Register(image_handler)
	dashboard_handler := app.Group("/dashboard")
	dashboardhandler.Register(dashboard_handler)
	api.Get("/test", auth.Auth, func(c *fiber.Ctx) error { return c.SendString("hello") })
	api.Get("/testadmin", auth.AdminAuth, func(c *fiber.Ctx) error { return c.SendString("hello") })
	api.Get("/teststudent", auth.Auth, func(c *fiber.Ctx) error { return c.SendString("hello") })
	api.Get("/testdocker", func(c *fiber.Ctx) error {
		containers, err := utilities.Docker.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			panic(err)
		}

		for _, container := range containers {
			fmt.Printf("%s %s\n", container.ID[:10], container.Image)
		}
		json, err := json.Marshal(containers)
		return c.Send(json)
	})
}
