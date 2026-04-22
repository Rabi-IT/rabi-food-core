package fixtures

import (
	"fmt"
	"net/http"
	"testing"

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

	userID := Auth.SignUp(t, SignUpInput{
		Email:    body.UserEmail,
		Password: body.UserPassword,
		Name:     body.UserName,
		Phone:    body.UserPhone,
		TaxID:    body.UserTaxID,
	})

	token := Auth.UserToken(t, userID)

	tenantID := httpexpect.Default(t, ApiURL).
		Request(http.MethodPost, "/tenant").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]string{"name": body.Name}).
		Expect().Status(http.StatusCreated).
		Body().Raw()

	require.NotEmpty(t, tenantID)

	return &TenantCreateOutput{ID: tenantID, UserID: userID}
}

func (tenantFixture) EnrollCustomer(t *testing.T, tenantID, token string) {
	t.Helper()

	httpexpect.Default(t, ApiURL).
		Request(http.MethodPost, Tenant.URI+tenantID+"/enroll-customers").
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusNoContent)
}

func (tenantFixture) GetMe(t *testing.T, token string) g.GetByIDOutput {
	t.Helper()

	found := g.GetByIDOutput{}

	obj := httpexpect.Default(t, ApiURL).
		Request(http.MethodGet, Tenant.URI+"me").
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK)

	response := obj.Raw()
	obj.JSON().Object().Decode(&found)

	err := response.Body.Close()
	require.NoError(t, err)

	return found
}
