package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/devphaseX/hotel-reservation-api/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	uri         = "mongodb://localhost:27017"
	dbName      = "hotel-reservation"
	userColName = "users"
)

func main() {
	listenAddress := flag.String("listenAddress", ":5000", "The listen address of the api server")
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	userCol := client.Database(dbName).Collection(userColName)

	user := types.User{
		FirstName: "Ayomide",
		LastName:  "Lawal",
	}

	res, err := userCol.InsertOne(context.Background(), user)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)

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
