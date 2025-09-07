package db

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database wraps the GORM DB connection and provides methods to interact with the database.
type Db struct {
	*gorm.DB
	logger *zap.Logger
}

// NewDatabase initializes a new Database instance.
func New(logger *zap.Logger) (*Db, error) {
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
		logger.Error("Failed to open", zap.Error(err))
	}

	if err := pingDatabase(db, logger); err != nil {
		return nil, err
	}

	logger.Info("Database connection established successfully")
	return &Db{DB: db, logger: logger}, err
}

// pingDatabase ensures the database connection is alive by pinging it.
func pingDatabase(gormDB *gorm.DB, logger *zap.Logger) error {
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database object: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database ping successful")
	return nil
}

// Close disconnects the database connection.
func (d *Db) Close() error {
	if d == nil || d.DB == nil {
		return nil
	}

	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB from GORM: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	d.logger.Info("Database connection closed successfully.")
	return nil
}
