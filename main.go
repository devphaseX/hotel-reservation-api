package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/devphaseX/hotel-reservation-api/api"
	"github.com/devphaseX/hotel-reservation-api/config"

	"github.com/devphaseX/hotel-reservation-api/api/middleware"
	"github.com/devphaseX/hotel-reservation-api/db"
	"github.com/devphaseX/hotel-reservation-api/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create a new fiber instance with custom fiberConfig
var fiberConfig = fiber.Config{
	// Override default error handler
	ErrorHandler: utils.ErrorHandler,
}

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.EnvConfig.MongoDBUrl))
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	app := fiber.New(fiberConfig)
	var (
		//store initialization
		userStore    = db.NewMongoUserStore(client)
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewMongoBookingStore(client)

		store = &db.Store{
			User:    userStore,
			Hotel:   hotelStore,
			Room:    roomStore,
			Booking: bookingStore,
		}

		apiNonVersion = app.Group("/api")
		authApi       = apiNonVersion.Group("/auth")
		apiv1         = app.Group("/api/v1", middleware.JWTAuth(userStore))
		adminv1Api    = apiv1.Group("/admin", middleware.AdminAuth)
		//handler initialization
		userHandler    = api.NewUserHandler(userStore)
		authHandler    = api.NewAuthHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
	)

	//public route
	authApi.Post("/sign-in", authHandler.SignIn)

	//protected routes

	userApi := apiv1.Group("/users", middleware.JWTAuth(userStore))

	userApi.Get("/", userHandler.HandlerGetUsers)
	userApi.Post("/", userHandler.HandleCreateUser)
	userApi.Get("/:id", userHandler.HandleGetUser)
	userApi.Put("/:id", userHandler.HandleUpdateUser)
	userApi.Delete("/users/:id", userHandler.HandleDeleteUser)

	hotelv1Api := apiv1.Group("/hotels", middleware.JWTAuth(userStore))
	hotelv1Api.Get("/", hotelHandler.HandlerGets)
	hotelv1Api.Get("/:id", hotelHandler.HandleGet)

	hotelv1Api.Get("/:id/rooms", hotelHandler.HandleGetRooms)

	roomv1Api := apiv1.Group("/rooms", middleware.JWTAuth(userStore))
	roomv1Api.Get("/", roomHandler.HandleGetRooms)
	roomv1Api.Post("/:id/book", roomHandler.HandlerBookRoom)

	//admin bookings
	_ = adminv1Api
	adminbookingv1Api := adminv1Api.Group("/bookings")

	adminbookingv1Api.Get("/", bookingHandler.GetBookings)
	adminbookingv1Api.Get("/:id", bookingHandler.GetBooking)
	adminbookingv1Api.Post("/:id", bookingHandler.CancelBooking)

	userBookingv1Api := apiv1.Group("/bookings")

	userBookingv1Api.Get("/", bookingHandler.GetBookings)
	userBookingv1Api.Get("/:id", bookingHandler.GetBooking)
	userBookingv1Api.Post("/:id/cancel", bookingHandler.CancelBooking)

	app.Listen(fmt.Sprintf(":%s", config.EnvConfig.ServerPort))
}
