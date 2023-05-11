package repository

import (
	"context"

	"github.com/PickHD/singkatin-revamp/user/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/user/internal/v1/model"
	shortenerpb "github.com/PickHD/singkatin-revamp/user/pkg/api/v1/proto/shortener"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/protobuf/proto"
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
		Tracer   *trace.TracerProvider
		DB       *mongo.Database
		RabbitMQ *amqp.Channel
	}
)

// NewUserRepository return new instances user repository
func NewUserRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, tracer *trace.TracerProvider, db *mongo.Database, amqp *amqp.Channel) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		Context:  ctx,
		Config:   config,
		Logger:   logger,
		Tracer:   tracer,
		DB:       db,
		RabbitMQ: amqp,
	}
}

func (ur *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	tr := otel.GetTracerProvider().Tracer("User-FindByEmail Repository")
	_, span := tr.Start(ctx, "Start FindByEmail")
	defer span.End()

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
	tr := otel.GetTracerProvider().Tracer("User-PublishCreateUserShortener Repository")
	_, span := tr.Start(ctx, "Start PublishCreateUserShortener")
	defer span.End()

	ur.Logger.Info("data req before publish", req)

	// transform data to proto
	msg := ur.prepareProtoPublishCreateUserShortenerMessage(req)

	b, err := proto.Marshal(msg)
	if err != nil {
		ur.Logger.Error("UserRepositoryImpl.PublishCreateUserShortener Marshal proto CreateShortenerMessage ERROR, ", err)
		return err
	}

	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(b),
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

func (ur *UserRepositoryImpl) prepareProtoPublishCreateUserShortenerMessage(req *model.GenerateShortUserMessage) *shortenerpb.CreateShortenerMessage {
	return &shortenerpb.CreateShortenerMessage{
		FullUrl:  req.FullURL,
		UserId:   req.UserID,
		ShortUrl: req.ShortURL,
	}
}
