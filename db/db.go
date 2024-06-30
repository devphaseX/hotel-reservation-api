package db

const (
	URI         = "mongodb://localhost:27017"
	DBNAME      = "hotel-reservation"
	TEST_DBNAME = "hotel-reservation-test"
)

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

type Record map[string]any
type KeyValueRecord []struct {
	Key   string
	Value any
}
