package db

import (
	"context"

	"github.com/devphaseX/hotel-reservation-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomCollName = "room"

type RoomStore interface {
	Insert(ctx context.Context, hotel *types.Room) (*types.Room, error)
	GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(DBNAME).Collection(roomCollName),
		HotelStore: hotelStore,
	}
}

func (h *MongoRoomStore) Insert(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := h.coll.InsertOne(ctx, room)

	if err != nil {
		return nil, err
	}

	room.ID = res.InsertedID.(primitive.ObjectID)

	filter := bson.M{"_id": room.HotelId}
	update := types.UpdateHotelParams{RoomID: room.ID}

	if err = h.HotelStore.Update(ctx, filter, update); err != nil {
		return nil, err
	}

	return room, nil
}

func (h *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	res, err := h.coll.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	var rooms []*types.Room

	if err := res.All(ctx, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}
