package fixtures

import (
	"net/http"
	"testing"

	"github.com/Rabi-IT/rabi-food-core/config"
	"github.com/Rabi-IT/rabi-food-core/domain/auth"
	auth_usecases "github.com/Rabi-IT/rabi-food-core/features/auth/usecases"

	"github.com/gavv/httpexpect/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

type SignUpInput struct {
	Email    string
	Password string
	Name     string
	Phone    string
	TaxID    string
}

type RequestContext struct {
	Token    string
	TenantID string
}

type authFixture struct{}

var Auth = authFixture{}

func (*authFixture) BackofficeToken(t *testing.T, userID string) string {
	t.Helper()

	claims := jwt.MapClaims{
		"sub":   userID,
		"email": "backoffice@backoffice.com",
		"app_metadata": map[string]any{
			"role":      string(auth.Backoffice),
			"tenant_id": "system",
		},
		"user_metadata": map[string]any{
			"name": "backoffice",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tk, err := token.SignedString([]byte(config.AuthSecret))
	require.NoError(t, err)

	return tk
}

func (a *authFixture) UserToken(t *testing.T, id string) string {
	t.Helper()
	return a.signToken(t, id, auth.User)
}

func (a *authFixture) TenantOwnerToken(t *testing.T, id, tenantID string) string {
	t.Helper()

	claims := jwt.MapClaims{
		"sub": id,
		"app_metadata": map[string]any{
			"role":      string(auth.TenantOwner),
			"tenant_id": tenantID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tk, err := token.SignedString([]byte(config.AuthSecret))
	require.NoError(t, err)

	return tk
}

func (a *authFixture) UserAuth(t *testing.T, tenantID, userID string) RequestContext {
	t.Helper()
	return RequestContext{Token: a.UserToken(t, userID), TenantID: tenantID}
}

func (a *authFixture) TenantOwnerAuth(t *testing.T, tenantID, userID string) RequestContext {
	t.Helper()
	return RequestContext{Token: a.TenantOwnerToken(t, userID, tenantID), TenantID: tenantID}
}

func (a *authFixture) signToken(t *testing.T, id string, role auth.Role) string {
	t.Helper()

	claims := jwt.MapClaims{
		"sub": id,
		"app_metadata": map[string]any{
			"role": string(role),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tk, err := token.SignedString([]byte(config.AuthSecret))
	require.NoError(t, err)

	return tk
}

func (a *authFixture) SystemToken(t *testing.T) string {
	t.Helper()

	claims := jwt.MapClaims{
		"sub":   "system",
		"email": "system@system.com",
		"app_metadata": map[string]any{
			"role": string(auth.System),
		},
		"user_metadata": map[string]any{
			"name": "system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tk, err := token.SignedString([]byte(config.AuthSecret))
	require.NoError(t, err)

	return tk
}

func (a *authFixture) SignUp(t *testing.T, input SignUpInput) string {
	t.Helper()

	out := &auth_usecases.SignUpOutput{}
	httpexpect.Default(t, AppURL).
		Request(http.MethodPost, "/auth/signup").
		WithJSON(auth_usecases.SignUpInput{
			Email:    input.Email,
			Password: input.Password,
			Name:     input.Name,
			Phone:    input.Phone,
			TaxID:    input.TaxID,
		}).
		Expect().Status(http.StatusCreated).
		JSON().Object().Decode(out)

	require.NotEmpty(t, out.ID)

	return out.ID
}

