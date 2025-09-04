package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database wraps the GORM DB connection and provides methods to interact with the database.
type Db struct {
	db *gorm.DB
}

// NewDatabase initializes a new Database instance.
func New() (*Db, error) {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: databaseURL,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	if err := pingDatabase(db); err != nil {
		return nil, err
	}
	
	return &Db{db: db}, err
}

// pingDatabase ensures the database connection is alive by pinging it.
func pingDatabase(gormDB *gorm.DB) error {
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database object: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection is alive.")
	return nil
}

// GetDb returns the db instance
func (d *Db) GetDB() *gorm.DB {
	return d.db
}

// Close disconnects the database connection.
func (d *Db) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB from GORM: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	log.Println("Database connection closed successfully.")
	return nil
}
