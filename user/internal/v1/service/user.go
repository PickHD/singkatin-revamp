package service

import (
	"context"

	"github.com/PickHD/singkatin-revamp/user/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/user/internal/v1/model"
	"github.com/PickHD/singkatin-revamp/user/internal/v1/repository"
	shortenerpb "github.com/PickHD/singkatin-revamp/user/pkg/api/v1/proto/shortener"
	"github.com/sirupsen/logrus"
)

type (
	// UserService is an interface that has all the function to be implemented inside user service
	UserService interface {
		GetUserDetail(email string) (*model.User, error)
		GetUserShorts(userID string) ([]model.UserShorts, error)
	}

	// UserServiceImpl is an app user struct that consists of all the dependencies needed for user service
	UserServiceImpl struct {
		Context      context.Context
		Config       *config.Configuration
		Logger       *logrus.Logger
		UserRepo     repository.UserRepository
		ShortClients shortenerpb.ShortenerServiceClient
	}
)

// NewUserService return new instances user service
func NewUserService(ctx context.Context, config *config.Configuration, logger *logrus.Logger, userRepo repository.UserRepository, shortClients shortenerpb.ShortenerServiceClient) *UserServiceImpl {
	return &UserServiceImpl{
		Context:      ctx,
		Config:       config,
		Logger:       logger,
		UserRepo:     userRepo,
		ShortClients: shortClients,
	}
}

func (us *UserServiceImpl) GetUserDetail(email string) (*model.User, error) {
	return us.UserRepo.FindByEmail(us.Context, email)
}

func (us *UserServiceImpl) GetUserShorts(userID string) ([]model.UserShorts, error) {
	data, err := us.ShortClients.GetListShortenerByUserID(us.Context, &shortenerpb.ListShortenerRequest{
		UserId: userID})
	if err != nil {
		us.Logger.Error("UserServiceImpl.GetUserShorts ShortClients ERROR, ", err)
		return nil, err
	}

	if len(data.Shorteners) < 1 {
		return nil, nil
	}

	shorteners := make([]model.UserShorts, len(data.Shorteners))

	for i, q := range data.Shorteners {
		shorteners[i] = model.UserShorts{
			ID:       q.GetId(),
			FullURL:  q.GetFullUrl(),
			ShortURL: q.GetShortUrl(),
			Visited:  q.GetVisited(),
		}
	}

	return shorteners, nil
}
