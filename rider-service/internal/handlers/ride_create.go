package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	rider "rider-service/internal/generated/schema"
	otel2 "rider-service/internal/otel"
	"rider-service/internal/services/order"
)

func (h *RideImpl) PostOrders(w http.ResponseWriter, r *http.Request, params rider.PostOrdersParams) {
	if params.XUserID <= 0 {
		writeAuthError(w)
		return
	}
	orderData := rider.CreateOrder{}
	if err := json.NewDecoder(r.Body).Decode(&orderData); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	ctx, span := otel2.GetTracer().Start(r.Context(), "createOrder", trace.WithAttributes(attribute.String("userID", strconv.Itoa(params.XUserID))))
	defer span.End()

	createdOrder, err := h.orderService.Create(ctx, order.OrderCreate{
		PickupLocation: order.Location{
			Latitude:  orderData.PickupLocation.Latitude,
			Longitude: orderData.PickupLocation.Longitude,
		},
		DropoffLocation: order.Location{
			Latitude:  orderData.DropoffLocation.Latitude,
			Longitude: orderData.DropoffLocation.Longitude,
		},
		UserID:         params.XUserID,
		IdempotencyKey: orderData.IdempotencyKey,
	})

	if err != nil {
		span.SetStatus(codes.Error, "failed create order")
		span.RecordError(err)
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	responseOrder := rider.Order{
		CreatedAt: createdOrder.CreatedAt.String(),
		DropoffLocation: rider.Location{
			Latitude:  createdOrder.DropoffLocation.Latitude,
			Longitude: createdOrder.DropoffLocation.Longitude,
		},
		Id: createdOrder.ID,
		PickupLocation: rider.Location{
			Latitude:  createdOrder.PickupLocation.Latitude,
			Longitude: createdOrder.PickupLocation.Longitude,
		},
		TotalPrice: createdOrder.TotalPrice,
	}
	h.log.Info("Created order", "id", responseOrder.Id)
	_ = json.NewEncoder(w).Encode(responseOrder)
}
