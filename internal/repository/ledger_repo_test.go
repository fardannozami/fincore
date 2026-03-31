package repository

import (
	"testing"

	"github.com/fardannozami/fincore/internal/domain"
	"github.com/fardannozami/fincore/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLedgerRepository_Create(t *testing.T) {
	db := testutil.SetupDB(t, &domain.Ledger{})
	repo := NewLedgerRepository(db)

	tx := db.Begin()
	id := uuid.NewString()
	walletID := uuid.NewString()
	refID := uuid.NewString()

	ledger := &domain.Ledger{
		ID:       id,
		WalletID: walletID,
		Amount:   1000,
		Type:     "CREDIT",
		RefID:    refID,
	}

	err := repo.Create(tx, ledger)
	assert.NoError(t, err)

	tx.Commit()

	// verify ke DB
	var result domain.Ledger
	err = db.First(&result, "id = ?", id).Error

	assert.NoError(t, err)
	assert.Equal(t, int64(1000), result.Amount)
	assert.Equal(t, "CREDIT", result.Type)
	assert.Equal(t, walletID, result.WalletID)
}
