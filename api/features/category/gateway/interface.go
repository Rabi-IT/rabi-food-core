package gateway

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/libs/database"
)

type CategoryGateway interface {
	Create(ctx context.Context, input CreateInput) (string, error)
	GetByID(ctx context.Context, filter GetByIDFilter) (*GetByIDOutput, error)
	Patch(ctx context.Context, filter PatchFilter, values PatchValues) (bool, error)
	Paginate(ctx context.Context, filter PaginateFilter, paginate database.PaginateInput) (PaginateOutput, error)
	Delete(ctx context.Context, filter DeleteFilter) (bool, error)
}

type CreateInput struct {
	TenantID    string `json:"-"`
	Name        string `json:"name"        validate:"required"`
	Description string `json:"description"`
}

type GetByIDFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"-"`
}

type GetByIDOutput struct {
	ID          string `db:"id"          json:"id"`
	Name        string `db:"name"        json:"name"`
	Description string `db:"description" json:"description"`
	TenantID    string `db:"tenant_id"   json:"tenantId"`
}

type PatchFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"-"`
}

type PatchValues struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type PaginateFilter struct {
	TenantID string `json:"tenantId"`
	Name     string `json:"name"`
}

type PaginateData struct {
	ID          string `db:"id"          json:"id"`
	Name        string `db:"name"        json:"name"`
	Description string `db:"description" json:"description"`
}

type PaginateOutput struct {
	Data     []PaginateData `json:"data"`
	MaxPages int            `json:"maxPages"`
}

type DeleteFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"-"`
}
