package repository

import (
	"context"

	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/config"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	// HealthCheckRepository is an interface that has all the function to be implemented inside health check repository
	HealthCheckRepository interface {
		Check() (bool, error)
	}

	// HealthCheckRepositoryImpl is an app health check struct that consists of all the dependencies needed for health check repository
	HealthCheckRepositoryImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		DB      *mongo.Client
		Redis   *redis.Client
	}
)

// NewHealthCheckRepository return new instances health check repository
func NewHealthCheckRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, db *mongo.Client, redis *redis.Client) *HealthCheckRepositoryImpl {
	return &HealthCheckRepositoryImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		DB:      db,
		Redis:   redis,
	}
}

func (hr *HealthCheckRepositoryImpl) Check() (bool, error) {
	if err := hr.DB.Ping(hr.Context, nil); err != nil {
		hr.Logger.Error("HealthCheckRepositoryImpl.Check() Ping DB ERROR, ", err)
		return false, nil
	}

	if err := hr.Redis.Ping(hr.Context).Err(); err != nil {
		hr.Logger.Error("HealthCheckRepositoryImpl.Check() Ping Redis ERROR, ", err)
		return false, nil
	}

	return true, nil
}
