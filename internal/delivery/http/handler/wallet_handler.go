package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fardannozami/fincore/internal/delivery/http/dto"
	"github.com/fardannozami/fincore/internal/domain"
	"github.com/fardannozami/fincore/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletHandler struct {
	db   *gorm.DB
	repo *repository.WalletRepository
}

func NewWalletHandler(db *gorm.DB, repo *repository.WalletRepository) *WalletHandler {
	return &WalletHandler{db, repo}
}

func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {

	wallet := domain.Wallet{
		ID:      uuid.NewString(),
		UserID:  "user-1", // nanti bisa dari auth
		Balance: 0,
	}

	if err := h.repo.Create(h.db, &wallet); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dto.CreateWalletResponse{
		ID:      wallet.ID,
		UserID:  wallet.UserID,
		Balance: wallet.Balance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
