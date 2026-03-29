package domain

type Ledger struct {
	ID       string `gorm:"primaryKey"`
	WalletID string
	Amount   int64
	Type     string
	RefID    string
}
