package controller

import (
	"context"
	"net/http"

	"github.com/PickHD/singkatin-revamp/user/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/user/internal/v1/helper"
	"github.com/PickHD/singkatin-revamp/user/internal/v1/service"
	"github.com/gofiber/fiber/v2"
)

type (
	// HealthCheckController is an interface that has all the function to be implemented inside health check controller
	HealthCheckController interface {
		Check(ctx *fiber.Ctx) error
	}

	// HealthCheckControllerImpl is an app health check struct that consists of all the dependencies needed for health check controller
	HealthCheckControllerImpl struct {
		Context        context.Context
		Config         *config.Configuration
		HealthCheckSvc service.HealthCheckService
	}
)

// NewHealthCheckController return new instances health check controller
func NewHealthCheckController(ctx context.Context, config *config.Configuration, healthCheckSvc service.HealthCheckService) *HealthCheckControllerImpl {
	return &HealthCheckControllerImpl{
		Context:        ctx,
		Config:         config,
		HealthCheckSvc: healthCheckSvc,
	}
}

// Check godoc
// @Summary      Checking Health Services
// @Tags         Health Check
// @Accept       json
// @Produce      json
// @Success      200  {object}  helper.BaseResponse
// @Failure      500  {object}  helper.BaseResponse
// @Router       /health-check [get]
func (hc *HealthCheckControllerImpl) Check(ctx *fiber.Ctx) error {
	ok, err := hc.HealthCheckSvc.Check()
	if err != nil || !ok {
		return helper.NewResponses[any](ctx, http.StatusInternalServerError, "not OK", ok, err, nil)
	}

	return helper.NewResponses[any](ctx, http.StatusOK, "OK", ok, nil, nil)
}
