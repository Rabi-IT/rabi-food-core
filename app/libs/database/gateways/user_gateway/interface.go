package user_gateway

import (
	"rabi-food-core/domain"
	"rabi-food-core/libs/database"
)

type UserGateway interface {
	Create(input CreateInput) (string, error)
	GetByID(filter GetByIDFilter) (*GetByIDOutput, error)
	Patch(filter PatchFilter, values PatchValues) (bool, error)
	Paginate(filter PaginateFilter, paginate database.PaginateInput) (PaginateOutput, error)
	Delete(filter DeleteFilter) (bool, error)
}

type CreateInput struct {
	TenantID     string
	State        string
	ZIP          string
	Phone        string
	City         string
	Photo        string
	TaxID        string
	SocialID     string
	Street       string
	Complement   string
	Name         string
	Email        string
	Neighborhood string
	Role         domain.Role
}

type GetByIDFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
}

type GetByIDOutput struct {
	ID         string `json:"id"`
	TenantID   string `json:"tenantId"`
	Phone      string `json:"phone"`
	City       string `json:"city"`
	State      string `json:"state"`
	ZIP        string `json:"zip"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Photo      string `json:"photo"`
	TaxID      string `json:"taxId"`
	SocialID   string `json:"socialId"`
	Street     string `json:"street"`
	Complement string `json:"complement"`
}

type PatchFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
}

type PatchValues struct {
	ZIP        string `json:"zip"`
	Phone      string `json:"phone"`
	City       string `json:"city"`
	State      string `json:"state"`
	TaxID      string `json:"taxId"`
	SocialID   string `json:"socialId"`
	Street     string `json:"street"`
	Complement string `json:"complement"`
	Name       string `json:"name"`
	Email      string `validate:"omitempty,email"`
	Photo      string `validate:"omitempty,url"`
}

type PaginateFilter struct {
	TenantID *string
	State    *string
	City     *string
	Name     *string
}

type PaginateData struct {
	ID    string `json:"id"`
	Photo string `json:"photo"`
	Name  string `json:"name"`
	State string `json:"state"`
	City  string `json:"city"`
}

type PaginateOutput struct {
	Data     []PaginateData `json:"data"`
	MaxPages int            `json:"maxPages"`
}

type DeleteFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
}
