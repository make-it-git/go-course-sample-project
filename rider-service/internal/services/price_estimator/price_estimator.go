package price_estimator

import (
	"context"
)

type PriceEstimatorService struct {
}

func NewPriceEstimatorService() PriceEstimatorService {
	return PriceEstimatorService{}
}

func (p PriceEstimatorService) Estimate(ctx context.Context, lat1 float32, lng1 float32, lat2 float32, lng2 float32) (int, error) {
	return 100500, nil
}
