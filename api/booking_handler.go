package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/devphaseX/hotel-reservation-api/db"
	"github.com/devphaseX/hotel-reservation-api/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) GetBookings(c *fiber.Ctx) error {
	user, err := getAuthUser(c)

	if err != nil {
		return utils.ErrUnauthorized()

	}

	filter := db.Record{}

	if !user.IsAdmin {
		filter["userId"] = bson.M{"$eq": user.ID}
	}

	bookings, err := h.store.Booking.GetBookings(c.Context(), filter)

	fmt.Println(bookings, err)
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

func (h *BookingHandler) GetBooking(c *fiber.Ctx) error {
	user, err := getAuthUser(c)

	if err != nil {
		return utils.ErrUnauthorized()
	}

	oid, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		fmt.Println("filter oid", oid, err)
		return utils.ErrInvalidID()
	}

	filter := db.Record{"_id": bson.M{"$eq": oid}}

	if !user.IsAdmin {
		filter["userId"] = bson.M{"$eq": user.ID}
	}

	booking, err := h.store.Booking.GetBooking(c.Context(), filter)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return utils.NewError(http.StatusNotFound, "booking not found")
		}

		return errors.New("failed to get booking")
	}

	return c.JSON(booking)
}

func (h *BookingHandler) CancelBooking(c *fiber.Ctx) error {
	user, err := getAuthUser(c)

	if err != nil {
		return utils.ErrUnauthorized()

	}

	oid, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return utils.ErrInvalidID()
	}

	filter := db.Record{"_id": bson.M{"$eq": oid}}

	if !user.IsAdmin {
		filter["userId"] = bson.M{"$eq": user.ID}
	}

	if err = h.store.Booking.UpdateBooking(c.Context(), filter, db.Record{"cancel": true}); err != nil {
		return utils.NewError(http.StatusNotFound, "booking not found")
	}

	return c.JSON(map[string]any{"message": "updated"})
}
