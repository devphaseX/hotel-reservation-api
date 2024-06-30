package db

import (
	"context"
	"fmt"

	"github.com/devphaseX/hotel-reservation-api/config"
	"github.com/devphaseX/hotel-reservation-api/types"
	"github.com/devphaseX/hotel-reservation-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColName = "users"

type Dropper interface {
	Drop(ctx context.Context) error
}

type UserStore interface {
	Dropper
	GetUserById(ctx context.Context, id string) (*types.User, error)
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
	GetUsers(ctx context.Context) ([]*types.User, error)
	CreateUser(ctx context.Context, user *types.User) (*types.User, error)
	RemoveUser(ctx context.Context, id string) error
	PutUser(ctx context.Context, filter Record, update types.UpdateUserParams) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(config.EnvConfig.MongoDBName).Collection(userColName),
	}
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("----dropping user collection")
	return s.coll.Drop(ctx)
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, utils.ErrInvalidID()
	}

	var user types.User
	if err = s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	users := []*types.User{}

	cur, err := s.coll.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	if err = cur.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *MongoUserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	user.ID = res.InsertedID.(primitive.ObjectID)

	return user, nil
}

func (s *MongoUserStore) RemoveUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	//TODO:log the deleted user to see if an actual record got removed for this query
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})

	if err != nil {
		return err
	}

	return nil
}

func (s *MongoUserStore) PutUser(ctx context.Context, filter Record, values types.UpdateUserParams) error {

	update := bson.D{{Key: "$set", Value: values.ToBSON()}}

	res, err := s.coll.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	fmt.Println(res.ModifiedCount)

	return nil
}
