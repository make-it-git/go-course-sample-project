package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	rider "rider-service/internal/generated/schema"
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

	order := rider.Order{
		CompletedAt: nil,
		CreatedAt:   openapi_types.Date{h.now()},
		PickupLocation: rider.Location{
			Latitude:  orderData.PickupLocation.Latitude,
			Longitude: orderData.PickupLocation.Longitude,
		},
		Id: uuid.New().String(),
		DropoffLocation: rider.Location{
			Latitude:  orderData.DropoffLocation.Latitude,
			Longitude: orderData.DropoffLocation.Longitude,
		},
		TotalPrice: 0, // TODO
	}
	h.log.Info("Created order", "id", order.Id)
	_ = json.NewEncoder(w).Encode(order)

}
