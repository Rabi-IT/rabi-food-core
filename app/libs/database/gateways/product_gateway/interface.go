package product_gateway

import (
	"rabi-food-core/libs/database"
)

type ProductGateway interface {
	Create(input CreateInput) (string, error)
	GetByID(filter GetByIDFilter) (*GetByIDOutput, error)
	Patch(filter PatchFilter, values PatchValues) (bool, error)
	List(filter ListFilter) ([]ListOutput, error)
	Paginate(filter PaginateFilter, paginate database.PaginateInput) (PaginateOutput, error)
	Delete(filter DeleteFilter) (bool, error)
}

type ListFilter struct {
	IDs      []string `json:"ids"`
	IsActive *bool    `json:"isActive"`
	TenantID string   `json:"tenantId"`
}

type ListOutput struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price uint   `json:"price"`
}

type CreateInput struct {
	TenantID    string `json:"tenantId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
	CategoryID  string `json:"categoryId"  validate:"required"`
	Unit        string `json:"unit"`
	Price       uint   `json:"price"`
	IsActive    bool   `json:"isActive"`
}

type GetByIDFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
}

type GetByIDOutput struct {
	ID           string `json:"id"`
	TenantID     string `json:"tenantId"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Photo        string `json:"photo"`
	CategoryID   string `json:"categoryId"`
	CategoryName string `json:"categoryName"`
	Unit         string `json:"unit"`
	Price        uint   `json:"price"`
	IsActive     bool   `json:"isActive"`
}

type PatchFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
}

type PatchValues struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
	CategoryID  string `json:"categoryId"`
	Unit        string `json:"unit"`
	Price       uint   `json:"price"`
	IsActive    bool   `json:"isActive"`
}

type PaginateFilter struct {
	TenantID   *string `json:"tenantId"`
	Name       *string `json:"name"`
	CategoryID *string `json:"categoryId"`
	IsActive   *bool   `json:"isActive"`
}

type PaginateData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
	CategoryID  string `json:"categoryId"`
	Unit        string `json:"unit"`
	Price       uint   `json:"price"`
	IsActive    bool   `json:"isActive"`
}

type PaginateOutput struct {
	Data     []PaginateData `json:"data"`
	MaxPages int            `json:"maxPages"`
}

type DeleteFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
}
