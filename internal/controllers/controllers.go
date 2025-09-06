package controllers

import (
	"golang/internal/services"
)

type Controllers struct {
	NucleotideController *NucleotideController
}

func NewControllers(services *services.Services) *Controllers {
	return &Controllers{
		NucleotideController: NewNucleotideController(services.NucleotideService),
	}
}
