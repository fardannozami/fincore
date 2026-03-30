package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/fardannozami/fincore/internal/domain"
)

func setupTransactionTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&domain.Transaction{})
	assert.NoError(t, err)

	return db
}

func TestTransactionRepository_Create(t *testing.T) {
	db := setupTransactionTestDB(t)
	repo := NewTransactionRepository(db)

	tx := db.Begin()

	trx := &domain.Transaction{
		ID:     "trx-1",
		FromID: "wallet-1",
		ToID:   "wallet-2",
		Amount: 1000,
		Status: "SUCCESS",
	}

	err := repo.Create(tx, trx)
	assert.NoError(t, err)

	tx.Commit()

	// verify ke DB
	var result domain.Transaction
	err = db.First(&result, "id = ?", "trx-1").Error

	assert.NoError(t, err)
	assert.Equal(t, "wallet-1", result.FromID)
	assert.Equal(t, "wallet-2", result.ToID)
	assert.Equal(t, int64(1000), result.Amount)
	assert.Equal(t, "SUCCESS", result.Status)
}
