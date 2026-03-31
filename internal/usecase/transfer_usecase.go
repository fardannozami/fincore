package usecase

import (
	"errors"
	"sort"

	"github.com/fardannozami/fincore/internal/domain"
	"github.com/fardannozami/fincore/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransferUsecase struct {
	db              *gorm.DB
	walletRepo      *repository.WalletRepository
	ledgerRepo      *repository.LedgerRepository
	transactionRepo *repository.TransactionRepository
}

func NewTransferUsecase(
	db *gorm.DB,
	w *repository.WalletRepository,
	l *repository.LedgerRepository,
	t *repository.TransactionRepository,
) *TransferUsecase {
	return &TransferUsecase{db, w, l, t}
}

func (u *TransferUsecase) Transfer(fromID, toID string, amount int64, refID string) error {

	if fromID == toID {
		return errors.New("cannot transfer to same wallet")
	}

	if amount <= 0 {
		return errors.New("invalid amount")
	}

	// 🛠️ Optimasi: Generasi UUID di luar transaksi untuk menghemat waktu lock
	debitID := uuid.NewString()
	creditID := uuid.NewString()

	return u.db.Transaction(func(tx *gorm.DB) error {
		// 🔐 Anti-Deadlock: Selalu urutkan dari ID terkecil
		ids := []string{fromID, toID}
		sort.Strings(ids)

		for _, id := range ids {
			updateAmount := amount
			if id == fromID {
				updateAmount = -amount
			}

			// Atomic update: Cek saldo + Update saldo dalam 1 perintah SQL
			wallet, err := u.walletRepo.UpdateBalanceAtomic(tx, id, updateAmount)
			if err != nil {
				return err
			}
			if wallet.ID == "" {
				if id == fromID {
					return errors.New("insufficient balance or wallet not found")
				}
				return errors.New("wallet not found")
			}
		}

		// 🧾 ledger debit & credit
		if err := u.ledgerRepo.Create(tx, &domain.Ledger{
			ID:       debitID,
			WalletID: fromID,
			Amount:   -amount,
			Type:     "DEBIT",
			RefID:    refID,
		}); err != nil {
			return err
		}

		if err := u.ledgerRepo.Create(tx, &domain.Ledger{
			ID:       creditID,
			WalletID: toID,
			Amount:   amount,
			Type:     "CREDIT",
			RefID:    refID,
		}); err != nil {
			return err
		}

		// save transaction
		return u.transactionRepo.Create(tx, &domain.Transaction{
			ID:     refID,
			FromID: fromID,
			ToID:   toID,
			Amount: amount,
			Status: "SUCCESS",
		})
	})
}
