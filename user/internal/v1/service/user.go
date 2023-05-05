package service

import (
	"context"

	"github.com/PickHD/singkatin-revamp/user/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/user/internal/v1/model"
	"github.com/PickHD/singkatin-revamp/user/internal/v1/repository"
	"github.com/sirupsen/logrus"
)

type (
	// UserService is an interface that has all the function to be implemented inside user service
	UserService interface {
		GetUserDetail(email string) (*model.User, error)
	}

	// UserServiceImpl is an app user struct that consists of all the dependencies needed for user service
	UserServiceImpl struct {
		Context  context.Context
		Config   *config.Configuration
		Logger   *logrus.Logger
		UserRepo repository.UserRepository
	}
)

// NewUserService return new instances user service
func NewUserService(ctx context.Context, config *config.Configuration, logger *logrus.Logger, userRepo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		Context:  ctx,
		Config:   config,
		Logger:   logger,
		UserRepo: userRepo,
	}
}

func (us *UserServiceImpl) GetUserDetail(email string) (*model.User, error) {
	return us.UserRepo.FindByEmail(us.Context, email)
}
