package controller

import (
	"context"
	"net/http"
	"strings"

	"github.com/PickHD/singkatin-revamp/auth/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/auth/internal/v1/helper"
	"github.com/PickHD/singkatin-revamp/auth/internal/v1/model"
	"github.com/PickHD/singkatin-revamp/auth/internal/v1/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// Authcontroller is an interface that has all the function to be implemented inside auth controller
	AuthController interface {
		Register(ctx *gin.Context)
		Login(ctx *gin.Context)
	}

	// AuthcontrollerImpl is an app auth struct that consists of all the dependencies needed for auth controller
	AuthControllerImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		Tracer  *trace.TracerProvider
		AuthSvc service.AuthService
	}
)

// NewAuthController return new instances auth controller
func NewAuthController(ctx context.Context, config *config.Configuration, logger *logrus.Logger, tracer *trace.TracerProvider, authSvc service.AuthService) *AuthControllerImpl {
	return &AuthControllerImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		Tracer:  tracer,
		AuthSvc: authSvc,
	}
}

// Check godoc
// @Summary      Register users
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user body model.RegisterRequest true "register user"
// @Success      201  {object}  helper.BaseResponse
// @Failure      400  {object}  helper.BaseResponse
// @Failure      500  {object}  helper.BaseResponse
// @Router       /register [post]
func (ac *AuthControllerImpl) Register(ctx *gin.Context) {
	var req model.RegisterRequest

	tr := otel.GetTracerProvider().Tracer("Auth-Register Controller")
	_, span := tr.Start(ctx, "Start Register")
	defer span.End()

	if err := ctx.BindJSON(&req); err != nil {
		helper.NewResponses[any](ctx, http.StatusBadRequest, "Invalid request", req, err, nil)
		return
	}

	data, err := ac.AuthSvc.RegisterUser(ctx, &req)
	if err != nil {
		if strings.Contains(err.Error(), string(model.Validation)) {
			helper.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), req, err, nil)
			return
		}

		helper.NewResponses[any](ctx, http.StatusInternalServerError, "Failed register user", data, err, nil)
		return
	}

	helper.NewResponses[any](ctx, http.StatusCreated, "Success register user", data, nil, nil)
}

// Check godoc
// @Summary      Login users
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user body model.LoginRequest true "login user"
// @Success      200  {object}  helper.BaseResponse
// @Failure      400  {object}  helper.BaseResponse
// @Failure      404  {object}  helper.BaseResponse
// @Failure      500  {object}  helper.BaseResponse
// @Router       /login [post]
func (ac *AuthControllerImpl) Login(ctx *gin.Context) {
	var req model.LoginRequest

	tr := otel.GetTracerProvider().Tracer("Auth-Login Controller")
	_, span := tr.Start(ctx, "Start Login")
	defer span.End()

	if err := ctx.BindJSON(&req); err != nil {
		helper.NewResponses[any](ctx, http.StatusBadRequest, "Invalid request", req, err, nil)
		return
	}

	data, err := ac.AuthSvc.LoginUser(ctx, &req)
	if err != nil {
		if strings.Contains(err.Error(), string(model.Validation)) {
			helper.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), req, err, nil)
			return
		}

		if strings.Contains(err.Error(), string(model.NotFound)) {
			helper.NewResponses[any](ctx, http.StatusNotFound, err.Error(), req, err, nil)
			return
		}

		helper.NewResponses[any](ctx, http.StatusInternalServerError, "Failed login user", data, err, nil)
		return
	}

	helper.NewResponses[any](ctx, http.StatusOK, "Success login user", data, nil, nil)
}
