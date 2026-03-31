package testutil

import (
	"fmt"
	"testing"

	infraDB "github.com/fardannozami/fincore/internal/infrastructure/db"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupDB creates a test database and migrates the provided models.
// It uses PostgreSQL to support advanced features like atomic updates with RETURNING.
func SetupDB(t *testing.T, models ...interface{}) *gorm.DB {
	user := "postgres"
	password := ""
	host := "localhost"
	port := "5432"
	dbname := "fincore_test"
	sslmode := "disable"

	// Ensure the test database exists
	infraDB.EnsureDBExists(host, port, user, password, dbname, sslmode)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Migrate provided models
	if len(models) > 0 {
		err = db.AutoMigrate(models...)
		assert.NoError(t, err)
	}

	return db
}
