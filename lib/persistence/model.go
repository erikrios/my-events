package persistence

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	First    string             `json:"first"`
	Last     string             `json:"last"`
	Age      int                `json:"age"`
	Bookings []Booking          `json:"bookings"`
}

func (u *User) String() string {
	return fmt.Sprintf("id: %s, first_name: %s, last_name: %s, Age: %d, Bookings: %v", u.ID, u.First, u.Last, u.Age, u.Bookings)
}

type Booking struct {
	Date    int64  `json:"date"`
	EventID []byte `json:"eventID"`
	Seats   int    `json:"seats"`
}

type Event struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `json:"name"`
	Duration  int                `json:"duration"`
	StartDate int64              `json:"startDate"`
	EndDate   int64              `json:"endDate"`
	Location  Location           `json:"location"`
}

type Location struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `json:"name"`
	Address   string             `json:"address"`
	Country   string             `json:"country"`
	OpenTime  int                `json:"openTime"`
	CloseTime int                `json:"closeTime"`
	Halls     []Hall             `json:"halls"`
}

type Hall struct {
	Name     string `json:"name"`
	Location string `json:"location,omitempty"`
	Capacity int    `json:"capacity"`
}
