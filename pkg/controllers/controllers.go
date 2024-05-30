package controllers

import "golang/pkg/services"

type Controllers struct {
	TransactionController TransactionController
}

func NewControllers(services *services.Services) *Controllers {
	return &Controllers{
		TransactionController: *NewTransactionController(&services.TransactionService),
	}
}
