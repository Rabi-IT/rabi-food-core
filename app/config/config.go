package config

import (
	"os"

	"gorm.io/gorm/logger"
)

var (
	AppPort            = os.Getenv("APP_PORT")
	AuthSecret         = os.Getenv("AUTH_SECRET")
	Env                = environment(os.Getenv("ENV"))
	ProductionDatabase = &DatabaseConfig{
		Host:         os.Getenv("DATABASE_HOST"),
		DatabaseName: os.Getenv("DATABASE_NAME"),
		User:         os.Getenv("DATABASE_USER"),
		Password:     os.Getenv("DATABASE_PASSWORD"),
		Port:         os.Getenv("DATABASE_PORT"),
		LogLevel:     logger.Warn,
	}
)

type DatabaseConfig struct {
	Host string
	User string
	// Password contains database password - gosec ignore
	Password     string // #nosec G117
	DatabaseName string
	Port         string
	LogLevel     logger.LogLevel
}

func (d DatabaseConfig) String() string {
	return "host=" + d.Host + " user=" + d.User + " port=" + d.Port + " database=" + d.DatabaseName + " password=**"
}
