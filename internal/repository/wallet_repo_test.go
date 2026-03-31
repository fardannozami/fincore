package repository

import (
	"testing"

	"github.com/fardannozami/fincore/internal/domain"
	"github.com/fardannozami/fincore/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWalletRepository_Create(t *testing.T) {
	db := testutil.SetupDB(t, &domain.Wallet{})
	repo := NewWalletRepository(db)

	tx := db.Begin()
	id := uuid.NewString()

	wallet := &domain.Wallet{
		ID:      id,
		UserID:  "user-1",
		Balance: 1000,
	}

	err := repo.Create(tx, wallet)
	assert.NoError(t, err)

	tx.Commit()

	// verify
	var result domain.Wallet
	err = db.First(&result, "id = ?", id).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(1000), result.Balance)
}

func TestWalletRepository_FindByIDForUpdate(t *testing.T) {
	db := testutil.SetupDB(t, &domain.Wallet{})
	repo := NewWalletRepository(db)

	id := uuid.NewString()
	// seed data
	db.Create(&domain.Wallet{
		ID:      id,
		UserID:  "user-1",
		Balance: 500,
	})

	tx := db.Begin()

	wallet, err := repo.FindByIDForUpdate(tx, id)

	assert.NoError(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, int64(500), wallet.Balance)

	tx.Commit()
}

func TestWalletRepository_UpdateBalanceAtomic(t *testing.T) {
	db := testutil.SetupDB(t, &domain.Wallet{})
	repo := NewWalletRepository(db)

	walletID := uuid.NewString()
	db.Create(&domain.Wallet{
		ID:      walletID,
		UserID:  "user-1",
		Balance: 1000,
	})

	t.Run("Success Increment", func(t *testing.T) {
		tx := db.Begin()
		w, err := repo.UpdateBalanceAtomic(tx, walletID, 500)
		assert.NoError(t, err)
		if assert.NotNil(t, w) {
			assert.Equal(t, int64(1500), w.Balance)
		}
		tx.Commit()
	})

	t.Run("Success Decrement", func(t *testing.T) {
		tx := db.Begin()
		w, err := repo.UpdateBalanceAtomic(tx, walletID, -200)
		assert.NoError(t, err)
		if assert.NotNil(t, w) {
			assert.Equal(t, int64(1300), w.Balance)
		}
		tx.Commit()
	})

	t.Run("Fail Decrement Insufficient Balance", func(t *testing.T) {
		tx := db.Begin()
		w, err := repo.UpdateBalanceAtomic(tx, walletID, -2000)
		assert.Error(t, err)
		assert.Nil(t, w)
		assert.Equal(t, ErrInsufficientBalance, err)
		tx.Commit()
	})
}
