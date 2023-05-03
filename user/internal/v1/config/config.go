package config

import "github.com/PickHD/singkatin-revamp/user/internal/v1/helper"

type (
	Configuration struct {
		Common   *Common
		Database *Database
		RabbitMQ *RabbitMQ
	}

	Common struct {
		Port int
	}

	Database struct {
		Port int
		Host string
	}

	RabbitMQ struct {
		ConnURL              string
		QueueCreateShortener string
		QueueUpdateVisitor   string
	}
)

func loadConfiguration() *Configuration {
	return &Configuration{
		Common: &Common{
			Port: helper.GetEnvInt("APP_PORT"),
		},
		Database: &Database{
			Port: helper.GetEnvInt("DB_PORT"),
			Host: helper.GetEnvString("DB_HOST"),
		},
		RabbitMQ: &RabbitMQ{
			ConnURL:              helper.GetEnvString("AMQP_SERVER_URL"),
			QueueCreateShortener: helper.GetEnvString("AMQP_QUEUE_CREATE_SHORTENER"),
			QueueUpdateVisitor:   helper.GetEnvString("AMQP_QUEUE_UPDATE_VISITOR"),
		},
	}
}

func NewConfig() *Configuration {
	return loadConfiguration()
}
