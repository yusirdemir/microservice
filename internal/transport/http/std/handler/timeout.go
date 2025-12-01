package handler

import (
	"net/http"
	"time"
)

type TimeoutHandler struct{}

func NewTimeoutHandler() *TimeoutHandler {
	return &TimeoutHandler{}
}

func (h *TimeoutHandler) Timeout(w http.ResponseWriter, r *http.Request) {
	time.Sleep(6 * time.Second)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Finished successfully"}`))
}
