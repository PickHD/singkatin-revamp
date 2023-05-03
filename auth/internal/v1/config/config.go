package config

import "github.com/PickHD/singkatin-revamp/auth/internal/v1/helper"

type (
	Configuration struct {
		Common   *Common
		Database *Database
	}

	Common struct {
		Port int
	}

	Database struct {
		Port int
		Host string
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
	}
}

func NewConfig() *Configuration {
	return loadConfiguration()
}
