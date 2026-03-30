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

	return u.db.Transaction(func(tx *gorm.DB) error {

		// 🔐 Hindari deadlock
		ids := []string{fromID, toID}
		sort.Strings(ids)

		w1, err := u.walletRepo.FindByIDForUpdate(tx, ids[0])
		if err != nil {
			return err
		}

		w2, err := u.walletRepo.FindByIDForUpdate(tx, ids[1])
		if err != nil {
			return err
		}

		var fromWallet, toWallet *domain.Wallet
		if w1.ID == fromID {
			fromWallet = w1
			toWallet = w2
		} else {
			fromWallet = w2
			toWallet = w1
		}

		// 💥 cek saldo
		if fromWallet.Balance < amount {
			return errors.New("insufficient balance")
		}

		// 🧾 ledger debit
		if err := u.ledgerRepo.Create(tx, &domain.Ledger{
			ID:       uuid.NewString(),
			WalletID: fromWallet.ID,
			Amount:   -amount,
			Type:     "DEBIT",
			RefID:    refID,
		}); err != nil {
			return err
		}

		// 🧾 ledger credit
		if err := u.ledgerRepo.Create(tx, &domain.Ledger{
			ID:       uuid.NewString(),
			WalletID: toWallet.ID,
			Amount:   amount,
			Type:     "CREDIT",
			RefID:    refID,
		}); err != nil {
			return err
		}

		// update balance
		fromWallet.Balance -= amount
		toWallet.Balance += amount

		if err := u.walletRepo.Update(tx, fromWallet); err != nil {
			return err
		}

		if err := u.walletRepo.Update(tx, toWallet); err != nil {
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
