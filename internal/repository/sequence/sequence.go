package repository

import (
	"golang/internal/models"

	"gorm.io/gorm"
)

type SequenceRepositoryI interface {
	Create(sequence *models.Sequence) error
}

type SequenceRepository struct {
	db *gorm.DB
}

func NewSequenceRepository(db *gorm.DB) *SequenceRepository {
	return &SequenceRepository{
		db: db,
	}
}

func (r *SequenceRepository) Create(sequence *models.Sequence) error {
	return r.db.Create(sequence).Error
}
