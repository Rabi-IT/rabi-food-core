package gateway

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/libs/database"
)

type TenantGateway interface {
	Patch(filter PatchFilter, values PatchValues) (bool, error)
	Create(input CreateInput) (string, error)
	GetByID(id string) (*GetByIDOutput, error)
	Paginate(filter PaginateFilter, paginate database.PaginateInput) (PaginateOutput, error)
	RegisterCustomer(ctx context.Context, tenantID, userID string) error
}

type PatchFilter struct {
	ID *string
}

type PatchValues struct {
	Name *string
}

type CreateInput struct {
	Name string
}

type GetByIDOutput struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type PaginateFilter struct {
	Name *string `json:"name"`
}

type PaginateData struct {
	ID   string `json:"id"   db:"id"`
	Name string `json:"name" db:"name"`
}

type PaginateOutput struct {
	Data     []PaginateData `json:"data"`
	MaxPages int            `json:"maxPages"`
}
