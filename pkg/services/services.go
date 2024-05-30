package services

import (
	"golang/pkg/repository"
)

type Services struct {
	TransactionService TransactionService
	CsvService         CSVService
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		TransactionService: *NewTransactionService(repos.TransactionRepo, NewCSVService()),
		CsvService:         *NewCSVService(),
	}
}
