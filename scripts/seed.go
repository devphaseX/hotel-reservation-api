package main

import (
	"context"
	"fmt"
	"log"

	"github.com/devphaseX/hotel-reservation-api/api"
	"github.com/devphaseX/hotel-reservation-api/db"
	"github.com/devphaseX/hotel-reservation-api/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.URI))

	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		fmt.Println(err)
	}

	hotelStore := db.NewMongoHotelStore(client)

	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Hotel:   hotelStore,
		Room:    db.NewMongoRoomStore(client, hotelStore),
	}

	user1 := fixtures.AddUser(store, "Ayomide", "Lawal", true)
	userToken, _ := api.CreateTokenClaim(user1)

	fmt.Printf("%s:%s -> %s\n", user1.FirstName, user1.ID, userToken)
	hotel1 := fixtures.AddHotel(store, "Don't die in your sleep", "London", nil)
	fmt.Printf("%s -> %s\n", hotel1.Name, hotel1.ID)

	fixtures.AddRoom(store, "size", true, 89.9, hotel1.ID)
	fixtures.AddRoom(store, "medium", true, 189.9, hotel1.ID)
	room1 := fixtures.AddRoom(store, "size", false, 289.9, hotel1.ID)

	fmt.Printf("%s -> %s\n", room1.Size, room1.ID)

	fmt.Println("seeding my db")
}
