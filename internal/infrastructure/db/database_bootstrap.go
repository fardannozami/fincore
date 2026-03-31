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
	EnsureDBExists(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
}

func EnsureDBExists(host, port, user, password, dbname, sslmode string) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		host,
		port,
		user,
		password,
		sslmode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed connect postgres default:", err)
	}

	defer db.Close()

	var exists bool
	query := "SELECT 1 FROM pg_database WHERE datname = $1"
	err = db.QueryRow(query, dbname).Scan(&exists)

	if err != nil {
		// database belum ada → create
		_, err = db.Exec("CREATE DATABASE " + dbname)
		if err != nil {
			log.Fatal("Failed create database:", err)
		}
		log.Printf("✅ Database %s created\n", dbname)
	} else {
		log.Printf("✅ Database %s already exists\n", dbname)
	}
}
