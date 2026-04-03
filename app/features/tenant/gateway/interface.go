package gateway

type TenantGateway interface {
	Patch(filter PatchFilter, values PatchValues) (bool, error)
	Create(input CreateInput) (string, error)
	GetByID(id string) (*GetByIDOutput, error)
}

type PatchFilter struct {
	ID *string
}

type PatchValues struct {
	Name string
}

type CreateInput struct {
	Name string
}

type GetByIDOutput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
