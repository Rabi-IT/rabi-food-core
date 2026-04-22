package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

// Paginate executes count and data queries in parallel and scans results into []T.
func Paginate[T any](
	ctx context.Context,
	pool *pgxpool.Pool,
	countSQL string,
	countArgs []any,
	dataSQL string,
	dataArgs []any,
) ([]T, int64, error) {
	var (
		count int64
		rows  []T
		eg    errgroup.Group
	)

	eg.Go(func() error {
		return pool.QueryRow(ctx, countSQL, countArgs...).Scan(&count)
	})

	eg.Go(func() error {
		r, err := pool.Query(ctx, dataSQL, dataArgs...)
		if err != nil {
			return err
		}

		rows, err = pgx.CollectRows(r, pgx.RowToStructByName[T])

		return err
	})

	err := eg.Wait()
	if err != nil {
		return nil, 0, err
	}

	return rows, count, nil
}
