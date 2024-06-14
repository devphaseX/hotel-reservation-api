package main

import (
	"context"
	"flag"
	"log"

	"github.com/devphaseX/hotel-reservation-api/api"
	"github.com/devphaseX/hotel-reservation-api/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	uri         = "mongodb://localhost:27017"
	dbName      = "hotel-reservation"
	userColName = "users"
)

// Create a new fiber instance with custom config
var config = fiber.Config{
	// Override default error handler
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddress := flag.String("listenAddress", ":5000", "The listen address of the api server")
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	app := fiber.New(config)

	apiv1 := app.Group("/api/v1")

	//handler initialization
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client, dbName))

	app.Get("/foo", handleFoo)
	apiv1.Get("/users", userHandler.HandlerGetUsers)
	apiv1.Post("/users", userHandler.HandleCreateUser)
	apiv1.Get("/users/:id", userHandler.HandleGetUser)
	apiv1.Put("/users/:id", userHandler.HandleUpdateUser)
	apiv1.Delete("/users/:id", userHandler.HandleDeleteUser)
	app.Listen(*listenAddress)
}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "Hello, it worked"})
}
