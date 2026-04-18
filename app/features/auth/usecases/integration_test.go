package usecases_test

import (
	"net/http"
	"testing"

	"github.com/Rabi-IT/rabi-food-core/fixtures"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite

	app *fixtures.App
}

func (t *TestSuite) SetupSuite() {
	t.app = fixtures.NewApp()
	t.app.Start(t.T())
}

func (t *TestSuite) SetupSubTest() {
	fixtures.CleanDatabase(t.T())
}

func (t *TestSuite) TearDownSuite() {
	t.app.Stop(t.T())
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (t *TestSuite) Test_AuthIntegration_UserPaginate() {
	t.Run("GET /user should not be accessible by regular users", func() {
		userID := fixtures.Auth.SignUp(t.T(), fixtures.SignUpInput{
			Email:    "user@email.com",
			Password: "password123",
			Name:     "Regular User",
			Phone:    "11988888888",
			TaxID:    "12345678901",
		})
		token := fixtures.Auth.UserToken(t.T(), userID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, "/user").
			WithHeader("Authorization", "Bearer "+token).
			Expect().Status(http.StatusNotFound)
	})
}
