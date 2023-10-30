package handlers

import (
	"encoding/json"
	"net/http"

	rider "rider-service/internal/generated/schema"
	"rider-service/internal/logger"
	"rider-service/internal/now_time"
)

type RideImpl struct {
	log logger.Log
	now now_time.NowType
}

var _ rider.ServerInterface = (*RideImpl)(nil)

func New(log logger.Log, now now_time.NowType) *RideImpl {
	return &RideImpl{
		log: log,
		now: now,
	}
}

func writeError(w http.ResponseWriter, statusCode int, err error) {
	message := err.Error()
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(rider.Error{Message: &message})
}

func writeAuthError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
}
