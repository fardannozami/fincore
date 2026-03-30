package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fardannozami/fincore/internal/delivery/http/dto"
)

// ===== MOCK USECASE =====
type mockUsecase struct {
	err error
}

func (m *mockUsecase) Transfer(fromID, toID string, amount int64, refID string) error {
	return m.err
}

// ===== TEST SUCCESS =====
func TestTransferHandler_Success(t *testing.T) {
	mock := &mockUsecase{err: nil}
	handler := NewTransferHandler(mock)

	reqBody := dto.TransferRequest{
		FromID: "A",
		ToID:   "B",
		Amount: 100,
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/transfer", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Transfer(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var resp dto.TransferResponse
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Message)
}

// ===== TEST INVALID JSON =====
func TestTransferHandler_InvalidBody(t *testing.T) {
	mock := &mockUsecase{}
	handler := NewTransferHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/transfer", bytes.NewBuffer([]byte("invalid")))
	w := httptest.NewRecorder()

	handler.Transfer(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ===== TEST INVALID INPUT =====
func TestTransferHandler_InvalidInput(t *testing.T) {
	mock := &mockUsecase{}
	handler := NewTransferHandler(mock)

	reqBody := dto.TransferRequest{
		FromID: "",
		ToID:   "B",
		Amount: 0,
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/transfer", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Transfer(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ===== TEST USECASE ERROR =====
func TestTransferHandler_UsecaseError(t *testing.T) {
	mock := &mockUsecase{
		err: errors.New("insufficient balance"),
	}
	handler := NewTransferHandler(mock)

	reqBody := dto.TransferRequest{
		FromID: "A",
		ToID:   "B",
		Amount: 1000,
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/transfer", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Transfer(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
