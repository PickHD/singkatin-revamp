package config

import "github.com/PickHD/singkatin-revamp/user/internal/v1/helper"

type (
	Configuration struct {
		Server   *Server
		Common   *Common
		Database *Database
		RabbitMQ *RabbitMQ
		Secret   *Secret
		Tracer   *Tracer
	}

	Common struct {
		JWTExpire int
		GRPCPort  string
	}

	Server struct {
		AppPort int
		AppEnv  string
		AppName string
		AppID   string
	}

	Database struct {
		Port                 int
		Host                 string
		Name                 string
		UsersCollection      string
		ShortenersCollection string
	}

	RabbitMQ struct {
		ConnURL              string
		QueueCreateShortener string
		QueueUpdateVisitor   string
	}

	Secret struct {
		JWTSecret string
	}

	Tracer struct {
		JaegerURL string
	}
)

func loadConfiguration() *Configuration {
	return &Configuration{
		Common: &Common{
			JWTExpire: helper.GetEnvInt("JWT_EXPIRE"),
			GRPCPort:  helper.GetEnvString("GRPC_SHORTENER_HOST"),
		},
		Server: &Server{
			AppPort: helper.GetEnvInt("APP_PORT"),
			AppEnv:  helper.GetEnvString("APP_ENV"),
			AppName: helper.GetEnvString("APP_NAME"),
			AppID:   helper.GetEnvString("APP_ID"),
		},
		Database: &Database{
			Port:                 helper.GetEnvInt("DB_PORT"),
			Host:                 helper.GetEnvString("DB_HOST"),
			Name:                 helper.GetEnvString("DB_NAME"),
			UsersCollection:      helper.GetEnvString("DB_COLLECTION_USERS"),
			ShortenersCollection: helper.GetEnvString("DB_COLLECTION_SHORTENERS"),
		},
		RabbitMQ: &RabbitMQ{
			ConnURL:              helper.GetEnvString("AMQP_SERVER_URL"),
			QueueCreateShortener: helper.GetEnvString("AMQP_QUEUE_CREATE_SHORTENER"),
			QueueUpdateVisitor:   helper.GetEnvString("AMQP_QUEUE_UPDATE_VISITOR"),
		},
		Secret: &Secret{
			JWTSecret: helper.GetEnvString("JWT_SECRET"),
		},
		Tracer: &Tracer{
			JaegerURL: helper.GetEnvString("JAEGER_URL"),
		},
	}
}

func NewConfig() *Configuration {
	return loadConfiguration()
}
