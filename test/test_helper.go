package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/fardannozami/fincore/internal/delivery/http/handler"
	"github.com/fardannozami/fincore/internal/domain"
	infraDB "github.com/fardannozami/fincore/internal/infrastructure/db"
	"github.com/fardannozami/fincore/internal/repository"
	"github.com/fardannozami/fincore/internal/usecase"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupTestServer() (*httptest.Server, *gorm.DB) {
	database := setupTestDB()

	// repo
	walletRepo := repository.NewWalletRepository(database)
	ledgerRepo := repository.NewLedgerRepository(database)
	transactionRepo := repository.NewTransactionRepository(database)

	// usecase
	transferUsecase := usecase.NewTransferUsecase(
		database,
		walletRepo,
		ledgerRepo,
		transactionRepo,
	)

	// handler
	walletHandler := handler.NewWalletHandler(database, walletRepo)
	transferHandler := handler.NewTransferHandler(transferUsecase)

	// router
	mux := http.NewServeMux()
	mux.HandleFunc("/wallet", walletHandler.CreateWallet)
	mux.HandleFunc("/transfer", transferHandler.Transfer)

	return httptest.NewServer(mux), database
}

func setupTestDB() *gorm.DB {
	user := "postgres"
	password := ""
	host := "localhost"
	port := "5432"
	dbname := "fincore_test"
	sslmode := "disable"

	infraDB.EnsureDBExists(host, port, user, password, dbname, sslmode)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)

	// Buat logger khusus untuk mengatur SlowThreshold
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Reset ke 200ms untuk verifikasi optimasi
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger, // Gunakan logger yang baru
	})
	if err != nil {
		log.Fatal("failed connect db:", err)
	}

	sqlDB, _ := db.DB()

	// 🔥 penting untuk concurrency test
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)

	// migrate
	err = db.AutoMigrate(
		&domain.Wallet{},
		&domain.Ledger{},
		&domain.Transaction{},
	)
	if err != nil {
		log.Fatal("failed migrate:", err)
	}

	return db
}

func topUp(dbConn *gorm.DB, walletID string, amount int64) {
	err := dbConn.Create(&domain.Ledger{
		ID:       uuid.NewString(), // Generasi ID agar tidak duplicate key
		WalletID: walletID,
		Amount:   amount,
		Type:     "CREDIT",
	}).Error

	if err != nil {
		log.Fatalf("failed topUp: %v", err)
	}

	err = dbConn.Model(&domain.Wallet{}).
		Where("id = ?", walletID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error

	if err != nil {
		log.Fatalf("failed update balance in topUp: %v", err)
	}
}

func createWallet(baseURL string) string {
	resp, err := http.Post(baseURL+"/wallet", "application/json", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	return res["id"].(string)
}

func getBalance(db *gorm.DB, walletID string) int64 {
	var wallet domain.Wallet

	err := db.First(&wallet, "id = ?", walletID).Error
	if err != nil {
		panic(err)
	}

	return wallet.Balance
}
