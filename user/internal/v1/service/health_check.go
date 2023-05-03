package service

import (
	"context"

	"github.com/PickHD/singkatin-revamp/user/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/user/internal/v1/repository"
)

type (
	// HealthCheckService is an interface that has all the function to be implemented inside health check service
	HealthCheckService interface {
		Check() (bool, error)
	}

	// HealthCheckServiceImpl is an app health check struct that consists of all the dependencies needed for health check service
	HealthCheckServiceImpl struct {
		Context         context.Context
		Config          *config.Configuration
		HealthCheckRepo repository.HealthCheckRepository
	}
)

// NewHealthCheckService return new instances health check service
func NewHealthCheckService(ctx context.Context, config *config.Configuration, healthCheckRepo repository.HealthCheckRepository) *HealthCheckServiceImpl {
	return &HealthCheckServiceImpl{
		Context:         ctx,
		Config:          config,
		HealthCheckRepo: healthCheckRepo,
	}
}

func (hs *HealthCheckServiceImpl) Check() (bool, error) {
	return hs.HealthCheckRepo.Check()
}
