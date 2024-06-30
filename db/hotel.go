package db

import (
	"context"

	"github.com/devphaseX/hotel-reservation-api/config"
	"github.com/devphaseX/hotel-reservation-api/types"
	"github.com/devphaseX/hotel-reservation-api/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const hotelCollName = "hotel"

type GetHotelsQueryParam struct {
	Room   *bool `json:"room"`
	Rating *int
}

type ReceivedFilter map[string]any

func (q *GetHotelsQueryParam) GetReceivedFilter() (receivedFilter ReceivedFilter, includeRoom bool) {
	receivedFilter = map[string]any{}

	if q.Rating != nil {
		receivedFilter["rating"] = q.Rating
	}

	if q.Room != nil {
		includeRoom = *q.Room
	}

	return
}

func NewGetHotelsQueryParams(c *fiber.Ctx) (*GetHotelsQueryParam, error) {
	var query GetHotelsQueryParam

	if err := c.QueryParser(&query); err != nil {
		return nil, err
	}

	return &query, nil
}

type HotelStore interface {
	Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error)
	Update(context.Context, Record, types.UpdateHotelParams) error
	GetMany(ctx context.Context, paginateQuery *utils.PaginateQuery, filter *GetHotelsQueryParam) ([]*types.Hotel, error)
	GetOne(ctx context.Context, id string) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(config.EnvConfig.MongoDBName).Collection(hotelCollName),
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

func (h *MongoHotelStore) Update(ctx context.Context, filter Record, values types.UpdateHotelParams) error {
	_, err := h.coll.UpdateOne(ctx, filter, values.ToBSON())
	if err != nil {
		return err
	}

	return nil
}

func (h *MongoHotelStore) GetMany(ctx context.Context, paginateQuery *utils.PaginateQuery, filter *GetHotelsQueryParam) ([]*types.Hotel, error) {
	opts := options.FindOptions{}
	offset, limit := paginateQuery.ApplyPagination()

	opts.SetSkip(offset)
	opts.SetLimit(limit)

	receivedFilter, includeRoom := filter.GetReceivedFilter()

	_ = includeRoom
	res, err := h.coll.Find(ctx, receivedFilter, &opts)

	if err != nil {
		return nil, err
	}

	hotels := []*types.Hotel{}
	if err = res.All(ctx, &hotels); err != nil {
		return nil, err
	}

	return hotels, nil
}

func (h *MongoHotelStore) GetOne(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, utils.ErrInvalidID()
	}

	res := h.coll.FindOne(ctx, bson.M{"_id": oid})

	var hotel types.Hotel
	if err := res.Decode(&hotel); err != nil {
		return nil, err
	}

	return &hotel, nil
}
