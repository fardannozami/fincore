package dto

type TransferRequest struct {
	FromID string `json:"from_id"`
	ToID   string `json:"to_id"`
	Amount int64  `json:"amount"`
}

type TransferResponse struct {
	Message string `json:"message"`
}
