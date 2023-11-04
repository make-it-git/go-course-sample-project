package driver_search

import (
	"context"
	"errors"

	"golang.org/x/exp/rand"
)

type DriverSearchService struct {
}

var ErrDriverNotFound = errors.New("can not find driver")

func NewDriverSearchService() DriverSearchService {
	return DriverSearchService{}
}

func (d DriverSearchService) FindDriver(ctx context.Context, lat1 float32, lng1 float32, lat2 float32, lng2 float32) (int64, error) {
	if rand.Int()%10 > 5 {
		return 0, ErrDriverNotFound
	}
	return int64(rand.Int() % 100), nil
}
