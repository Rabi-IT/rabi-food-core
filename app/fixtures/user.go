package fixtures

import (
	"net/http"
	g "rabi-food-core/features/user/gateway"
	c "rabi-food-core/features/user/usecases"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

type userFixture struct {
	URI string
}

var User = userFixture{"/user/"}

func (userFixture) Create(t *testing.T, input *c.CreateInput, token string) string {
	t.Helper()
	Body := input
	if Body == nil {
		Body = &c.CreateInput{
			Name:         "Name",
			Photo:        "http://example.com/photo.png",
			TaxID:        "TaxID",
			City:         "City",
			State:        "State",
			Phone:        "Phone",
			ZIP:          "ZIP",
			SocialID:     "SocialID",
			Email:        "email@email.com",
			Neighborhood: "Neighborhood",
			Street:       "Street",
			Complement:   "Complement",
		}
	}

	id := ""
	httpexpect.Default(t, AppURL).
		Request(http.MethodPost, User.URI).
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(Body).
		Expect().
		Status(http.StatusCreated).
		Body().Decode(&id)

	return id
}

func (userFixture) GetByID(t *testing.T, id string, token string) (g.GetByIDOutput, int) {
	t.Helper()
	found := g.GetByIDOutput{}

	obj := httpexpect.Default(t, AppURL).
		Request(http.MethodGet, User.URI+id).
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK)

	response := obj.Raw()

	obj.JSON().Object().Decode(&found)

	err := response.Body.Close()
	require.NoError(t, err)

	return found, response.StatusCode
}

func (userFixture) GetMe(t *testing.T, token string) g.GetByIDOutput {
	t.Helper()
	found := g.GetByIDOutput{}

	obj := httpexpect.Default(t, AppURL).
		Request(http.MethodGet, User.URI+"me").
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK)

	response := obj.Raw()

	obj.JSON().Object().Decode(&found)

	err := response.Body.Close()
	require.NoError(t, err)

	return found
}
