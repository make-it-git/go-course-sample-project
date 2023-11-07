package handlers

import (
	"context"
	"io"

	"google.golang.org/protobuf/types/known/timestamppb"

	ride_order "ride-service/internal/generated/proto/ride.order"
)

type Handler struct {
	ride_order.UnimplementedRideServer
	rideService RideService
}

func NewHandler(rideService RideService) *Handler {
	return &Handler{
		rideService: rideService,
	}
}

func (h Handler) TrackOrder(stream ride_order.Ride_TrackOrderServer) error {
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&ride_order.TrackOrderResponse{})
		}
		if request == nil {
			return stream.SendAndClose(&ride_order.TrackOrderResponse{})
		}
		err = h.rideService.TrackOrder(
			stream.Context(),
			request.Id,
			request.CreatedAt.AsTime(),
			request.Latitude,
			request.Longitude,
		)
		if err != nil {
			return err
		}
	}
}

func (h Handler) GetTrack(ctx context.Context, request *ride_order.GetTrackRequest) (*ride_order.GetTrackResponse, error) {
	track, err := h.rideService.GetTrack(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	trackItems := make([]*ride_order.TrackItem, len(track))
	for i, t := range track {
		trackItems[i] = &ride_order.TrackItem{
			CreatedAt: timestamppb.New(t.Time),
			Latitude:  t.Latitude,
			Longitude: t.Longitude,
		}
	}
	return &ride_order.GetTrackResponse{Track: trackItems}, nil
}
