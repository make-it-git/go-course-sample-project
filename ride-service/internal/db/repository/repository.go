package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRideRepository(conn *redis.Client) RideRepositoryImpl {
	return RideRepositoryImpl{conn: conn}
}

type RideRepositoryImpl struct {
	conn *redis.Client
}

func (r RideRepositoryImpl) TrackPoint(ctx context.Context, id string, l Location) error {
	return r.conn.LPush(ctx, id, &l).Err()
}

func (r RideRepositoryImpl) GetTrack(ctx context.Context, id string) ([]*Location, error) {
	cmd := r.conn.LRange(ctx, id, 0, -1)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	var result []*Location
	err := cmd.ScanSlice(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
