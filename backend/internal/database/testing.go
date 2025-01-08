package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB initializes a test database connection
func SetupTestDB() *gorm.DB {
	// Use test database configuration
	dsn := "host=localhost user=postgres password=postgres dbname=share_ai_platform_test port=5432 sslmode=disable"
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Silent,
			},
		),
	})
	
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Set the global DB instance
	DB = db

	return db
}

// CleanupTestDB cleans up the test database
func CleanupTestDB(db *gorm.DB) {
	// Get the underlying SQL database
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error getting underlying SQL DB: %v", err)
		return
	}

	// Close the database connection
	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
} 