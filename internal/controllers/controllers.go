package controllers

import (
	"golang/internal/services"
)

type Controllers struct {
	SequenceController *SequenceController
}

func NewControllers(services *services.Services) *Controllers {
	return &Controllers{
		SequenceController: NewSequenceController(services.SequenceService),
	}
}
