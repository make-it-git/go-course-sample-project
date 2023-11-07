package handlers

import (
	"time"
)

type Location struct {
	CreatedAt time.Time `redis:"time"`
	Latitude  float32   `redis:"latitude"`
	Longitude float32   `redis:"longitude"`
}
