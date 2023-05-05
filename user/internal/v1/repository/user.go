package repository

import (
	"context"

	"github.com/PickHD/singkatin-revamp/user/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/user/internal/v1/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	// UserRepository is an interface that has all the function to be implemented inside user repository
	UserRepository interface {
		FindByEmail(ctx context.Context, email string) (*model.User, error)
	}

	// UserRepositoryImpl is an app user struct that consists of all the dependencies needed for user repository
	UserRepositoryImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		DB      *mongo.Database
	}
)

// NewUserRepository return new instances user repository
func NewUserRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, db *mongo.Database) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		DB:      db,
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
