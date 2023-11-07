package ride

import (
	"context"

	"ride-service/internal/db/repository"
)

type RideRepository interface {
	TrackPoint(ctx context.Context, id string, l repository.Location) error
	GetTrack(ctx context.Context, id string) ([]*repository.Location, error)
}
