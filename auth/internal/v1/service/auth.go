package service

import (
	"context"
	"time"

	"github.com/PickHD/singkatin-revamp/auth/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/auth/internal/v1/helper"
	"github.com/PickHD/singkatin-revamp/auth/internal/v1/model"
	"github.com/PickHD/singkatin-revamp/auth/internal/v1/repository"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// AuthService is an interface that has all the function to be implemented inside auth service
	AuthService interface {
		RegisterUser(ctx context.Context, req *model.RegisterRequest) (*model.RegisterResponse, error)
		LoginUser(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
	}

	// AuthServiceImpl is an app auth struct that consists of all the dependencies needed for auth service
	AuthServiceImpl struct {
		Context  context.Context
		Config   *config.Configuration
		Logger   *logrus.Logger
		Tracer   *trace.TracerProvider
		AuthRepo repository.AuthRepository
	}
)

// NewAuthService return new instances auth service
func NewAuthService(ctx context.Context, config *config.Configuration, logger *logrus.Logger, tracer *trace.TracerProvider, authRepo repository.AuthRepository) *AuthServiceImpl {
	return &AuthServiceImpl{
		Context:  ctx,
		Config:   config,
		Logger:   logger,
		Tracer:   tracer,
		AuthRepo: authRepo,
	}
}

func (as *AuthServiceImpl) RegisterUser(ctx context.Context, req *model.RegisterRequest) (*model.RegisterResponse, error) {
	tr := otel.GetTracerProvider().Tracer("Auth-RegisterUser service")
	_, span := tr.Start(ctx, "Start RegisterUser")
	defer span.End()

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

func (as *AuthServiceImpl) LoginUser(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	tr := otel.GetTracerProvider().Tracer("Auth-LoginUser service")
	_, span := tr.Start(ctx, "Start LoginUser")
	defer span.End()

	err := validateLoginUser(req)
	if err != nil {
		return nil, err
	}

	user, err := as.AuthRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// verify user password by comparing incoming request password with crypted password stored in database
	if !helper.CheckPasswordHash(user.Password, req.Password) {
		return nil, model.NewError(model.Validation, "invalid password")
	}

	// generate access token jwt
	token, expiredAt, err := as.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		AccessToken: token,
		Type:        "Bearer",
		ExpireAt:    time.Now().Add(expiredAt),
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

func validateLoginUser(req *model.LoginRequest) error {
	if !model.IsValidEmail.MatchString(req.Email) {
		return model.NewError(model.Validation, "invalid email")
	}

	return nil
}

func (as *AuthServiceImpl) generateJWT(user *model.User) (string, time.Duration, error) {
	var (
		payloadUserID   = "user_id"
		payloadFullName = "full_name"
		payloadEmail    = "email"
		payloadExpires  = "exp"
		JWTExpire       = time.Duration(as.Config.Common.JWTExpire) * time.Hour
	)

	claims := jwt.MapClaims{}
	claims[payloadUserID] = user.ID.Hex()
	claims[payloadFullName] = user.FullName
	claims[payloadEmail] = user.Email
	claims[payloadExpires] = time.Now().Add(JWTExpire).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(as.Config.Secret.JWTSecret))
	if err != nil {
		return "", 0, err
	}

	return signedToken, JWTExpire, nil
}
