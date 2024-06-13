package main

import (
	"flag"

	"github.com/devphaseX/hotel-reservation-api/types"
	"github.com/gofiber/fiber/v2"
)

func main() {
	listenAddress := flag.String("listenAddress", ":5000", "The listen address of the api server")

	flag.Parse()
	app := fiber.New()

	apiv1 := app.Group("/api/v1")

	app.Get("/foo", handleFoo)
	apiv1.Get("/user", handleUser)
	app.Listen(*listenAddress)
}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "Hello, it worked"})
}

func handleUser(c *fiber.Ctx) error {

	user := types.User{
		FirstName: "Ayomide",
		LastName:  "Lawal",
	}

	return c.JSON(user)
}
