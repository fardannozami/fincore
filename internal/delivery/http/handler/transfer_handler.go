package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fardannozami/fincore/internal/delivery/http/dto"
	"github.com/google/uuid"
)

type TransferUsecase interface {
	Transfer(fromID, toID string, amount int64, refID string) error
}

type TransferHandler struct {
	usecase TransferUsecase
}

func NewTransferHandler(u TransferUsecase) *TransferHandler {
	return &TransferHandler{u}
}

func (h *TransferHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req dto.TransferRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if req.FromID == "" || req.ToID == "" || req.Amount <= 0 {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	err := h.usecase.Transfer(
		req.FromID,
		req.ToID,
		req.Amount,
		uuid.NewString(),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(dto.TransferResponse{
		Message: "success",
	})
}
