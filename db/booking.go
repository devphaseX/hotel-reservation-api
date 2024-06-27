package db

import (
	"context"
	"fmt"

	"github.com/devphaseX/hotel-reservation-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingCollName = "booking"

type BookingStore interface {
	Insert(ctx context.Context, booking *types.Booking) (*types.Booking, error)
	GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error)
	GetBooking(ctx context.Context, filter bson.M) (*types.Booking, error)
	UpdateBooking(ctx context.Context, filter bson.M, update bson.M) error
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
	fmt.Println("booking", booking)
	return booking, nil
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
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

func (s *MongoBookingStore) GetBooking(ctx context.Context, filter bson.M) (*types.Booking, error) {
	res := s.coll.FindOne(ctx, filter)

	var booking *types.Booking

	if err := res.Decode(&booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, filter bson.M, update bson.M) error {
	setUpdate := bson.M{"$set": update}
	if _, err := s.coll.UpdateOne(ctx, filter, setUpdate); err != nil {
		return err
	}
	return nil
}
