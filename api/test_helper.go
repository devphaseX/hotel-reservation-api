package api

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/devphaseX/hotel-reservation-api/db"
	"github.com/devphaseX/hotel-reservation-api/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testStore struct {
	client *mongo.Client
	*db.Store
}

func (ts *testStore) tearDown(t *testing.T) {
	if err := ts.client.Database(db.DBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup() *testStore {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.URI))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)

	return &testStore{
		client: client,
		Store: &db.Store{
			User:    db.NewMongoUserStore(client),
			Room:    db.NewMongoRoomStore(client, hotelStore),
			Hotel:   hotelStore,
			Booking: db.NewMongoBookingStore(client),
		},
	}
}

var TestFiberConfig = fiber.Config{
	// Override default error handler
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		if err, ok := err.(utils.Error); ok {
			return ctx.Status(err.Code).JSON(err)
		}

		internalError := utils.NewError(http.StatusInternalServerError, err.Error())
		return ctx.Status(internalError.Code).JSON(internalError)
	},
}
