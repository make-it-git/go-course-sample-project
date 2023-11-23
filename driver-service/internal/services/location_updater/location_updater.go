package location_updater

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"driver-service/internal/db/repository"
	"driver-service/internal/logger"
	"driver-service/internal/otel"
)

type LocationUpdater struct {
	conn       *redis.Client
	repository OrderRepository
	log        logger.Log
}

func NewLocationUpdater(conn *redis.Client, repository OrderRepository, log logger.Log) LocationUpdater {
	return LocationUpdater{
		conn:       conn,
		repository: repository,
		log:        log,
	}
}

func (l LocationUpdater) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			cmd := l.conn.XRead(ctx, &redis.XReadArgs{
				Streams: []string{"tracks", "$"},
				Block:   time.Second,
			})
			if cmd.Err() != nil {
				if errors.Is(cmd.Err(), redis.Nil) {
					continue
				}
				l.log.WithError(cmd.Err(), "xread")
				continue
			}
			for _, v := range cmd.Val() {
				for _, m := range v.Messages {
					jsonData, err := json.Marshal(m.Values)
					if err != nil {
						l.log.WithError(err, "json marshal")
						continue
					}
					e := LocationEvent{}
					if err = json.Unmarshal(jsonData, &e); err != nil {
						l.log.WithError(err, "json marshal")
						continue
					}
					var span *trace.Span
					spanContext, err := otel.NewSpanContext(e.TraceID, e.SpanID)
					if err != nil {
						l.log.WithError(err, "create span context")
					} else {
						ctx = trace.ContextWithSpanContext(ctx, spanContext)
						spanCtx, newSpan := otel.GetTracer().Start(ctx, "updateCurrentLocation", trace.WithAttributes(attribute.String("orderID", e.Id)))
						ctx = spanCtx
						span = &newSpan
						(*span).AddEvent("update location")
					}
					err = l.repository.UpdateCurrentLocation(ctx, e.Id, repository.Location{
						Latitude:  e.Latitude,
						Longitude: e.Longitude,
					})
					if err != nil {
						if span != nil {
							(*span).SetStatus(codes.Error, "failed update location")
							(*span).RecordError(err)
							(*span).End()
						}
						l.log.WithError(err, "update location")
						continue
					}
					if span != nil {
						(*span).End()
					}
				}
			}
		}
	}
}
