package fixtures

import (
	"net/http"
	"rabi-food-core/config"
	"rabi-food-core/domain"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

type authFixture struct{}

var Auth = authFixture{}

func (*authFixture) BackofficeToken(t *testing.T, userId string) string {
	t.Helper()

	claims := jwt.MapClaims{
		"user_id":          userId,
		"tenant_id":        "system",
		"name":             "backoffice",
		"email":            "backoffice@backoffice.com",
		"role":             domain.BackofficeRole,
		"original_user_id": userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tk, err := token.SignedString([]byte(config.AuthSecret))
	require.NoError(t, err)

	return tk
}

func (auth *authFixture) UserToken(t *testing.T, id string) string {
	t.Helper()
	backofficeTk := auth.BackofficeToken(t, id)
	user, statusCode := User.GetByID(t, id, backofficeTk)
	require.Equal(t, http.StatusOK, statusCode)

	claims := jwt.MapClaims{
		"user_id":   id,
		"name":      user.Name,
		"email":     user.Email,
		"role":      domain.UserRole,
		"tenant_id": user.TenantID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tk, err := token.SignedString([]byte(config.AuthSecret))
	require.NoError(t, err)

	return tk
}

func (auth *authFixture) StaffToken(t *testing.T, id string) string {
	t.Helper()
	backofficeTk := auth.BackofficeToken(t, id)
	user, statusCode := User.GetByID(t, id, backofficeTk)
	require.Equal(t, http.StatusOK, statusCode)

	claims := jwt.MapClaims{
		"user_id":   id,
		"name":      user.Name,
		"email":     user.Email,
		"role":      domain.StaffRole,
		"tenant_id": user.TenantID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tk, err := token.SignedString([]byte(config.AuthSecret))
	require.NoError(t, err)

	return tk
}

func (auth *authFixture) SystemToken(t *testing.T, tenantID string) string {
	t.Helper()

	claims := jwt.MapClaims{
		"user_id":   "system",
		"name":      "system",
		"email":     "system@system.com",
		"role":      domain.SystemRole,
		"tenant_id": tenantID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tk, err := token.SignedString([]byte(config.AuthSecret))
	require.NoError(t, err)

	return tk
}
