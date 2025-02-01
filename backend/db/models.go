package db

import (
	"crypto-wallet/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		return nil, fmt.Errorf("config loading failed: %w", err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB_HOST,
		cfg.DB_PORT,
		cfg.DB_USER,
		cfg.DB_PASS,
		cfg.DB_NAME,
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	if err := db.AutoMigrate(&Transaction{}); err != nil {
		log.Printf("Failed to migrate database: %v", err)
		return db, fmt.Errorf("database migration failed: %w", err)
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		log.Printf("Failed to migrate database: %v", err)
		return db, fmt.Errorf("database migration failed: %w", err)
	}

	SynchronizeTransactions(db)

	return db, nil
}
