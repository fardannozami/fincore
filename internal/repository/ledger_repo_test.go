package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/fardannozami/fincore/internal/domain"
)

func setupLedgerTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&domain.Ledger{})
	assert.NoError(t, err)

	return db
}

func TestLedgerRepository_Create(t *testing.T) {
	db := setupLedgerTestDB(t)
	repo := NewLedgerRepository(db)

	tx := db.Begin()

	ledger := &domain.Ledger{
		ID:       "ledger-1",
		WalletID: "wallet-1",
		Amount:   1000,
		Type:     "CREDIT",
		RefID:    "trx-1",
	}

	err := repo.Create(tx, ledger)
	assert.NoError(t, err)

	tx.Commit()

	// verify ke DB
	var result domain.Ledger
	err = db.First(&result, "id = ?", "ledger-1").Error

	assert.NoError(t, err)
	assert.Equal(t, int64(1000), result.Amount)
	assert.Equal(t, "CREDIT", result.Type)
	assert.Equal(t, "wallet-1", result.WalletID)
}
