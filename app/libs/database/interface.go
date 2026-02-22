package database

import (
	"context"
	"math"
)

// Database defines the interface for a database with start and stop capabilities.
type Database interface {
	Start(ctx context.Context) error
	Stop() error
}

// PaginateInput holds pagination parameters for database queries.
type PaginateInput struct {
	Page     int
	PageSize int
}

// CalcMaxPages calculates the maximum number of pages based on the total count of items.
func (p *PaginateInput) CalcMaxPages(count int64) int {
	total := float64(count) / float64(p.PageSize)

	return int(math.Ceil(total))
}

// Offset calculates the offset for database queries based on the current page and page size.
func (p *PaginateInput) Offset() int {
	return p.Page * p.PageSize
}
