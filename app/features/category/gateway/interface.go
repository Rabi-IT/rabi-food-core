package gateway

import (
	"github.com/Rabi-IT/rabi-food-core/libs/database"
)

type CategoryGateway interface {
	Create(input CreateInput) (string, error)
	GetByID(filter GetByIDFilter) (*GetByIDOutput, error)
	Patch(filter PatchFilter, values PatchValues) (bool, error)
	Paginate(filter PaginateFilter, paginate database.PaginateInput) (PaginateOutput, error)
	Delete(filter DeleteFilter) (bool, error)
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
	ID          string `json:"id"          db:"id"`
	Name        string `json:"name"        db:"name"`
	Description string `json:"description" db:"description"`
	TenantID    string `json:"tenantId"    db:"tenant_id"`
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
	TenantID *string `json:"tenantId"`
	Name     *string `json:"name"`
}

type PaginateData struct {
	ID          string `json:"id"          db:"id"`
	Name        string `json:"name"        db:"name"`
	Description string `json:"description" db:"description"`
}

type PaginateOutput struct {
	Data     []PaginateData `json:"data"`
	MaxPages int            `json:"maxPages"`
}

type DeleteFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"-"`
}
