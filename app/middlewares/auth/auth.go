package auth

import (
	"app/utilities"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func Auth(c *fiber.Ctx) error {
	// Log the request method and path
	sess, err := utilities.Store.Get(c)
	if err != nil {
		panic(err)
	}
	name := sess.Get("name")
	if name == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	fmt.Println(name)

	println(c.Method(), c.Path())

	// Call the next middleware in the chain
	return c.Next()
}
func AdminAuth(c *fiber.Ctx) error {
	// Log the request method and path
	sess, err := utilities.Store.Get(c)
	if err != nil {
		panic(err)
	}
	name := sess.Get("name")
	if name == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	role := sess.Get("role")
	if role == nil || role != "admin" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	println(c.Method(), c.Path())

	// Call the next middleware in the chain
	return c.Next()
}

func FacultyAuth(c *fiber.Ctx) error {
	// Log the request method and path
	sess, err := utilities.Store.Get(c)
	if err != nil {
		panic(err)
	}
	name := sess.Get("name")
	if name == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	role := sess.Get("role")
	if role == nil || role != "faculty" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	println(c.Method(), c.Path())

	// Call the next middleware in the chain
	return c.Next()
}
func StudentAuth(c *fiber.Ctx) error {
	// Log the request method and path
	sess, err := utilities.Store.Get(c)
	if err != nil {
		panic(err)
	}
	name := sess.Get("name")
	if name == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	role := sess.Get("role")
	if role == nil || role != "student" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	println(c.Method(), c.Path())

	// Call the next middleware in the chain
	return c.Next()
}
