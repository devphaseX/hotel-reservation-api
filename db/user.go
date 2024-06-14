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
	GetUsers(ctx context.Context) ([]*types.User, error)
	CreateUser(ctx context.Context, user *types.User) (*types.User, error)
	RemoveUser(ctx context.Context, id string) error
	PutUser(ctx context.Context, filter bson.D, update bson.M) error
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

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	users := []*types.User{}

	cur, err := s.coll.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	if err = cur.All(ctx, &users); err != nil {
		return nil, err
	}

	// for cur.Next(ctx) {
	// 	var user *types.User

	// 	err = cur.Decode(&user)

	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	users = append(users, user)
	// }

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

func (s *MongoUserStore) PutUser(ctx context.Context, filter bson.D, values bson.M) error {
	formattedValues := bson.M{}
	for key, value := range values {
		newKey := toSnakeCase(key)
		formattedValues[newKey] = value
	}

	update := bson.D{{Key: "$set", Value: formattedValues}}

	res, err := s.coll.UpdateMany(ctx, filter, update)

	if err != nil {
		return err
	}

	fmt.Println(res.ModifiedCount)

	return nil
}
