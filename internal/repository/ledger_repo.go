package repository

import (
	"github.com/fardannozami/fincore/internal/domain"
	"gorm.io/gorm"
)

type LedgerRepository struct {
	db *gorm.DB
}

func NewLedgerRepository(db *gorm.DB) *LedgerRepository {
	return &LedgerRepository{db: db}
}

func (r *LedgerRepository) Create(tx *gorm.DB, ledger *domain.Ledger) error {
	return tx.Create(ledger).Error
}
