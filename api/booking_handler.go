package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/devphaseX/hotel-reservation-api/db"
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
		return c.Status(http.StatusUnauthorized).JSON(FailedResp{Type: "error", Message: "unauthorized"})

	}

	filter := bson.M{}

	if !user.IsAdmin {
		filter["userId"] = bson.M{"$eq": user.ID}
	}

	bookings, err := h.store.Booking.GetBookings(c.Context(), filter)

	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

func (h *BookingHandler) GetBooking(c *fiber.Ctx) error {
	user, err := getAuthUser(c)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(FailedResp{Type: "error", Message: "unauthorized"})

	}

	oid, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(FailedResp{Type: "error", Message: "not a valid id"})
	}

	filter := bson.M{"_id": bson.M{"$eq": oid}}

	if !user.IsAdmin {
		filter["userId"] = bson.M{"$eq": user.ID}
	}

	booking, err := h.store.Booking.GetBooking(c.Context(), filter)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(http.StatusNotFound).JSON(FailedResp{Type: "error", Message: "booking not found"})
		}

		fmt.Println(err)

		return errors.New("failed to get booking")
	}

	return c.JSON(booking)
}

func (h *BookingHandler) CancelBooking(c *fiber.Ctx) error {
	user, err := getAuthUser(c)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(FailedResp{Type: "error", Message: "unauthorized"})

	}

	oid, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(FailedResp{Type: "error", Message: "not a valid id"})
	}

	filter := bson.M{"_id": bson.M{"$eq": oid}}

	if !user.IsAdmin {
		filter["userId"] = bson.M{"$eq": user.ID}
	}

	if err = h.store.Booking.UpdateBooking(c.Context(), filter, bson.M{"cancel": true}); err != nil {
		return c.Status(http.StatusNotFound).JSON(FailedResp{Type: "error", Message: "not found"})
	}

	return c.JSON(map[string]any{"message": "updated"})
}
