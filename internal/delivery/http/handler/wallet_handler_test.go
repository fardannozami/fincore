package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"github.com/fardannozami/fincore/internal/domain"
	"github.com/fardannozami/fincore/internal/repository"
	"github.com/fardannozami/fincore/internal/testutil"
)

// setup DB test
func setupWalletHandlerTest(t *testing.T) (*WalletHandler, *gorm.DB) {
	db := testutil.SetupDB(t, &domain.Wallet{})

	repo := repository.NewWalletRepository(db)
	handler := NewWalletHandler(db, repo)

	return handler, db
}

func TestWalletHandler_CreateWallet(t *testing.T) {
	handler, db := setupWalletHandlerTest(t)

	// buat request HTTP
	req := httptest.NewRequest(http.MethodPost, "/wallet", nil)
	w := httptest.NewRecorder()

	// call handler
	handler.CreateWallet(w, req)

	res := w.Result()
	defer res.Body.Close()

	// cek status code
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// decode response
	var response map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)

	// cek response field
	assert.NotEmpty(t, response["id"])
	assert.Equal(t, "user-1", response["user_id"])
	assert.Equal(t, float64(0), response["balance"]) // JSON number = float64

	// cek benar-benar masuk DB
	var wallet domain.Wallet
	err = db.First(&wallet, "id = ?", response["id"]).Error
	assert.NoError(t, err)

	assert.Equal(t, "user-1", wallet.UserID)
	assert.Equal(t, int64(0), wallet.Balance)
}
