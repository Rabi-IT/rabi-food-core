package fixtures

import (
	"net/http"
	"testing"

	"github.com/Rabi-IT/rabi-food-core/config"
	"github.com/Rabi-IT/rabi-food-core/domain/auth"

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
		"role":             auth.Backoffice,
		"original_user_id": userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tk, err := token.SignedString([]byte(config.AuthSecret))
	require.NoError(t, err)

	return tk
}

func (a *authFixture) UserToken(t *testing.T, id string) string {
	t.Helper()
	backofficeTk := a.BackofficeToken(t, id)
	user, statusCode := User.GetByID(t, id, backofficeTk)
	require.Equal(t, http.StatusOK, statusCode)

	claims := jwt.MapClaims{
		"user_id":   id,
		"name":      user.Name,
		"email":     user.Email,
		"role":      auth.User,
		"tenant_id": user.TenantID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tk, err := token.SignedString([]byte(config.AuthSecret))
	require.NoError(t, err)

	return tk
}

func (a *authFixture) StaffToken(t *testing.T, id string) string {
	t.Helper()
	backofficeTk := a.BackofficeToken(t, id)
	user, statusCode := User.GetByID(t, id, backofficeTk)
	require.Equal(t, http.StatusOK, statusCode)

	claims := jwt.MapClaims{
		"user_id":   id,
		"name":      user.Name,
		"email":     user.Email,
		"role":      auth.Staff,
		"tenant_id": user.TenantID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tk, err := token.SignedString([]byte(config.AuthSecret))
	require.NoError(t, err)

	return tk
}

func (a *authFixture) SystemToken(t *testing.T) string {
	t.Helper()

	claims := jwt.MapClaims{
		"user_id": "system",
		"name":    "system",
		"email":   "system@system.com",
		"role":    auth.System,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tk, err := token.SignedString([]byte(config.AuthSecret))
	require.NoError(t, err)

	return tk
}
