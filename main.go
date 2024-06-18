package main

import (
	"context"
	"flag"
	"log"

	"github.com/devphaseX/hotel-reservation-api/api"
	"github.com/devphaseX/hotel-reservation-api/api/middleware"
	"github.com/devphaseX/hotel-reservation-api/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	var err error
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(db.URI))
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	app := fiber.New(config)
	var (
		//store initialization
		userStore  = db.NewMongoUserStore(client)
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)

		store = &db.Store{
			User:  userStore,
			Hotel: hotelStore,
			Room:  roomStore,
		}

		apiNonVersion = app.Group("/api")
		authApi       = apiNonVersion.Group("/auth")
		apiv1         = app.Group("/api/v1")
		//handler initialization
		userHandler  = api.NewUserHandler(userStore)
		authHandler  = api.NewAuthHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
	)

	//public route
	authApi.Post("/sign-in", authHandler.SignIn)

	//protected routes

	userApi := apiv1.Group("/users", middleware.JWTAuth)

	userApi.Get("/", userHandler.HandlerGetUsers)
	userApi.Post("/", userHandler.HandleCreateUser)
	userApi.Get("/users/:id", userHandler.HandleGetUser)
	userApi.Put("/users/:id", userHandler.HandleUpdateUser)
	userApi.Delete("/users/:id", userHandler.HandleDeleteUser)

	hotelv1Api := apiv1.Get("/hotels", middleware.JWTAuth)
	hotelv1Api.Get("/hotels", hotelHandler.HandlerGets)
	hotelv1Api.Get("/hotels/:id/rooms", hotelHandler.HandleGetRooms)

	app.Listen(*listenAddress)
}
