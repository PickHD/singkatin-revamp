package repository

import (
	"context"

	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	// ShortRepository is an interface that has all the function to be implemented inside short repository
	ShortRepository interface {
		GetListShortenerByUserID(ctx context.Context, userID string) ([]model.Short, error)
		Create(ctx context.Context, req *model.Short) error
	}

	// ShortRepositoryImpl is an app short struct that consists of all the dependencies needed for short repository
	ShortRepositoryImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		DB      *mongo.Database
	}
)

// NewShortRepository return new instances short repository
func NewShortRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, db *mongo.Database) *ShortRepositoryImpl {
	return &ShortRepositoryImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		DB:      db,
	}
}

func (sr *ShortRepositoryImpl) GetListShortenerByUserID(ctx context.Context, userID string) ([]model.Short, error) {
	shorts := []model.Short{}

	cur, err := sr.DB.Collection(sr.Config.Database.ShortenersCollection).Find(ctx, bson.D{{Key: "user_id", Value: userID}})
	if err != nil {
		sr.Logger.Error("ShortRepositoryImpl.GetListShortenerByUserID Find ERROR, ", err)
		return nil, err
	}

	for cur.Next(ctx) {
		var short model.Short

		err := cur.Decode(&short)
		if err != nil {
			sr.Logger.Error("ShortRepositoryImpl.GetListShortenerByUserID Decode ERROR, ", err)
		}

		shorts = append(shorts, short)
	}

	if err := cur.Err(); err != nil {
		sr.Logger.Error("ShortRepositoryImpl.GetListShortenerByUserID Cursors ERROR, ", err)
		return nil, err
	}

	return shorts, nil
}

func (sr *ShortRepositoryImpl) Create(ctx context.Context, req *model.Short) error {
	_, err := sr.DB.Collection(sr.Config.Database.ShortenersCollection).InsertOne(ctx,
		bson.D{{Key: "full_url", Value: req.FullURL}, {Key: "user_id", Value: req.UserID}, {Key: "short_url", Value: req.ShortURL}, {Key: "visited", Value: 0}})
	if err != nil {
		sr.Logger.Error("ShortRepositoryImpl.Create InsertOne ERROR, ", err)
		return err
	}

	return nil
}
