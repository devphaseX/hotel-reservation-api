package db

import (
	"context"

	"github.com/devphaseX/hotel-reservation-api/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingCollName = "booking"

type BookingStore interface {
	Insert(ctx context.Context, booking *types.Booking) (*types.Booking, error)
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(bookingCollName),
	}
}

func (s *MongoBookingStore) Insert(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	res, err := s.coll.InsertOne(ctx, booking)

	if err != nil {
		return nil, err
	}

	booking.ID = res.InsertedID.(primitive.ObjectID)

	return booking, nil
}
