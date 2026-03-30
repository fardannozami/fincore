package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/fardannozami/fincore/internal/domain"
	"github.com/fardannozami/fincore/internal/repository"
)

func setupTest(t *testing.T) (*TransferUsecase, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// migrate semua tabel
	err = db.AutoMigrate(
		&domain.Wallet{},
		&domain.Ledger{},
		&domain.Transaction{},
	)
	assert.NoError(t, err)

	// repo
	walletRepo := repository.NewWalletRepository(db)
	ledgerRepo := repository.NewLedgerRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// usecase
	uc := NewTransferUsecase(db, walletRepo, ledgerRepo, transactionRepo)

	return uc, db
}

func TestTransfer_Success(t *testing.T) {
	uc, db := setupTest(t)

	// seed wallet
	db.Create(&domain.Wallet{ID: "A", Balance: 1000})
	db.Create(&domain.Wallet{ID: "B", Balance: 500})

	err := uc.Transfer("A", "B", 300, "trx-1")
	assert.NoError(t, err)

	// cek saldo
	var a, b domain.Wallet
	db.First(&a, "id = ?", "A")
	db.First(&b, "id = ?", "B")

	assert.Equal(t, int64(700), a.Balance)
	assert.Equal(t, int64(800), b.Balance)

	// cek ledger
	var ledgers []domain.Ledger
	db.Find(&ledgers)

	assert.Len(t, ledgers, 2)

	// cek transaction
	var trx domain.Transaction
	err = db.First(&trx, "id = ?", "trx-1").Error
	assert.NoError(t, err)
	assert.Equal(t, "SUCCESS", trx.Status)
}

func TestTransfer_InsufficientBalance(t *testing.T) {
	uc, db := setupTest(t)

	db.Create(&domain.Wallet{ID: "A", Balance: 100})
	db.Create(&domain.Wallet{ID: "B", Balance: 0})

	err := uc.Transfer("A", "B", 200, "trx-2")
	assert.Error(t, err)

	// saldo tidak berubah
	var a, b domain.Wallet
	db.First(&a, "id = ?", "A")
	db.First(&b, "id = ?", "B")

	assert.Equal(t, int64(100), a.Balance)
	assert.Equal(t, int64(0), b.Balance)
}
