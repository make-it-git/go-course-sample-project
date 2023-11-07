package repository

import (
	"encoding/json"
	"time"
)

type Location struct {
	CreatedAt time.Time `json:"time" redis:"time"`
	Latitude  float32   `json:"latitude" redis:"latitude"`
	Longitude float32   `json:"longitude" redis:"longitude"`
}

func (l *Location) MarshalBinary() (data []byte, err error) {
	return json.Marshal(l)
}

func (l *Location) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, l)
}
