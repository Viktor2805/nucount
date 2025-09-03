// Package services provides the application's service layer, which includes business logic.
package services

import (
	"golang/internal/repository"
	basesCounter "golang/internal/services/dnaBaseCounter"
)

// Services interface
type Services struct {
	BasesCounterService basesCounter.BasesCounterServiceI
}

// NewServices initializes and returns a new Services instance with all required components.
func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		BasesCounterService: basesCounter.NewBasesCounter(),
	}
}
