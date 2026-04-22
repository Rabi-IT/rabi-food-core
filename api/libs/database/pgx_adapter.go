package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Rabi-IT/rabi-food-core/config"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// PgxAdapter is the pgx implementation of the Database interface.
type PgxAdapter struct {
	Pool   *pgxpool.Pool
	config *config.DatabaseConfig
}

// NewPgx creates a new instance of PgxAdapter with the given database configuration.
func NewPgx(c *config.DatabaseConfig) Database {
	return &PgxAdapter{config: c}
}

// migrate runs pending database migrations using Goose.
func (g *PgxAdapter) migrate() error {
	sqlDB := stdlib.OpenDBFromPool(g.Pool)

	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.Up(sqlDB, "migrations")
}

// Connect establishes a connection pool to the database.
func (g *PgxAdapter) Connect(ctx context.Context) error {
	time.Local = time.UTC

	logger.Get(ctx).Info().Msg("Connecting to database: " + g.config.String())

	pool, err := pgxpool.New(ctx, parseDSN(g.config))
	if err != nil {
		return err
	}

	g.Pool = pool

	return nil
}

// Start initializes the database connection and runs migrations.
func (g *PgxAdapter) Start(ctx context.Context) error {
	if err := g.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := g.migrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

// Stop closes the database connection pool.
func (g *PgxAdapter) Stop() error {
	if g.Pool != nil {
		g.Pool.Close()
	}

	return nil
}

func parseDSN(d *config.DatabaseConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s",
		d.Host, d.User, d.Password, d.Port, d.DatabaseName)
}
