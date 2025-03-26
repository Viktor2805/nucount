// Package services provides the application's service layer, which includes business logic.
package services

import (
	"golang/internal/repository"
	sequence "golang/internal/services/sequenceService"
)

// Services interface
type Services struct {
	SequenceService sequence.SequenceServiceI
}

// NewServices initializes and returns a new Services instance with all required components.
func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		SequenceService: sequence.NewSequenceService(repos.SequenceRepo), // Pass the repository
	}
}
