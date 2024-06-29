package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/devphaseX/hotel-reservation-api/db"
	"github.com/devphaseX/hotel-reservation-api/types"
	"github.com/devphaseX/hotel-reservation-api/utils"
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

func (p *BookRoomParams) Validate() error {
	now := time.Now()

	if now.After(p.FromDate) || p.ToDate.Before(p.FromDate) {
		return errors.New("cannot book a room in the past")
	}

	if p.NumPersons < 1 {
		return errors.New("require to provide a single person for the room")
	}

	return nil
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{store: store}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	_, ok := c.Context().Value("user").(*types.User)

	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(utils.FailedResp{
			Type:    "error",
			Message: "Unauthorized",
		})
	}

	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})

	if err != nil {
		return err
	}

	return c.JSON(rooms)
}

func (h *RoomHandler) HandlerBookRoom(c *fiber.Ctx) error {
	user, ok := c.Context().Value("user").(*types.User)

	if !ok {
		return utils.ErrUnauthorized()
	}

	oid, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return utils.ErrInvalidID()
	}

	var params BookRoomParams

	if err := c.QueryParser(&params); err != nil {
		return utils.ErrBadJSON()
	}

	if err := params.Validate(); err != nil {
		return utils.NewValidationError(err)

	}

	available, err := h.RoomAvailable(c.Context(), oid, params)

	if err != nil {
		return errors.New("failed to book room")

	}

	if !available {
		return utils.NewError(http.StatusConflict, "room already booked")
	}

	room, err := h.store.Room.GetRoom(c.Context(), bson.M{"_id": oid})

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return utils.NewError(http.StatusNotFound, "room not found")
		}

		return errors.New("failed to set a booking")
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

func (h *RoomHandler) RoomAvailable(ctx context.Context, roomId primitive.ObjectID, params BookRoomParams) (bool, error) {

	bookedRoomFilter := bson.M{
		"fromDate": bson.M{"$gte": params.FromDate},
		"toDate":   bson.M{"$lte": params.ToDate},
		"roomId":   bson.M{"$eq": roomId},
	}

	bookedRooms, err := h.store.Booking.GetBookings(ctx, bookedRoomFilter)
	if err != nil {
		return false, err
	}

	return len(bookedRooms) == 0, nil
}
