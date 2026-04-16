package gateway

import "context"

type AuthGateway interface {
	SignUp(ctx context.Context, input SignUpInput) (*SignUpOutput, error)
	SignIn(ctx context.Context, input SignInInput) (*TokenOutput, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenOutput, error)
	SignOut(ctx context.Context, accessToken string) error
	GetByID(ctx context.Context, id string) (*UserOutput, error)
	Patch(ctx context.Context, id string, input PatchInput) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
	Paginate(ctx context.Context, input PaginateInput) (*PaginateOutput, error)
}

type SignUpInput struct {
	Email        string
	Password     string
	Name         string
	Phone        string
	TaxID        string
	SocialID     string
	City         string
	State        string
	ZIP          string
	Street       string
	Complement   string
	Neighborhood string
	Role string
}

type SignUpOutput struct {
	ID    string
	Email string
}

type SignInInput struct {
	Email    string
	Password string
}

type TokenOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type UserOutput struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	TaxID        string `json:"taxId"`
	SocialID     string `json:"socialId"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZIP          string `json:"zip"`
	Street       string `json:"street"`
	Complement   string `json:"complement"`
	Neighborhood string `json:"neighborhood"`
	Role         string `json:"role"`
}

type PatchInput struct {
	Email        *string
	Name         *string
	Phone        *string
	TaxID        *string
	SocialID     *string
	City         *string
	State        *string
	ZIP          *string
	Street       *string
	Complement   *string
	Neighborhood *string
}

type PaginateInput struct {
	Page     int
	PageSize int
}

type PaginateOutput struct {
	Data     []UserOutput `json:"data"`
	MaxPages int          `json:"maxPages"`
}
