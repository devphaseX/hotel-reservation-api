package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/devphaseX/hotel-reservation-api/db"
	"github.com/devphaseX/hotel-reservation-api/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandlerGets(c *fiber.Ctx) error {
	pq, err := utils.NewPaginate(c)
	if err != nil {
		fmt.Println(err)
		return utils.NewError(http.StatusBadRequest, "invalid query")
	}

	filterQuery, err := db.NewGetHotelsQueryParams(c)

	if err != nil {
		fmt.Println(err)
		return utils.NewError(http.StatusBadRequest, "invalid query")
	}

	fmt.Println(pq, filterQuery)

	hotels, err := h.store.Hotel.GetMany(c.Context(), pq, filterQuery)
	fmt.Println(hotels)

	if err != nil {
		return err
	}

	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGet(c *fiber.Ctx) error {

	hotel, err := h.store.Hotel.GetOne(c.Context(), c.Params("id"))

	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return utils.NewError(http.StatusNotFound, "hotel not found")
		}
		return err
	}

	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return utils.ErrInvalidID()
	}

	filters := db.Record{"hotelId": oid}

	rooms, err := h.store.Room.GetRooms(c.Context(), filters)

	if err != nil {
		return err
	}

	return c.JSON(rooms)
}
