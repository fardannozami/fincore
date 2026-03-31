package repository

import (
	"testing"

	"github.com/fardannozami/fincore/internal/domain"
	"github.com/fardannozami/fincore/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTransactionRepository_Create(t *testing.T) {
	db := testutil.SetupDB(t, &domain.Transaction{})
	repo := NewTransactionRepository(db)

	tx := db.Begin()
	id := uuid.NewString()
	fromID := uuid.NewString()
	toID := uuid.NewString()

	trx := &domain.Transaction{
		ID:     id,
		FromID: fromID,
		ToID:   toID,
		Amount: 1000,
		Status: "SUCCESS",
	}

	err := repo.Create(tx, trx)
	assert.NoError(t, err)

	tx.Commit()

	// verify ke DB
	var result domain.Transaction
	err = db.First(&result, "id = ?", id).Error

	assert.NoError(t, err)
	assert.Equal(t, id, result.ID)
	assert.Equal(t, fromID, result.FromID)
	assert.Equal(t, toID, result.ToID)
	assert.Equal(t, int64(1000), result.Amount)
	assert.Equal(t, "SUCCESS", result.Status)
}
