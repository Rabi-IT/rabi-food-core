package gateway

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/domain/tenant"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
)

type TenantGateway interface {
	Patch(ctx context.Context, filter PatchFilter, values PatchValues) (bool, error)
	Create(ctx context.Context, input CreateInput) (string, error)
	GetByID(ctx context.Context, id string) (*GetByIDOutput, error)
	GetBySlug(ctx context.Context, slug string) (*GetBySlugOutput, error)
	Paginate(ctx context.Context, filter PaginateFilter, paginate database.PaginateInput) (PaginateOutput, error)
	CreateCustomer(ctx context.Context, tenantID, userID string) error
	CreateMember(ctx context.Context, input CreateMemberInput) error
	GetMember(ctx context.Context, filter GetMemberFilter) (*GetMemberOutput, error)
	GetCustomer(ctx context.Context, filter GetCustomerFilter) (*GetCustomerOutput, error)
}

type PatchFilter struct {
	ID string
}

type PatchValues struct {
	Name *string
}

type InitialMemberInput struct {
	UserID string
	Role   tenant.Role
}

type CreateInput struct {
	Name          string
	Slug          string
	Language      string
	InitialMember InitialMemberInput
}

type GetBySlugOutput struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	Slug     string `db:"slug"`
	Language string `db:"language"`
}

type GetByIDOutput struct {
	ID   string `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
}

type PaginateFilter struct {
	Name string `json:"name"`
}

type PaginateData struct {
	ID   string `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
}

type PaginateOutput struct {
	Data     []PaginateData `json:"data"`
	MaxPages int            `json:"maxPages"`
}

type CreateMemberInput struct {
	TenantID string
	UserID   string
	Role     tenant.Role
}

type GetCustomerFilter struct {
	UserID   string
	TenantID string
}

type GetCustomerOutput struct {
	UserID   string
	TenantID string
}

type GetMemberFilter struct {
	UserID   string
	TenantID string
}

type GetMemberOutput struct {
	UserID   string
	TenantID string
	Role     tenant.Role
}
