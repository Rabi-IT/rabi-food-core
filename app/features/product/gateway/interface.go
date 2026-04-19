package gateway

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/libs/database"
	"github.com/Rabi-IT/rabi-food-core/libs/database/filter"
)

type ProductGateway interface {
	Create(ctx context.Context, input CreateInput) (string, error)
	GetByID(ctx context.Context, filter GetByIDFilter) (*GetByIDOutput, error)
	Patch(ctx context.Context, filter PatchFilter, values PatchValues) (bool, error)
	List(ctx context.Context, filter ListFilter) ([]ListOutput, error)
	Paginate(ctx context.Context, filter PaginateFilter, paginate database.PaginateInput) (PaginateOutput, error)
	Delete(ctx context.Context, filter DeleteFilter) (bool, error)
}

type ListFilter struct {
	IDs      []string    `json:"ids"`
	IsActive filter.Bool `json:"isActive"`
	TenantID string      `json:"tenantId"`
}

type DiscountRule struct {
	QuantityThreshold uint `json:"quantityThreshold"`
	Discount          uint `json:"discount"`
}

type ListOutput struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Price         uint           `json:"price"`
	DiscountRules []DiscountRule `json:"discountRules"`
}

type CreateInput struct {
	TenantID    string `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
	CategoryID  string `json:"categoryId"    validate:"required"`
	// Unit of measurement (e.g., "kg", "liter", "piece")
	Unit          string         `json:"unit"`
	Price         uint           `json:"price"`
	IsActive      bool           `json:"isActive"`
	DiscountRules []DiscountRule `json:"discountRules"`
}

type GetByIDFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"-"`
}

type GetByIDOutput struct {
	ID           string `json:"id"           db:"id"`
	TenantID     string `json:"tenantId"     db:"tenant_id"`
	Name         string `json:"name"         db:"name"`
	Description  string `json:"description"  db:"description"`
	Photo        string `json:"photo"        db:"photo"`
	CategoryID   string `json:"categoryId"   db:"category_id"`
	CategoryName string `json:"categoryName" db:"category_name"`
	Unit         string `json:"unit"         db:"unit"`
	Price        uint   `json:"price"        db:"price"`
	IsActive     bool   `json:"isActive"     db:"is_active"`
}

type PatchFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"-"`
}

type PatchValues struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Photo       *string `json:"photo"`
	CategoryID  *string `json:"categoryId"`
	Unit        *string `json:"unit"`
	Price       *uint   `json:"price"`
	IsActive    *bool   `json:"isActive"`
}

type PaginateFilter struct {
	TenantID   string      `json:"tenantId"`
	Name       string      `json:"name"`
	CategoryID string      `json:"categoryId"`
	IsActive   filter.Bool `json:"isActive"`
}

type PaginateData struct {
	ID          string `json:"id"          db:"id"`
	Name        string `json:"name"        db:"name"`
	Description string `json:"description" db:"description"`
	Photo       string `json:"photo"       db:"photo"`
	CategoryID  string `json:"categoryId"  db:"category_id"`
	Unit        string `json:"unit"        db:"unit"`
	Price       uint   `json:"price"       db:"price"`
	IsActive    bool   `json:"isActive"    db:"is_active"`
}

type PaginateOutput struct {
	Data     []PaginateData `json:"data"`
	MaxPages int            `json:"maxPages"`
}

type DeleteFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"-"`
}
