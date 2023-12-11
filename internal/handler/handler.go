package handler

import (
	"cart-backend/internal/service"
	"context"
	"encoding/json"
	"net/http"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service}
}

func (h *Handler) ListTxRecordByAddress(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "address is required", http.StatusBadRequest)
		return
	}
	txRecords, err := h.service.List(ctx, address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(txRecords)
}
