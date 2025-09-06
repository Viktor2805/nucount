// Package services provides the application's service layer, which includes business logic.
package services

import (
	"golang/internal/repository"
	nucleotide "golang/internal/services/nucleotide"
)

// Services interface
type Services struct {
	NucleotideService nucleotide.Service
}

// NewServices initializes and returns a new Services instance with all required components.
func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		NucleotideService: nucleotide.NewCounter(),
	}
}
