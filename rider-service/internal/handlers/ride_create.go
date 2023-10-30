package handlers

import (
	"encoding/json"
	"net/http"

	rider "rider-service/internal/generated/schema"
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

	createdOrder, err := h.orderService.Create(r.Context(), order.OrderCreate{
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
