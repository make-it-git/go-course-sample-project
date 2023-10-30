package handlers

import (
	"encoding/json"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"

	rider "rider-service/internal/generated/schema"
)

func (h *RideImpl) GetOrders(w http.ResponseWriter, r *http.Request, params rider.GetOrdersParams) {
	if params.XUserID <= 0 {
		writeAuthError(w)
		return
	}

	mockOrders := make([]rider.Order, 0, 1)
	mockOrders = append(mockOrders, rider.Order{
		CompletedAt: nil,
		CreatedAt:   openapi_types.Date{h.now()},
		DropoffLocation: rider.Location{
			Latitude:  1,
			Longitude: 1,
		},
		Id: "",
		PickupLocation: rider.Location{
			Latitude:  1,
			Longitude: 1,
		},
		TotalPrice: 0,
	})

	_ = json.NewEncoder(w).Encode(&mockOrders)
}
