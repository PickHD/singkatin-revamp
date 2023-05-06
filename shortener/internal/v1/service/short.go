package service

import (
	"context"

	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/model"
	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/repository"
	"github.com/sirupsen/logrus"
)

type (
	// ShortService is an interface that has all the function to be implemented inside short service
	ShortService interface {
		GetListShortenerByUserID(ctx context.Context, userID string) ([]model.Short, error)
	}

	// ShortServiceImpl is an app short struct that consists of all the dependencies needed for short repository
	ShortServiceImpl struct {
		Context   context.Context
		Config    *config.Configuration
		Logger    *logrus.Logger
		ShortRepo repository.ShortRepository
	}
)

// NewShortService return new instances short service
func NewShortService(ctx context.Context, config *config.Configuration, logger *logrus.Logger, shortRepo repository.ShortRepository) *ShortServiceImpl {
	return &ShortServiceImpl{
		Context:   ctx,
		Config:    config,
		Logger:    logger,
		ShortRepo: shortRepo,
	}
}

func (ss *ShortServiceImpl) GetListShortenerByUserID(ctx context.Context, userID string) ([]model.Short, error) {
	data, err := ss.ShortRepo.GetListShortenerByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return data, nil
}
