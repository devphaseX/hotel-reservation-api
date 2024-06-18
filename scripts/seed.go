package main

import (
	"context"
	"fmt"
	"log"

	"github.com/devphaseX/hotel-reservation-api/db"
	"github.com/devphaseX/hotel-reservation-api/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	hotelStore db.HotelStore
	roomStore  db.RoomStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedUser(firstName, lastName, email, password string) {
	params := types.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}

	user, err := types.NewUserFromParams(params)

	if err != nil {
		log.Fatal(err)
	}

	if _, err = userStore.CreateUser(ctx, user); err != nil {
		log.Fatal(err)
	}

}

func seedHotel(name, location string) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
	}

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	rooms := []*types.Room{
		{
			Size:    "small",
			Price:   99.9,
			HotelId: insertedHotel.ID,
		},
		{
			Size:    "normal",
			Price:   122.9,
			HotelId: insertedHotel.ID,
		},
		{
			Size:    "kingsize",
			Price:   234.9,
			HotelId: insertedHotel.ID,
		},
	}

	for _, room := range rooms {

		insertedRoom, err := roomStore.Insert(ctx, room)

		if err != nil {
			log.Fatal(insertedRoom)
		}

	}
}

func main() {
	seedHotel("Bellucia", "France")
	seedHotel("The cosy Hotel", "The Netherland")
	seedHotel("Don't die in your sleep", "London")
	seedUser("Ayomide", "Lawal", "ayomidelawal800@gmail.com", "supersecure")

	fmt.Println("seeding my db")
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.URI))

	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		fmt.Println(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}
