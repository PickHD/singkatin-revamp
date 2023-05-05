package application

import (
	"context"
	"fmt"

	"github.com/PickHD/singkatin-revamp/user/internal/v1/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// App ...
type App struct {
	Application *fiber.App
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

	amqpConn, err := amqp.Dial(app.Config.RabbitMQ.ConnURL)
	if err != nil {
		app.Logger.Error("failed dial RabbitMQ, error :", err)
		return app, err
	}
	defer amqpConn.Close()

	amqpClient, err := amqpConn.Channel()
	if err != nil {
		app.Logger.Error("failed open RabbitMQ Channels, error :", err)
		return app, err
	}
	defer amqpClient.Close()

	queues := []string{app.Config.RabbitMQ.QueueCreateShortener, app.Config.RabbitMQ.QueueUpdateVisitor}

	for _, q := range queues {
		_, err = amqpClient.QueueDeclare(
			q,     // queue name
			true,  // durable
			false, // auto delete
			false, // exclusive
			false, // no wait
			nil,   // arguments
		)
		if err != nil {
			app.Logger.Error("failed queue declare Channels, error :", err)
			return nil, err
		}
	}

	app.Application = fiber.New()
	app.Application.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Logger.Info("APP RUN SUCCESSFULLY")

	return app, nil
}

// Close method will close any instances before app terminated
func (a *App) Close() {
	a.Logger.Info("APP CLOSED SUCCESSFULLY")

	defer func() {
		if err := a.DB.Client().Disconnect(a.Context); err != nil {
			panic(err)
		}
	}()
}
