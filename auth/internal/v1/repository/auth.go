package repository

import (
	"context"

	"github.com/PickHD/singkatin-revamp/auth/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/auth/internal/v1/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		DB      *mongo.Database
	}
)

// NewAuthRepository return new instances auth repository
func NewAuthRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, db *mongo.Database) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		DB:      db,
	}
}

func (ar *AuthRepositoryImpl) CreateUser(ctx context.Context, req *model.User) (*model.User, error) {
	// check data users by email is already exists or not
	err := ar.DB.Collection(ar.Config.Database.UsersCollection).FindOne(ctx, bson.D{{Key: "email", Value: req.Email}}).Err()
	if err != nil {
		// if doc not exists, create new one
		if err == mongo.ErrNoDocuments {
			res, err := ar.DB.Collection(ar.Config.Database.UsersCollection).InsertOne(ctx, model.User{
				FullName: req.FullName,
				Email:    req.Email,
				Password: req.Password,
			})
			if err != nil {
				ar.Logger.Error("AuthRepositoryImpl.CreateUser InsertOne ERROR, ", err)
				return nil, err
			}

			id, ok := res.InsertedID.(primitive.ObjectID)
			if !ok {
				ar.Logger.Error("AuthRepositoryImpl.CreateUser Type Assertion ERROR, ", err)
				return nil, model.NewError("Type", "type assertion error")
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
