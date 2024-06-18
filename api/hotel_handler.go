package api

import (
	"fmt"

	"github.com/devphaseX/hotel-reservation-api/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

type GetHotelsQueryParam struct {
	Room   bool `json:"room"`
	Rating int
}

func (h *HotelHandler) HandlerGets(c *fiber.Ctx) error {
	var query GetHotelsQueryParam

	if err := c.QueryParser(&query); err != nil {
		return err
	}

	fmt.Println(query)

	hotels, err := h.store.Hotel.GetMany(c.Context())

	if err != nil {
		return err
	}

	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGet(c *fiber.Ctx) error {

	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	hotel, err := h.store.Hotel.GetOne(c.Context(), oid)

	if err != nil {
		return err
	}

	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	filters := bson.M{"hotelId": oid}

	rooms, err := h.store.Room.GetRooms(c.Context(), filters)

	if err != nil {
		return err
	}

	return c.JSON(rooms)
}
