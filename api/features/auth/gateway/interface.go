package gateway

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/domain/auth"
)

type AuthGateway interface {
	CreateUser(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error)
	SignIn(ctx context.Context, input SignInInput) (*TokenOutput, error)
	SendOTP(ctx context.Context, input SendOTPInput) error
	VerifyOTP(ctx context.Context, input VerifyOTPInput) (*TokenOutput, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenOutput, error)
	SignOut(ctx context.Context, accessToken string) error
	GetByID(ctx context.Context, id string) (*UserOutput, error)
	Patch(ctx context.Context, id string, input PatchInput) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
	Paginate(ctx context.Context, input PaginateInput) (*PaginateOutput, error)
}

type CreateUserInput struct {
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
	Role         auth.Role
}

type CreateUserOutput struct {
	ID    string
	Email string
}

type SignInInput struct {
	Email    string
	Password string
}

type SendOTPInput struct {
	Email string
	Name  string
}

type VerifyOTPInput struct {
	Email string
	Code  string
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
