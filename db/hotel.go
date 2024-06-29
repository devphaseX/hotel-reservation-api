package db

import (
	"context"

	"github.com/devphaseX/hotel-reservation-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const hotelCollName = "hotel"

type HotelStore interface {
	Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error)
	Update(context.Context, bson.M, types.UpdateHotelParams) error
	GetMany(ctx context.Context) ([]*types.Hotel, error)
	GetOne(ctx context.Context, id primitive.ObjectID) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(hotelCollName),
	}
}

func (h *MongoHotelStore) Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := h.coll.InsertOne(ctx, hotel)

	if err != nil {
		return nil, err
	}

	hotel.ID = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (h *MongoHotelStore) Update(ctx context.Context, filter bson.M, values types.UpdateHotelParams) error {
	_, err := h.coll.UpdateOne(ctx, filter, values.ToBSON())
	if err != nil {
		return err
	}

	return nil
}

func (h *MongoHotelStore) GetMany(ctx context.Context) ([]*types.Hotel, error) {
	res, err := h.coll.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	hotels := []*types.Hotel{}
	if err = res.All(ctx, &hotels); err != nil {
		return nil, err
	}

	return hotels, nil
}

func (h *MongoHotelStore) GetOne(ctx context.Context, id primitive.ObjectID) (*types.Hotel, error) {
	res := h.coll.FindOne(ctx, bson.M{"_id": id})

	var hotel types.Hotel
	if err := res.Decode(&hotel); err != nil {
		return nil, err
	}

	return &hotel, nil
}
