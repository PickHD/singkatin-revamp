package repository

import (
	"context"
	"time"

	"github.com/PickHD/singkatin-revamp/auth/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/auth/internal/v1/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// AuthRepository is an interface that has all the function to be implemented inside auth repository
	AuthRepository interface {
		CreateUser(ctx context.Context, req *model.User) (*model.User, error)
		FindByEmail(ctx context.Context, email string) (*model.User, error)
	}

	// AuthRepositoryImpl is an app auth struct that consists of all the dependencies needed for auth repository
	AuthRepositoryImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		Tracer  *trace.TracerProvider
		DB      *mongo.Database
	}
)

// NewAuthRepository return new instances auth repository
func NewAuthRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, tracer *trace.TracerProvider, db *mongo.Database) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		Tracer:  tracer,
		DB:      db,
	}
}

func (ar *AuthRepositoryImpl) CreateUser(ctx context.Context, req *model.User) (*model.User, error) {
	tr := otel.GetTracerProvider().Tracer("Auth-CreateUser repository")
	_, span := tr.Start(ctx, "Start CreateUser")
	defer span.End()

	// check data users by email is already exists or not
	err := ar.DB.Collection(ar.Config.Database.UsersCollection).FindOne(ctx, bson.D{{Key: "email", Value: req.Email}}).Err()
	if err != nil {
		// if doc not exists, create new one
		if err == mongo.ErrNoDocuments {
			res, err := ar.DB.Collection(ar.Config.Database.UsersCollection).InsertOne(ctx, model.User{
				FullName:  req.FullName,
				Email:     req.Email,
				Password:  req.Password,
				CreatedAt: time.Now(),
			})
			if err != nil {
				ar.Logger.Error("AuthRepositoryImpl.CreateUser InsertOne ERROR, ", err)
				return nil, err
			}

			id, ok := res.InsertedID.(primitive.ObjectID)
			if !ok {
				ar.Logger.Error("AuthRepositoryImpl.CreateUser Type Assertion ERROR, ", err)
				return nil, model.NewError(model.Type, "type assertion error")
			}
			req.ID = id

			return req, nil
		}

		ar.Logger.Error("AuthRepositoryImpl.CreateUser FindOne ERROR, ", err)
		return nil, err

	}

	return nil, model.NewError(model.Validation, "email already exists")
}

func (ar *AuthRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	tr := otel.GetTracerProvider().Tracer("Auth-FindByEmail repository")
	_, span := tr.Start(ctx, "Start FindByEmail")
	defer span.End()

	user := model.User{}

	err := ar.DB.Collection(ar.Config.Database.UsersCollection).FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.NewError(model.NotFound, "users not found")
		}

		ar.Logger.Error("AuthRepositoryImpl.FindByEmail FindOne ERROR, ", err)
		return nil, err
	}

	return &user, nil
}
