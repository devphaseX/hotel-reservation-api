package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/devphaseX/hotel-reservation-api/api/middleware"
	"github.com/devphaseX/hotel-reservation-api/db/fixtures"
	"github.com/devphaseX/hotel-reservation-api/types"
	"github.com/gofiber/fiber/v2"
)

func TestGettBookings(t *testing.T) {
	db := setup()

	defer db.tearDown(t)

	var (
		ownerUser = fixtures.AddUser(db.Store, "james", "foo", true)
		hotel     = fixtures.AddHotel(db.Store, "lafresh", "lagos", nil)
		room      = fixtures.AddRoom(db.Store, "small", false, 234, hotel.ID)

		from           = time.Now()
		till           = from.AddDate(0, 0, 5)
		booking        = fixtures.BookRoom(db.Store, ownerUser.ID, room.ID, 2, from, till, false)
		app            = fiber.New()
		bookingHandler = NewBookingHandler(db.Store)
		admin          = app.Group("/", middleware.JWTAuth(db.Store.User), middleware.AdminAuth)
	)

	_ = booking

	admin.Get("/", bookingHandler.GetBookings)

	req := httptest.NewRequest("GET", "/", nil)
	token, err := CreateTokenClaim(ownerUser)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("x-api-token", token)

	res, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res.StatusCode)

	if res.StatusCode != http.StatusOK {

		var failedResp any
		if err = json.NewDecoder(res.Body).Decode(&failedResp); err != nil {
			t.Fatal(err)
		}

		fmt.Println(failedResp)

		return
	}

	var bookings []*types.Booking

	if err = json.NewDecoder(res.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	fmt.Println(bookings)

	if len(bookings) != 1 {
		t.Fatalf("expects 1 booking but got %d\n", len(bookings))
	}

	have := bookings[0]

	if have.ID != booking.ID {
		t.Fatalf("expect %s got %s", booking.ID, have.ID)
	}

	if have.UserID != booking.UserID {
		t.Fatalf("expect %s got %s", booking.ID, have.ID)
	}

	//booking can be access by owner

	owner := app.Group("owner/", middleware.JWTAuth(db.Store.User))

	owner.Get("/", bookingHandler.GetBookings)
	{

		nonAdmin := fixtures.AddUser(db.Store, "james", "foo", false)
		req := httptest.NewRequest("GET", "/", nil)
		token, err := CreateTokenClaim(nonAdmin)

		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("x-api-token", token)

		res, err := app.Test(req)

		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode == http.StatusOK {

			var failedResp any
			if err = json.NewDecoder(res.Body).Decode(&failedResp); err != nil {
				t.Fatal("non admin user are allow access to non owned booking")
			}

			fmt.Println(failedResp)

			return
		}
	}

}

func TestUserGetBooking(t *testing.T) {
	db := setup()

	defer db.tearDown(t)

	var (
		nonOwnerUser = fixtures.AddUser(db.Store, "bad", "hacker", false)
		ownerUser    = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel        = fixtures.AddHotel(db.Store, "lafresh", "lagos", nil)
		room         = fixtures.AddRoom(db.Store, "small", false, 234, hotel.ID)

		from           = time.Now()
		till           = from.AddDate(0, 0, 5)
		booking        = fixtures.BookRoom(db.Store, ownerUser.ID, room.ID, 2, from, till, false)
		app            = fiber.New()
		bookingHandler = NewBookingHandler(db.Store)
	)

	app.Get("/:id", middleware.JWTAuth(db.Store.User), bookingHandler.GetBooking)

	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	token, err := CreateTokenClaim(ownerUser)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("x-api-token", token)

	res, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res.StatusCode)

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expects status code %d but got %d", http.StatusOK, res.StatusCode)
	}

	var ownedBooking types.Booking

	if err = json.NewDecoder(res.Body).Decode(&ownedBooking); err != nil {
		t.Fatal(err)
	}

	if ownedBooking.ID != booking.ID {
		t.Fatalf("expect %s got %s", booking.ID, ownedBooking.ID)
	}

	if ownedBooking.UserID != booking.UserID {
		t.Fatalf("expect %s got %s", booking.ID, ownedBooking.ID)
	}

	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	token, err = CreateTokenClaim(nonOwnerUser)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("x-api-token", token)

	res, err = app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusNotFound {
		fmt.Println(res.StatusCode)
		t.Fatalf("owner booking not to be access to other user")
	}

}
