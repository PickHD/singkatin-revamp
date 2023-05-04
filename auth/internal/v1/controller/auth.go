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
)

type (
	// Authcontroller is an interface that has all the function to be implemented inside auth controller
	AuthController interface {
		Register(ctx *gin.Context)
	}

	// AuthcontrollerImpl is an app auth struct that consists of all the dependencies needed for auth controller
	AuthControllerImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		AuthSvc service.AuthService
	}
)

// NewAuthController return new instances auth controller
func NewAuthController(ctx context.Context, config *config.Configuration, logger *logrus.Logger, authSvc service.AuthService) *AuthControllerImpl {
	return &AuthControllerImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
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
