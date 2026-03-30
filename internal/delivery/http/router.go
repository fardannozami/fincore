package http

import (
	"net/http"

	"github.com/fardannozami/fincore/internal/delivery/http/handler"
)

func RegisterRoutes(
	walletHandler *handler.WalletHandler,
	transferHandler *handler.TransferHandler,
) {
	http.HandleFunc("/wallet", walletHandler.CreateWallet)
	http.HandleFunc("/transfer", transferHandler.Transfer)
}
