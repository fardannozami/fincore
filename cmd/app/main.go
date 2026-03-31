package main

import (
	"log"
	"net/http"

	"github.com/fardannozami/fincore/configs"
	httpDelivery "github.com/fardannozami/fincore/internal/delivery/http"
	"github.com/fardannozami/fincore/internal/delivery/http/handler"
	"github.com/fardannozami/fincore/internal/infrastructure/db"
	"github.com/fardannozami/fincore/internal/repository"
	"github.com/fardannozami/fincore/internal/usecase"
)

func main() {
	cfg := configs.LoadConfig()

	// 🔥 ensure db exists
	db.EnsureDatabase(cfg)

	// 🔌 connect db
	database := db.NewPostgres(cfg.DBUrl)

	// 🧱 auto migrate
	if cfg.AutoMigrate {
		db.AutoMigrate(database)
	}

	// ==============================
	// 🧱 INIT REPOSITORY
	// ==============================
	walletRepo := repository.NewWalletRepository(database)
	ledgerRepo := repository.NewLedgerRepository(database)
	transactionRepo := repository.NewTransactionRepository(database)

	// ==============================
	// 🧠 INIT USECASE
	// ==============================
	transferUsecase := usecase.NewTransferUsecase(
		database,
		walletRepo,
		ledgerRepo,
		transactionRepo,
	)

	// ==============================
	// 🌐 INIT HANDLER
	// ==============================
	walletHandler := handler.NewWalletHandler(database, walletRepo)
	transferHandler := handler.NewTransferHandler(transferUsecase)

	// ==============================
	// 🔌 REGISTER ROUTES
	// ==============================
	httpDelivery.RegisterRoutes(walletHandler, transferHandler)

	// ==============================
	// 🚀 START SERVER
	// ==============================
	port := cfg.Port

	addr := ":" + port

	log.Printf("🚀 Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
