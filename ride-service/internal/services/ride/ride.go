package ride

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"ride-service/internal/db/repository"
	"ride-service/internal/otel"
)

type RideService struct {
	rideRepository RideRepository
	conn           *redis.Client
}

func NewRideService(rideRepository RideRepository, conn *redis.Client) RideService {
	return RideService{rideRepository: rideRepository, conn: conn}
}

func (s RideService) TrackOrder(ctx context.Context, id string, t time.Time, latitude float32, longitude float32) error {
	ctx, span := otel.GetTracer().Start(ctx, "trackOrder", trace.WithAttributes(attribute.String("orderID", id)))
	defer span.End()

	err := s.rideRepository.TrackPoint(ctx, id, repository.Location{
		CreatedAt: t,
		Latitude:  latitude,
		Longitude: longitude,
	})
	if err != nil {
		span.SetStatus(codes.Error, "failed track order")
		span.RecordError(err)
		return err
	}

	return s.conn.XAdd(ctx, &redis.XAddArgs{
		Stream: "tracks",
		MaxLen: 1000,
		Approx: true,
		Values: map[string]interface{}{
			"lat":     latitude,
			"lng":     longitude,
			"id":      id,
			"time":    t.Unix(),
			"traceID": span.SpanContext().TraceID().String(),
			"spanID":  span.SpanContext().SpanID().String(),
		},
	}).Err()
}

func (s RideService) GetTrack(ctx context.Context, id string) ([]Location, error) {
	track, err := s.rideRepository.GetTrack(ctx, id)
	if err != nil {
		return nil, err
	}

	result := make([]Location, len(track))
	for i, t := range track {
		result[i] = Location{
			Time:      t.CreatedAt,
			Latitude:  t.Latitude,
			Longitude: t.Longitude,
		}
	}

	return result, nil
}
