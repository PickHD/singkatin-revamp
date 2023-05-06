package controller

import (
	"context"

	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/service"
	shortenerpb "github.com/PickHD/singkatin-revamp/shortener/pkg/api/v1/proto/shortener"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	// ShortController is an interface that has all the function to be implemented inside short controller
	ShortController interface {
		GetListShortenerByUserID(ctx context.Context, req *shortenerpb.ListShortenerRequest) (*shortenerpb.ListShortenerResponse, error)
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
