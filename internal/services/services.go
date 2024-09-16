// Package services provides the application's service layer, which includes business logic.
package services

import (
	"golang/internal/repository"
	dna "golang/internal/services/dnaService"
)

// Services interface
type Services struct {
	DNAService dna.DNAService
}

// NewServices initializes and returns a new Services instance with all required components.
func NewServices(_ *repository.Repositories) *Services {
	return &Services{
		DNAService: *dna.NewDNAService(),
	}
}
