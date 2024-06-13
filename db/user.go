package db

import (
	"context"
	"fmt"

	"github.com/devphaseX/hotel-reservation-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const DBNAME = "hotel-reservation"
const userColName = "users"

type UserStore interface {
	GetUserById(ctx context.Context, id string) (*types.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(userColName),
	}
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, fmt.Errorf("provided id not valid: %w", err)
	}
	var user types.User
	err = s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
