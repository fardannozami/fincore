package db

import (
	"github.com/fardannozami/fincore/internal/domain"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&domain.Wallet{},
		&domain.Ledger{},
		&domain.Transaction{},
	)
}
