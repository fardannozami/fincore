package domain

type Transaction struct {
	ID     string `gorm:"primaryKey"`
	FromID string
	ToID   string
	Amount int64
	Status string
}
