// Package repository initializes
package repository

import (
	dna "golang/internal/repository/dnaRepository"

	"gorm.io/gorm"
)

type Repositories struct {
	DNAAssemblyStatsRepo dna.AssemblyStatsRepository
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
		DNAAssemblyStatsRepo: *dna.NewAssemblyStatsRepository(db),
	}
}
