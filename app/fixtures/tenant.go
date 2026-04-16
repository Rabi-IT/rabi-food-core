package fixtures

import (
	"fmt"
	"net/http"
	"testing"

	auth_usecases "github.com/Rabi-IT/rabi-food-core/features/auth/usecases"
	g "github.com/Rabi-IT/rabi-food-core/features/tenant/gateway"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type TenantCreateInput struct {
	Name         string
	UserName     string
	UserPhone    string
	UserEmail    string
	UserPassword string
	UserTaxID    string
}

type TenantCreateOutput struct {
	ID     string
	UserID string
}

type tenantFixture struct {
	URI string
}

var Tenant = tenantFixture{"/tenant/"}

func (tenantFixture) Create(t *testing.T, input *TenantCreateInput) *TenantCreateOutput {
	t.Helper()

	body := input
	if body == nil {
		body = &TenantCreateInput{
			Name:         "Name",
			UserName:     "UserName",
			UserPhone:    "11999999999",
			UserEmail:    fmt.Sprintf("%s@email.com", uuid.NewString()),
			UserPassword: "password123",
			UserTaxID:    "12345678901",
		}
	}

	out := &auth_usecases.SignUpOutput{}
	httpexpect.Default(t, AppURL).
		Request(http.MethodPost, "/auth/signup").
		WithJSON(auth_usecases.SignUpInput{
			Email:    body.UserEmail,
			Password: body.UserPassword,
			Name:     body.UserName,
			Phone:    body.UserPhone,
			TaxID:    body.UserTaxID,
			Tenant:   &auth_usecases.TenantInput{Name: body.Name},
		}).
		Expect().Status(http.StatusCreated).
		JSON().Object().Decode(out)

	require.NotEmpty(t, out.ID)
	require.NotEmpty(t, out.TenantID)

	return &TenantCreateOutput{ID: out.TenantID, UserID: out.ID}
}

func (tenantFixture) EnrollCustomer(t *testing.T, tenantID, token string) {
	t.Helper()

	httpexpect.Default(t, AppURL).
		Request(http.MethodPost, Tenant.URI+tenantID+"/enroll-customers").
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusNoContent)
}

func (tenantFixture) GetMe(t *testing.T, token string) g.GetByIDOutput {
	t.Helper()

	found := g.GetByIDOutput{}

	obj := httpexpect.Default(t, AppURL).
		Request(http.MethodGet, Tenant.URI+"me").
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK)

	response := obj.Raw()
	obj.JSON().Object().Decode(&found)

	err := response.Body.Close()
	require.NoError(t, err)

	return found
}
