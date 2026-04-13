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
	TenantID string // optional — leave empty for global customers
}

type authFixture struct{}

var Auth = authFixture{}

func (*authFixture) BackofficeToken(t *testing.T, userID string) string {
	t.Helper()

	claims := jwt.MapClaims{
		"sub":   userID,
		"email": "backoffice@backoffice.com",
		"app_metadata": map[string]any{
			"role":             string(auth.Backoffice),
			"tenant_id":        "system",
			"original_user_id": userID,
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

	user := a.fetchUser(t, id)

	claims := jwt.MapClaims{
		"sub":   id,
		"email": user.Email,
		"app_metadata": map[string]any{
			"role":      string(auth.User),
			"tenant_id": user.TenantID,
		},
		"user_metadata": map[string]any{
			"name": user.Name,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tk, err := token.SignedString([]byte(config.AuthSecret))
	require.NoError(t, err)

	return tk
}

func (a *authFixture) StaffToken(t *testing.T, id string) string {
	t.Helper()

	user := a.fetchUser(t, id)

	claims := jwt.MapClaims{
		"sub":   id,
		"email": user.Email,
		"app_metadata": map[string]any{
			"role":      string(auth.Staff),
			"tenant_id": user.TenantID,
		},
		"user_metadata": map[string]any{
			"name": user.Name,
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
			TenantID: input.TenantID,
		}).
		Expect().Status(http.StatusCreated).
		JSON().Object().Decode(out)

	require.NotEmpty(t, out.ID)

	return out.ID
}

// fetchUser retrieves user data from the API using a backoffice token.
func (a *authFixture) fetchUser(t *testing.T, id string) *auth_usecases.GetByIDOutput {
	t.Helper()

	backofficeTk := a.BackofficeToken(t, id)

	out := &auth_usecases.GetByIDOutput{}
	httpexpect.Default(t, AppURL).
		Request(http.MethodGet, "/user/"+id).
		WithHeader("Authorization", "Bearer "+backofficeTk).
		Expect().Status(http.StatusOK).
		JSON().Object().Decode(out)

	return out
}
