package app_plan_gateway

type AppPlanGateway interface {
	Create(input CreateInput) (string, error)
	GetByID(filter GetByIDFilter) (*GetByIDOutput, error)
	List(input ListFilter) ([]ListOutput, error)
	Patch(filter PatchFilter, values PatchValues) (bool, error)
	Delete(filter DeleteFilter) (bool, error)
}

type CreateInput struct {
	ID          string `json:"id"`
	Name        string `json:"name"        validate:"required"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	IsActive    bool   `json:"isActive"`
}

type GetByIDFilter struct {
	ID string `json:"id"`
}

type GetByIDOutput struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
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

type PatchFilter struct {
	ID string `json:"id"`
}

type PatchValues struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DeleteFilter struct {
	ID string `json:"id"`
}
