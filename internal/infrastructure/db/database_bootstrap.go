package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/fardannozami/fincore/configs"
	_ "github.com/lib/pq"
)

func EnsureDatabase(cfg *configs.Config) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed connect postgres default:", err)
	}

	defer db.Close()

	var exists bool
	query := "SELECT 1 FROM pg_database WHERE datname = $1"
	err = db.QueryRow(query, os.Getenv("DB_NAME")).Scan(&exists)

	if err != nil {
		// database belum ada → create
		_, err = db.Exec("CREATE DATABASE " + os.Getenv("DB_NAME"))
		if err != nil {
			log.Fatal("Failed create database:", err)
		}
		log.Println("✅ Database created")
	} else {
		log.Println("✅ Database already exists")
	}
}
