package gorm_adapter

import (
	"context"
	"fmt"
	"rabi-food-core/config"
	"rabi-food-core/libs/database"
	"rabi-food-core/libs/database/gorm_adapter/models"
	"rabi-food-core/libs/logger"
	"time"

	gormLogger "gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GormAdapter is the GORM implementation of the Database interface.
type GormAdapter struct {
	Conn   *gorm.DB
	config *config.DatabaseConfig
}

// New creates a new instance of GormAdapter with the given database configuration.
func New(c *config.DatabaseConfig) database.Database {
	return &GormAdapter{config: c}
}

// Migrate performs automatic migration for the database models.
func (g *GormAdapter) Migrate() error {
	return g.Conn.AutoMigrate(
		&models.User{},
		&models.Tenant{},
		&models.Product{},
		&models.Category{},
		&models.Order{},
		&models.AppPlan{},
	)
}

// Connect establishes a connection to the database.
func (g *GormAdapter) Connect(ctx context.Context) error {
	time.Local = time.UTC

	logger.Get(ctx).Info().Msg("Connecting to database with DSN: " + g.config.String())
	dsn := parseDSN(g.config)
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		return err
	}

	g.Conn = db
	g.Conn.Logger = gormLogger.Default.LogMode(g.config.LogLevel)

	return nil
}

// CreateDatabase creates the database if it does not exist.
func (g *GormAdapter) CreateDatabase(ctx context.Context) error {
	dbConfig := &config.DatabaseConfig{
		Host:         g.config.Host,
		User:         g.config.User,
		Password:     g.config.Password,
		Port:         g.config.Port,
		DatabaseName: "postgres",
	}

	logger.Get(ctx).Info().Msg("Creating database with DSN: " + dbConfig.String())
	var dsn = parseDSN(dbConfig)
	conn, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return err
	}

	result := 0
	err = conn.Raw("SELECT 1 from pg_database WHERE datname=?", g.config.DatabaseName).Scan(&result).Error
	if err != nil {
		return err
	}

	hasDatabase := result > 0
	if hasDatabase {
		return nil
	}

	err = conn.Exec("CREATE DATABASE " + g.config.DatabaseName).Error
	if err != nil {
		return err
	}

	sql, err := conn.DB()
	if err != nil {
		return err
	}

	return sql.Close()
}

// Start initializes the database connection and performs migrations.
func (g *GormAdapter) Start(ctx context.Context) error {
	err := g.CreateDatabase(ctx)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	err = g.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	err = g.Migrate()
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

// Stop closes the database connection.
func (g *GormAdapter) Stop() error {
	if g.Conn != nil {
		sqlDB, err := g.Conn.DB()
		if err != nil {
			return err
		}

		return sqlDB.Close()
	}

	return nil
}

func parseDSN(d *config.DatabaseConfig) string {
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s",
		d.Host,
		d.User,
		d.Password,
		d.Port,
		d.DatabaseName,
	)

	return dsn
}
