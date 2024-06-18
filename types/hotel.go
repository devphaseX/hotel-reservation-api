package types

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateHotelParams struct {
	Name     string             `json:"name"`
	Location string             `json:"location"`
	RoomID   primitive.ObjectID `json:"roomId"`
}

func (u *UpdateHotelParams) ToBSON() bson.M {
	b := bson.M{}

	if len(u.Name) > 0 {
		b["name"] = u.Name
	}

	if len(u.Location) > 0 {
		b["location"] = u.Location
	}

	if len(u.RoomID) > 0 {
		b["$addToSet"] = bson.M{"rooms": u.RoomID}
	}

	return b
}

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
}

type Room struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Size    string             `bson:"size" json:"type"`
	Seaside bool               `bson:"seaside" json:"seaside"`
	Price   float64            `bson:"price" json:"price"`
	HotelId primitive.ObjectID `bson:"hotelId" json:"hotelId"`
}
