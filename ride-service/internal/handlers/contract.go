package handlers

import (
	"context"
	"time"

	"ride-service/internal/services/ride"
)

type RideService interface {
	TrackOrder(ctx context.Context, id string, t time.Time, latitude float32, longitude float32) error
	GetTrack(ctx context.Context, id string) ([]ride.Location, error)
}
