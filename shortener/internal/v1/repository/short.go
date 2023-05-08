package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/model"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	// ShortRepository is an interface that has all the function to be implemented inside short repository
	ShortRepository interface {
		GetListShortenerByUserID(ctx context.Context, userID string) ([]model.Short, error)
		Create(ctx context.Context, req *model.Short) error
		GetByShortURL(ctx context.Context, shortURL string) (*model.Short, error)
		GetFullURLByKey(ctx context.Context, shortURL string) (string, error)
		SetFullURLByKey(ctx context.Context, shortURL string, fullURL string, duration time.Duration) error
		PublishUpdateVisitorCount(ctx context.Context, req *model.UpdateVisitorRequest) error
		UpdateVisitorByShortURL(ctx context.Context, req *model.UpdateVisitorRequest, lastVisitedCount int64) error
	}

	// ShortRepositoryImpl is an app short struct that consists of all the dependencies needed for short repository
	ShortRepositoryImpl struct {
		Context  context.Context
		Config   *config.Configuration
		Logger   *logrus.Logger
		DB       *mongo.Database
		Redis    *redis.Client
		RabbitMQ *amqp.Channel
	}
)

// NewShortRepository return new instances short repository
func NewShortRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, db *mongo.Database, rds *redis.Client, amqp *amqp.Channel) *ShortRepositoryImpl {
	return &ShortRepositoryImpl{
		Context:  ctx,
		Config:   config,
		Logger:   logger,
		DB:       db,
		Redis:    rds,
		RabbitMQ: amqp,
	}
}

func (sr *ShortRepositoryImpl) GetListShortenerByUserID(ctx context.Context, userID string) ([]model.Short, error) {
	shorts := []model.Short{}

	cur, err := sr.DB.Collection(sr.Config.Database.ShortenersCollection).Find(ctx,
		bson.D{{Key: "user_id", Value: userID}},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}, {Key: "_id", Value: -1}}))
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
		bson.D{{Key: "full_url", Value: req.FullURL},
			{Key: "user_id", Value: req.UserID},
			{Key: "short_url", Value: req.ShortURL},
			{Key: "visited", Value: 0}, {Key: "created_at", Value: time.Now()}})
	if err != nil {
		sr.Logger.Error("ShortRepositoryImpl.Create InsertOne ERROR, ", err)
		return err
	}

	return nil
}

func (sr *ShortRepositoryImpl) GetByShortURL(ctx context.Context, shortURL string) (*model.Short, error) {
	short := &model.Short{}

	err := sr.DB.Collection(sr.Config.Database.ShortenersCollection).FindOne(ctx, bson.D{{Key: "short_url", Value: shortURL}}).Decode(&short)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.NewError(model.NotFound, "short_url not found")
		}

		sr.Logger.Error("ShortRepositoryImpl.GetByShortURL FindOne ERROR,", err)
		return nil, err
	}

	return short, nil
}

func (sr *ShortRepositoryImpl) GetFullURLByKey(ctx context.Context, shortURL string) (string, error) {
	result := sr.Redis.Get(ctx, fmt.Sprintf(model.KeyShortURL, shortURL))
	if result.Err() != nil {
		sr.Logger.Error("ShortRepositoryImpl.GetFullURLByKey Get ERROR, ", result.Err())

		return "", result.Err()
	}

	return result.String(), nil
}

func (sr *ShortRepositoryImpl) SetFullURLByKey(ctx context.Context, shortURL string, fullURL string, duration time.Duration) error {
	err := sr.Redis.SetEx(ctx, fmt.Sprintf(model.KeyShortURL, shortURL), fullURL, duration).Err()
	if err != nil {
		sr.Logger.Error("ShortRepositoryImpl.SetFullURLByKey SetEx ERROR, ", err)

		return err
	}

	return nil
}

func (sr *ShortRepositoryImpl) PublishUpdateVisitorCount(ctx context.Context, req *model.UpdateVisitorRequest) error {
	sr.Logger.Info("data req before publish", req)

	b, err := json.Marshal(&req)
	if err != nil {
		sr.Logger.Error("ShortRepositoryImpl.PublishUpdateVisitorCount Marshal JSON ERROR, ", err)
		return err
	}

	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        b,
	}

	// Attempt to publish a message to the queue.
	if err := sr.RabbitMQ.Publish(
		"",                                    // exchange
		sr.Config.RabbitMQ.QueueUpdateVisitor, // queue name
		false,                                 // mandatory
		false,                                 // immediate
		message,                               // message to publish
	); err != nil {
		sr.Logger.Error("ShortRepositoryImpl.PublishUpdateVisitorCount RabbitMQ.Publish ERROR, ", err)
		return err
	}

	sr.Logger.Info("Success Publish Update Visitor Count to Queue: ", sr.Config.RabbitMQ.QueueUpdateVisitor)

	return nil
}

func (sr *ShortRepositoryImpl) UpdateVisitorByShortURL(ctx context.Context, req *model.UpdateVisitorRequest, lastVisitedCount int64) error {
	_, err := sr.DB.Collection(sr.Config.Database.ShortenersCollection).UpdateOne(ctx,
		bson.D{{Key: "short_url", Value: req.ShortURL}}, bson.M{
			"$set": bson.D{{Key: "visited", Value: lastVisitedCount + 1}, {Key: "updated_at", Value: time.Now()}},
		})
	if err != nil {
		sr.Logger.Error("ShortRepositoryImpl.UpdateVisitorByShortURL UpdateOne ERROR, ", err)
		return err
	}

	return nil
}
