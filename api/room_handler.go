package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/devphaseX/hotel-reservation-api/db"
	"github.com/devphaseX/hotel-reservation-api/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomHandler struct {
	store *db.Store
}

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	ToDate     time.Time `json:"toDate"`
	NumPersons int       `json:"numPersons"`
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{store: store}
}

func (h *RoomHandler) HandlerBookRoom(c *fiber.Ctx) error {
	user, ok := c.Context().Value("user").(*types.User)

	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(FailedResp{
			Type:    "error",
			Message: "Unauthorized",
		})
	}

	oid, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return errors.New("not valid id")
	}

	room, err := h.store.Room.GetRoom(c.Context(), bson.M{"_id": oid})

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("room not found")
		}

		fmt.Printf("failed to get room: %v\n", err)
		return errors.New("failed to set a booking")
	}

	var params BookRoomParams

	if err := c.QueryParser(&params); err != nil {
		return err
	}

	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     room.ID,
		FromDate:   params.FromDate,
		ToDate:     params.ToDate,
		NumPersons: params.NumPersons,
	}

	_, err = h.store.Booking.Insert(c.Context(), &booking)

	if err != nil {
		return err
	}

	return c.JSON(booking)
}
