package repository

import (
	"github.com/fardannozami/fincore/internal/domain"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *gorm.DB, trx *domain.Transaction) error {
	return tx.Create(trx).Error
}
