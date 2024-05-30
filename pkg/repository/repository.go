package repository

import "gorm.io/gorm"

type Repositories struct {
	TransactionRepo TransactionRepository
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
		TransactionRepo: NewPostgresTransactionRepository(db),
	}
}
