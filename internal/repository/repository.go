package repository

import (
	sequence "golang/internal/repository/sequence"

	"gorm.io/gorm"
)

type Repositories struct {
	SequenceRepo sequence.SequenceRepositoryI
}

func Paginate(limit, offset int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if limit > 0 {
			db = db.Limit(limit)
		}
		if offset >= 0 {
			db = db.Offset(offset)
		}
		return db
	}
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		SequenceRepo: sequence.NewSequenceRepository(db),
	}
}
