package main

import (
	"github.com/fardannozami/fincore/configs"
	"github.com/fardannozami/fincore/internal/infrastructure/db"
)

func main() {
	cfg := configs.LoadConfig()

	// 🔥 Auto create DB
	db.EnsureDatabase(cfg)

	// connect db
	database := db.NewPostgres(cfg.DBUrl)

	if cfg.AutoMigrate {
		db.AutoMigrate(database)
	}

	println("🚀 App running...")
}
