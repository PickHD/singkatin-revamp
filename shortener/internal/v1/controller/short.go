package controller

import (
	"context"
	"net/http"
	"strings"

	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/helper"
	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/model"
	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/service"
	shortenerpb "github.com/PickHD/singkatin-revamp/shortener/pkg/api/v1/proto/shortener"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	// ShortController is an interface that has all the function to be implemented inside short controller
	ShortController interface {
		// grpc
		GetListShortenerByUserID(ctx context.Context, req *shortenerpb.ListShortenerRequest) (*shortenerpb.ListShortenerResponse, error)

		// http
		ClickShortener(ctx echo.Context) error

		// rabbitmq
		ProcessCreateShortUser(ctx context.Context, req *model.CreateShortRequest) error
		ProcessUpdateVisitorCount(ctx context.Context, req *model.UpdateVisitorRequest) error
	}

	// ShortControllerImpl is an app short struct that consists of all the dependencies needed for short controller
	ShortControllerImpl struct {
		Context  context.Context
		Config   *config.Configuration
		Logger   *logrus.Logger
		ShortSvc service.ShortService
		shortenerpb.UnimplementedShortenerServiceServer
	}
)

// NewShortController return new instances short controller
func NewShortController(ctx context.Context, config *config.Configuration, logger *logrus.Logger, shortSvc service.ShortService) *ShortControllerImpl {
	return &ShortControllerImpl{
		Context:  ctx,
		Config:   config,
		Logger:   logger,
		ShortSvc: shortSvc,
	}
}

func (sc *ShortControllerImpl) GetListShortenerByUserID(ctx context.Context, req *shortenerpb.ListShortenerRequest) (*shortenerpb.ListShortenerResponse, error) {
	data, err := sc.ShortSvc.GetListShortenerByUserID(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed Get List Shortener By UserID %s", err.Error())
	}

	if len(data) < 1 {
		return &shortenerpb.ListShortenerResponse{}, nil
	}

	shorteners := make([]*shortenerpb.Shortener, len(data))

	for i, q := range data {
		shorteners[i] = &shortenerpb.Shortener{
			Id:       q.ID.Hex(),
			FullUrl:  q.FullURL,
			ShortUrl: q.ShortURL,
			Visited:  q.Visited,
		}
	}

	return &shortenerpb.ListShortenerResponse{
		Shorteners: shorteners,
	}, nil
}

// Check godoc
// @Summary      Click Shorteners URL
// @Tags         Shortener
// @Accept       json
// @Produce      json
// @Param        short_url   path string  true  "short urls"
// @Success      301  {object}  helper.BaseResponse
// @Failure      400  {object}  helper.BaseResponse
// @Failure      404  {object}  helper.BaseResponse
// @Failure      500  {object}  helper.BaseResponse
// @Router       /{short_url} [get]
func (sc *ShortControllerImpl) ClickShortener(ctx echo.Context) error {
	data, err := sc.ShortSvc.ClickShort(ctx.Param("short_url"))
	if err != nil {
		if strings.Contains(err.Error(), string(model.Validation)) {
			return helper.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), ctx.Param("short_url"), err, nil)
		}

		if strings.Contains(err.Error(), string(model.NotFound)) {
			return helper.NewResponses[any](ctx, http.StatusNotFound, err.Error(), ctx.Param("short_url"), err, nil)
		}

		return helper.NewResponses[any](ctx, http.StatusInternalServerError, "failed click shortener", ctx.Param("short_url"), err, nil)
	}

	return ctx.Redirect(http.StatusTemporaryRedirect, data.FullURL)
}

func (sc *ShortControllerImpl) ProcessCreateShortUser(ctx context.Context, req *model.CreateShortRequest) error {
	err := sc.ShortSvc.CreateShort(ctx, req)
	if err != nil {
		return model.NewError(model.Internal, err.Error())
	}

	return nil
}

func (sc *ShortControllerImpl) ProcessUpdateVisitorCount(ctx context.Context, req *model.UpdateVisitorRequest) error {
	err := sc.ShortSvc.UpdateVisitorShort(ctx, req)
	if err != nil {
		return model.NewError(model.Internal, err.Error())
	}

	return nil
}
