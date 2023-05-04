package service

import (
	"context"

	"github.com/PickHD/singkatin-revamp/auth/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/auth/internal/v1/helper"
	"github.com/PickHD/singkatin-revamp/auth/internal/v1/model"
	"github.com/PickHD/singkatin-revamp/auth/internal/v1/repository"
	"github.com/sirupsen/logrus"
)

type (
	// AuthService is an interface that has all the function to be implemented inside auth service
	AuthService interface {
		RegisterUser(ctx context.Context, req *model.RegisterRequest) (*model.RegisterResponse, error)
	}

	// AuthServiceImpl is an app auth struct that consists of all the dependencies needed for auth service
	AuthServiceImpl struct {
		Context  context.Context
		Config   *config.Configuration
		Logger   *logrus.Logger
		AuthRepo repository.AuthRepository
	}
)

// NewAuthService return new instances auth service
func NewAuthService(ctx context.Context, config *config.Configuration, logger *logrus.Logger, authRepo repository.AuthRepository) *AuthServiceImpl {
	return &AuthServiceImpl{
		Context:  ctx,
		Config:   config,
		Logger:   logger,
		AuthRepo: authRepo,
	}
}

func (as *AuthServiceImpl) RegisterUser(ctx context.Context, req *model.RegisterRequest) (*model.RegisterResponse, error) {
	err := validateRegisterUser(req)
	if err != nil {
		return nil, err
	}

	hashPass, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	data, err := as.AuthRepo.CreateUser(ctx, &model.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: hashPass,
	})
	if err != nil {
		return nil, err
	}

	return &model.RegisterResponse{
		ID:    data.ID.Hex(),
		Email: data.Email,
	}, nil
}

func validateRegisterUser(req *model.RegisterRequest) error {
	if len(req.FullName) < 3 {
		return model.NewError(model.Validation, "full name must more than 3")
	}

	if req.FullName == "" {
		return model.NewError(model.Validation, "full name required")
	}

	if !model.IsValidEmail.MatchString(req.Email) {
		return model.NewError(model.Validation, "invalid email")
	}

	return nil
}
