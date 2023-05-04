package application

import (
	"context"
	"fmt"

	"github.com/PickHD/singkatin-revamp/auth/internal/v1/config"
	"github.com/PickHD/singkatin-revamp/auth/internal/v1/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// App ...
type App struct {
	Application *gin.Engine
	Context     context.Context
	Config      *config.Configuration
	Logger      *logrus.Logger
	DB          *mongo.Database
}

// SetupApplication configuring dependencies app needed
func SetupApplication(ctx context.Context) (*App, error) {
	var err error

	app := &App{}
	app.Context = context.TODO()
	app.Config = config.NewConfig()
	if err != nil {
		return app, err
	}

	// custom log app with logrus
	logWithLogrus := logrus.New()
	logWithLogrus.Formatter = &logrus.JSONFormatter{}
	logWithLogrus.ReportCaller = true
	app.Logger = logWithLogrus

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", app.Config.Database.Host, app.Config.Database.Port)))
	if err != nil {
		app.Logger.Error("failed connect mongoDB, error :", err)
		return app, err
	}

	db := mongoClient.Database(app.Config.Database.Name)
	app.DB = db

	app.Application = gin.New()
	app.Application.Use(middleware.CORSMiddleware())

	app.Logger.Info("APP RUN SUCCESSFULLY")

	return app, nil
}

// Close method will close any instances before app terminated
func (a *App) Close() {
	a.Logger.Info("APP CLOSED SUCCESSFULLY")

	defer func() {
		if err := a.DB.Client().Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()
}
