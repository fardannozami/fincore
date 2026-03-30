package repository

import (
	"github.com/fardannozami/fincore/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository struct {
	db *gorm.DB
}

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

func (r *WalletRepository) Update(tx *gorm.DB, wallet *domain.Wallet) error {
	return tx.Save(wallet).Error
}
