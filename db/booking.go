package db

import (
	"context"

	"github.com/devphaseX/hotel-reservation-api/config"
	"github.com/devphaseX/hotel-reservation-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingCollName = "booking"

type BookingStore interface {
	Insert(ctx context.Context, booking *types.Booking) (*types.Booking, error)
	GetBookings(ctx context.Context, filter Record) ([]*types.Booking, error)
	GetBooking(ctx context.Context, filter Record) (*types.Booking, error)
	UpdateBooking(ctx context.Context, filter Record, update Record) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(config.EnvConfig.MongoDBName).Collection(bookingCollName),
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

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter Record) ([]*types.Booking, error) {
	res, err := s.coll.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	var bookings []*types.Booking

	if err := res.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (s *MongoBookingStore) GetBooking(ctx context.Context, filter Record) (*types.Booking, error) {
	res := s.coll.FindOne(ctx, filter)

	var booking *types.Booking

	if err := res.Decode(&booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, filter Record, update Record) error {
	setUpdate := bson.M{"$set": update}
	if _, err := s.coll.UpdateOne(ctx, filter, setUpdate); err != nil {
		return err
	}
	return nil
}
