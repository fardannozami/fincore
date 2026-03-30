package repository

import (
	"testing"

	"github.com/fardannozami/fincore/internal/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&domain.Wallet{})
	assert.NoError(t, err)

	return db
}

func TestWalletRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewWalletRepository(db)

	tx := db.Begin()

	wallet := &domain.Wallet{
		ID:      "wallet-1",
		UserID:  "user-1",
		Balance: 1000,
	}

	err := repo.Create(tx, wallet)
	assert.NoError(t, err)

	tx.Commit()

	// verify
	var result domain.Wallet
	err = db.First(&result, "id = ?", "wallet-1").Error
	assert.NoError(t, err)
	assert.Equal(t, int64(1000), result.Balance)
}

func TestWalletRepository_FindByIDForUpdate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewWalletRepository(db)

	// seed data
	db.Create(&domain.Wallet{
		ID:      "wallet-1",
		UserID:  "user-1",
		Balance: 500,
	})

	tx := db.Begin()

	wallet, err := repo.FindByIDForUpdate(tx, "wallet-1")

	assert.NoError(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, int64(500), wallet.Balance)

	tx.Commit()
}

func TestWalletRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewWalletRepository(db)

	// seed
	db.Create(&domain.Wallet{
		ID:      "wallet-1",
		UserID:  "user-1",
		Balance: 100,
	})

	tx := db.Begin()

	wallet, _ := repo.FindByIDForUpdate(tx, "wallet-1")
	wallet.Balance += 200

	err := repo.Update(tx, wallet)
	assert.NoError(t, err)

	tx.Commit()

	// verify
	var result domain.Wallet
	db.First(&result, "id = ?", "wallet-1")

	assert.Equal(t, int64(300), result.Balance)
}
