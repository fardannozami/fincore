package usecase

import (
	"testing"

	"github.com/fardannozami/fincore/internal/domain"
	"github.com/fardannozami/fincore/internal/repository"
	"github.com/fardannozami/fincore/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) (*TransferUsecase, *gorm.DB) {
	// Use testutil to setup PostgreSQL test database
	db := testutil.SetupDB(t,
		&domain.Wallet{},
		&domain.Ledger{},
		&domain.Transaction{},
	)

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

	// Use unique IDs to avoid conflicts in persistent test DB
	idA := uuid.NewString()
	idB := uuid.NewString()
	trxID := uuid.NewString()

	// seed wallet
	db.Create(&domain.Wallet{ID: idA, UserID: "User-A", Balance: 1000})
	db.Create(&domain.Wallet{ID: idB, UserID: "User-B", Balance: 500})

	err := uc.Transfer(idA, idB, 300, trxID)
	assert.NoError(t, err)

	// cek saldo
	var a, b domain.Wallet
	db.First(&a, "id = ?", idA)
	db.First(&b, "id = ?", idB)

	assert.Equal(t, int64(700), a.Balance)
	assert.Equal(t, int64(800), b.Balance)

	// cek ledgers for these specific IDs
	var ledgers []domain.Ledger
	db.Find(&ledgers, "ref_id = ?", trxID)

	assert.Len(t, ledgers, 2)

	// cek transaction
	var trx domain.Transaction
	err = db.First(&trx, "id = ?", trxID).Error
	assert.NoError(t, err)
	assert.Equal(t, "SUCCESS", trx.Status)
}

func TestTransfer_InsufficientBalance(t *testing.T) {
	uc, db := setupTest(t)

	idA := uuid.NewString()
	idB := uuid.NewString()
	trxID := uuid.NewString()

	db.Create(&domain.Wallet{ID: idA, UserID: "User-A", Balance: 100})
	db.Create(&domain.Wallet{ID: idB, UserID: "User-B", Balance: 0})

	err := uc.Transfer(idA, idB, 200, trxID)
	assert.Error(t, err)

	// saldo tidak berubah
	var a, b domain.Wallet
	db.First(&a, "id = ?", idA)
	db.First(&b, "id = ?", idB)

	assert.Equal(t, int64(100), a.Balance)
	assert.Equal(t, int64(0), b.Balance)
}
