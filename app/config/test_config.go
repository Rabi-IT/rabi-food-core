package config

import (
	"os"

	"gorm.io/gorm/logger"
)

var (
	TestDatabase = &DatabaseConfig{
		Host:         testEnv("TEST_DATABASE_HOST", "localhost"),
		User:         testEnv("TEST_DATABASE_USER", "postgres"),
		Password:     testEnv("TEST_DATABASE_PASSWORD", "postgres"),
		DatabaseName: testEnv("TEST_DATABASE_NAME", "rabi_food_test"),
		Port:         testEnv("TEST_DATABASE_PORT", "5432"),
		LogLevel:     logger.Warn,
	}
)

func testEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
