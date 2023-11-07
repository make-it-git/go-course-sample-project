package location_updater

import (
	"context"

	"driver-service/internal/db/repository"
)

type OrderRepository interface {
	UpdateCurrentLocation(ctx context.Context, id string, l repository.Location) error
}
