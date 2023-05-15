package service

import (
	"context"

	"github.com/PickHD/singkatin-revamp/user/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/user/internal/v1/helper"
	"github.com/PickHD/singkatin-revamp/user/internal/v1/model"
	"github.com/PickHD/singkatin-revamp/user/internal/v1/repository"
	shortenerpb "github.com/PickHD/singkatin-revamp/user/pkg/api/v1/proto/shortener"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// UserService is an interface that has all the function to be implemented inside user service
	UserService interface {
		GetUserDetail(email string) (*model.User, error)
		GetUserShorts(userID string) ([]model.UserShorts, error)
		GenerateUserShorts(userID string, req *model.GenerateShortUserRequest) (*model.GenerateShortUserResponse, error)
		UpdateUserProfile(userID string, req *model.EditProfileRequest) error
	}

	// UserServiceImpl is an app user struct that consists of all the dependencies needed for user service
	UserServiceImpl struct {
		Context      context.Context
		Config       *config.Configuration
		Logger       *logrus.Logger
		Tracer       *trace.TracerProvider
		UserRepo     repository.UserRepository
		ShortClients shortenerpb.ShortenerServiceClient
	}
)

// NewUserService return new instances user service
func NewUserService(ctx context.Context, config *config.Configuration, logger *logrus.Logger, tracer *trace.TracerProvider, userRepo repository.UserRepository, shortClients shortenerpb.ShortenerServiceClient) *UserServiceImpl {
	return &UserServiceImpl{
		Context:      ctx,
		Config:       config,
		Logger:       logger,
		Tracer:       tracer,
		UserRepo:     userRepo,
		ShortClients: shortClients,
	}
}

func (us *UserServiceImpl) GetUserDetail(email string) (*model.User, error) {
	tr := us.Tracer.Tracer("User-GetUserDetail Service")
	_, span := tr.Start(us.Context, "Start GetUserDetail")
	defer span.End()

	return us.UserRepo.FindByEmail(us.Context, email)
}

func (us *UserServiceImpl) GetUserShorts(userID string) ([]model.UserShorts, error) {
	tr := us.Tracer.Tracer("User-GetUserShorts Service")
	_, span := tr.Start(us.Context, "Start GetUserShorts")
	defer span.End()

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

func (us *UserServiceImpl) GenerateUserShorts(userID string, req *model.GenerateShortUserRequest) (*model.GenerateShortUserResponse, error) {
	tr := us.Tracer.Tracer("User-GenerateUserShorts Service")
	_, span := tr.Start(us.Context, "Start GenerateUserShorts")
	defer span.End()

	msg := model.GenerateShortUserMessage{
		FullURL:  req.FullURL,
		UserID:   userID,
		ShortURL: helper.RandomStringBytesMaskImprSrcSB(8),
	}

	err := us.UserRepo.PublishCreateUserShortener(us.Context, &msg)
	if err != nil {
		return nil, err
	}

	return &model.GenerateShortUserResponse{
		ShortURL: msg.ShortURL,
	}, nil
}

func (us *UserServiceImpl) UpdateUserProfile(userID string, req *model.EditProfileRequest) error {
	tr := us.Tracer.Tracer("User-UpdateUserProfile Service")
	_, span := tr.Start(us.Context, "Start UpdateUserProfile")
	defer span.End()

	if req.FullName == "" {
		return model.NewError(model.Validation, "Full Name Required")
	}

	if len(req.FullName) < 3 {
		return model.NewError(model.Validation, "Full Name must more than 3")
	}

	err := us.UserRepo.UpdateProfileByID(us.Context, userID, req)
	if err != nil {
		return err
	}

	return nil
}
