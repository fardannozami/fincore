package repository

import (
	"errors"

	"github.com/fardannozami/fincore/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository struct {
	db *gorm.DB
}

var ErrInsufficientBalance = errors.New("insufficient balance")

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) FindByIDForUpdate(tx *gorm.DB, id string) (*domain.Wallet, error) {
	var wallet domain.Wallet

	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&wallet, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *WalletRepository) Create(tx *gorm.DB, wallet *domain.Wallet) error {
	return tx.Create(wallet).Error
}

func (r *WalletRepository) UpdateBalanceAtomic(
	tx *gorm.DB,
	id string,
	amount int64,
) (*domain.Wallet, error) {

	var wallet domain.Wallet

	result := tx.Raw(`
		UPDATE wallets
		SET balance = balance + ?
		WHERE id = ? AND (balance + ? >= 0)
		RETURNING id, balance
	`, amount, id, amount).Scan(&wallet)

	// 🔴 error dari DB
	if result.Error != nil {
		return nil, result.Error
	}

	// 🔴 tidak ada row ter-update = saldo tidak cukup / wallet tidak ada
	if result.RowsAffected == 0 {
		return nil, ErrInsufficientBalance
	}

	return &wallet, nil
}
