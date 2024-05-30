package repository

import (
	"fmt"
	"golang/pkg/models"
	"strings"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransactionInBatches(transactions []*models.Transaction, batchSize int) error
	GetTransactions(filters map[string]interface{}, limit, offset int) ([]models.Transaction, error)
}

type PostgresTransactionRepository struct {
	db *gorm.DB
}

func NewPostgresTransactionRepository(db *gorm.DB) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{db: db}
}

func (r *PostgresTransactionRepository) CreateTransaction(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *PostgresTransactionRepository) CreateTransactionInBatches(transactions []*models.Transaction, batchSize int) error {
	return r.db.CreateInBatches(transactions, batchSize).Error
}

func (r *PostgresTransactionRepository) GetTransactions(filters map[string]interface{}, limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := r.db.Model(&models.Transaction{})

	for key, val := range filters {
		switch key {
		case "transaction_id":
			query = query.Where("transaction_id = ?", val)
		case "terminal_id":
			ids, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("terminal_id must be a string")
			}
			terminalIDs := strings.Split(ids, ",")
			query = query.Where("terminal_id IN (?)", terminalIDs)
		case "status":
			query = query.Where("status = ?", val)
		case "payment_type":
			query = query.Where("payment_type = ?", val)
		case "date_from":
			query = query.Where("date_post >= ?", val)
		case "date_to":
			query = query.Where("date_post <= ?", val)
		case "payment_narrative":
			query = query.Where("payment_narrative ILIKE ?", "%"+val.(string)+"%")
		}
	}

	err := query.Scopes(Paginate(limit, offset)).Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
