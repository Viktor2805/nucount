package repository

import (
	"golang/internal/models"

	"gorm.io/gorm"
)

type AssemblyStatsRepositoryI interface {
	Create() (models.AssemblyStats, error)
}

type AssemblyStatsRepository struct {
	db *gorm.DB
}

func NewAssemblyStatsRepository(db *gorm.DB) *AssemblyStatsRepository {
	return &AssemblyStatsRepository{
		db: db,
	}
}

func (r *AssemblyStatsRepository) CreateTransaction(transaction *models.AssemblyStats) error {
	return r.db.Create(transaction).Error
}
