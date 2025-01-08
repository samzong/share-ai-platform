package database

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/samzong/share-ai-platform/internal/models"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		viper.GetString("database.postgres.host"),
		viper.GetString("database.postgres.port"),
		viper.GetString("database.postgres.user"),
		viper.GetString("database.postgres.password"),
		viper.GetString("database.postgres.dbname"),
		viper.GetString("database.postgres.sslmode"),
	)

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(viper.GetInt("database.postgres.max_idle_conns"))
	sqlDB.SetMaxOpenConns(viper.GetInt("database.postgres.max_open_conns"))
	sqlDB.SetConnMaxLifetime(viper.GetDuration("database.postgres.conn_max_lifetime") * time.Hour)

	// Auto migrate the schema
	err = autoMigrate(db)
	if err != nil {
		return fmt.Errorf("failed to auto migrate schema: %v", err)
	}

	DB = db
	log.Println("Database connection established")
	return nil
}

// autoMigrate automatically migrates the schema
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Image{},
		&models.Tag{},
		&models.Provider{},
		&models.ImageProvider{},
		&models.Collection{},
	)
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
} 