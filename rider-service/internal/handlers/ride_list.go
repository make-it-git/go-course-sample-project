package handlers

import (
	"encoding/json"
	"net/http"

	rider "rider-service/internal/generated/schema"
)

func (h *RideImpl) GetOrders(w http.ResponseWriter, r *http.Request, params rider.GetOrdersParams) {
	if params.XUserID <= 0 {
		writeAuthError(w)
		return
	}

	orders, err := h.orderService.List(r.Context(), params.XUserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	resultOrders := make([]rider.Order, 0, len(orders))
	for _, order := range orders {
		resultOrders = append(resultOrders, rider.Order{
			CreatedAt: order.CreatedAt.String(),
			DropoffLocation: rider.Location{
				Latitude:  order.DropoffLocation.Latitude,
				Longitude: order.DropoffLocation.Longitude,
			},
			Id: order.ID,
			PickupLocation: rider.Location{
				Latitude:  order.PickupLocation.Latitude,
				Longitude: order.PickupLocation.Longitude,
			},
			TotalPrice: order.TotalPrice,
		})
	}

	_ = json.NewEncoder(w).Encode(&orders)
}
