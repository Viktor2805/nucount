package controllers

import (
	"golang/internal/services"
)

type Controllers struct {
	BasesCounterController *BasesController
}

func NewControllers(services *services.Services) *Controllers {
	return &Controllers{
		BasesCounterController: NewBasesCounterController(services.BasesCounterService),
	}
}
