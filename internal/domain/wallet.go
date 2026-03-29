package domain

type Wallet struct {
	ID      string `gorm:"primaryKey"`
	UserID  string
	Balance int64
}
