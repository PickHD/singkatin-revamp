package controller

import (
	"context"
	"strings"

	"github.com/PickHD/singkatin-revamp/user/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/user/internal/v1/helper"
	"github.com/PickHD/singkatin-revamp/user/internal/v1/middleware"
	"github.com/PickHD/singkatin-revamp/user/internal/v1/model"
	"github.com/PickHD/singkatin-revamp/user/internal/v1/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type (
	// UserController is an interface that has all the function to be implemented inside user controller
	UserController interface {
		Profile(ctx *fiber.Ctx) error
	}

	// UserControllerImpl is an app user struct that consists of all the dependencies needed for user controller
	UserControllerImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		UserSvc service.UserService
	}
)

// NewUserController return new instances user controller
func NewUserController(ctx context.Context, config *config.Configuration, logger *logrus.Logger, userSvc service.UserService) *UserControllerImpl {
	return &UserControllerImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		UserSvc: userSvc,
	}
}

// Check godoc
// @Summary      Get Profiles
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Authorization Bearer <Place Access Token Here>"
// @Success      200  {object}  helper.BaseResponse
// @Failure      404  {object}  helper.BaseResponse
// @Failure      500  {object}  helper.BaseResponse
// @Router       /me [get]
func (uc *UserControllerImpl) Profile(ctx *fiber.Ctx) error {
	data := ctx.Locals(model.KeyJWTValidAccess)
	extData, err := middleware.Extract(data)
	if err != nil {
		return helper.NewResponses[any](ctx, fiber.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	detail, err := uc.UserSvc.GetUserDetail(extData.Email)
	if err != nil {
		if strings.Contains(err.Error(), string(model.NotFound)) {
			return helper.NewResponses[any](ctx, fiber.StatusNotFound, err.Error(), nil, err, nil)
		}

		return helper.NewResponses[any](ctx, fiber.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return helper.NewResponses[any](ctx, fiber.StatusOK, "Success get Profiles", detail, nil, nil)
}
