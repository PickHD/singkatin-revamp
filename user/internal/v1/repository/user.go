package repository

import (
	"context"
	"encoding/json"

	"github.com/PickHD/singkatin-revamp/user/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/user/internal/v1/model"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	// UserRepository is an interface that has all the function to be implemented inside user repository
	UserRepository interface {
		FindByEmail(ctx context.Context, email string) (*model.User, error)
		PublishCreateUserShortener(ctx context.Context, req *model.GenerateShortUserMessage) error
	}

	// UserRepositoryImpl is an app user struct that consists of all the dependencies needed for user repository
	UserRepositoryImpl struct {
		Context  context.Context
		Config   *config.Configuration
		Logger   *logrus.Logger
		DB       *mongo.Database
		RabbitMQ *amqp.Channel
	}
)

// NewUserRepository return new instances user repository
func NewUserRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, db *mongo.Database, amqp *amqp.Channel) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		Context:  ctx,
		Config:   config,
		Logger:   logger,
		DB:       db,
		RabbitMQ: amqp,
	}
}

func (ur *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := model.User{}

	err := ur.DB.Collection(ur.Config.Database.UsersCollection).FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.NewError(model.NotFound, "users not found")
		}

		ur.Logger.Error("UserRepositoryImpl.FindByEmail FindOne ERROR, ", err)
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepositoryImpl) PublishCreateUserShortener(ctx context.Context, req *model.GenerateShortUserMessage) error {
	ur.Logger.Info("data req before publish", req)

	b, err := json.Marshal(&req)
	if err != nil {
		ur.Logger.Error("UserRepositoryImpl.PublishCreateUserShortener Marshal JSON ERROR, ", err)
		return err
	}

	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        b,
	}

	// Attempt to publish a message to the queue.
	if err := ur.RabbitMQ.Publish(
		"",                                      // exchange
		ur.Config.RabbitMQ.QueueCreateShortener, // queue name
		false,                                   // mandatory
		false,                                   // immediate
		message,                                 // message to publish
	); err != nil {
		ur.Logger.Error("UserRepositoryImpl.PublishCreateUserShortener RabbitMQ.Publish ERROR, ", err)
		return err
	}

	ur.Logger.Info("Success Publish User Shortener to Queue: ", ur.Config.RabbitMQ.QueueCreateShortener)

	return nil
}
