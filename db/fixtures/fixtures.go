package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/devphaseX/hotel-reservation-api/db"
	"github.com/devphaseX/hotel-reservation-api/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BookRoom(store *db.Store, userId, roomId primitive.ObjectID, numPerson int, fromDate, tillDate time.Time, cancel bool) *types.Booking {
	booking := types.Booking{
		UserID:   userId,
		RoomID:   roomId,
		FromDate: fromDate,
		ToDate:   tillDate,
		Cancel:   cancel,
	}
	insertedBooking, err := store.Booking.Insert(context.TODO(), &booking)

	if err != nil {
		log.Fatal(err)
	}

	return insertedBooking
}

func AddRoom(store *db.Store, size string, ss bool, price float64, hotelId primitive.ObjectID) *types.Room {
	room := types.Room{
		Size:    size,
		Seaside: ss,
		Price:   price,
		HotelId: hotelId,
	}

	insertedRoom, err := store.Room.Insert(context.TODO(), &room)

	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func AddHotel(store *db.Store, name, loc string, roomIds []primitive.ObjectID) *types.Hotel {
	if roomIds == nil {
		roomIds = []primitive.ObjectID{}
	}

	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomIds,
	}

	insertedHotel, err := store.Hotel.Insert(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func AddUser(store *db.Store, firstName, lastName string, admin bool) *types.User {
	params := types.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     fmt.Sprintf("%s.%s@mail.com", firstName, lastName),
		Password:  fmt.Sprintf("%s_%s", firstName, lastName),
	}

	user, err := types.NewUserFromParams(params)

	user.IsAdmin = admin
	if err != nil {
		log.Fatal(err)
	}

	insertedUser, err := store.User.CreateUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}
