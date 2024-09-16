package controllers

import (
	"golang/internal/services"
)

type Controllers struct {
	DNAController DNAController
}

func NewControllers(services *services.Services) *Controllers {
	return &Controllers{
		DNAController: *NewDNAController(&services.DNAService),
	}
}
