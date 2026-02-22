package user_case_test

import (
	"net/http"
	"rabi-food-core/fixtures"
	"rabi-food-core/libs/database"
	g "rabi-food-core/libs/database/gateways/user_gateway"
	"rabi-food-core/libs/http/fiber_adapter/middlewares"
	"rabi-food-core/usecases/tenant_case"
	"rabi-food-core/usecases/user_case"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
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

func (t *TestSuite) Test_UserIntegration_Create() {
	t.Run("should be able to create", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := user_case.CreateInput{
			Name:         "Name",
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
			Photo:        "http://example.com/photo.png",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.User.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			Body().NotEmpty()
	})
	t.Run("should fail if required fields are empty", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := user_case.CreateInput{
			// Required fields left empty
			Name:  "",
			Email: "",

			// Optional fields
			TaxID:        "TaxID",
			City:         "City",
			State:        "State",
			Phone:        "Phone",
			ZIP:          "ZIP",
			SocialID:     "SocialID",
			Neighborhood: "Neighborhood",
			Street:       "Street",
			Photo:        "http://example.com/photo.png",
			Complement:   "Complement",
		}

		response := new(middlewares.ValidationErrorResponse)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.User.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Decode(&response)

		t.Len(response.Errors, 2)
		for _, e := range response.Errors {
			switch e.Field {
			case "Name":
				t.Equal("required", e.Tag)
			case "Email":
				t.Equal("required", e.Tag)
			default:
				t.Fail("unexpected validation error field: " + e.Field)
			}
		}
	})

	t.Run("should not fail if all optional fields are not filled in", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := user_case.CreateInput{
			// Required fields
			Name:  "Name",
			Email: "email@email.com",

			// Optional fields left empty
			TenantID:     "",
			Photo:        "",
			TaxID:        "",
			City:         "",
			State:        "",
			Phone:        "",
			ZIP:          "",
			SocialID:     "",
			Neighborhood: "",
			Street:       "",
			Complement:   "",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.User.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			Body().NotEmpty()
	})

	t.Run("should return an error if photo is invalid", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := user_case.CreateInput{
			Name:  "Name",
			Email: "email@email.com",
			Photo: "invalid-url",
		}

		response := new(middlewares.ValidationErrorResponse)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.User.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Decode(response)

		t.Len(response.Errors, 1)
		t.Equal("Photo", response.Errors[0].Field)
		t.Equal("url", response.Errors[0].Tag)
	})

	t.Run("should ignore provided tenantID and use token tenant when user is not backoffice", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		anotherTenant := fixtures.Tenant.Create(t.T(), nil)

		Body := user_case.CreateInput{
			Name:     "Name",
			Email:    "email@email.com",
			TenantID: anotherTenant.ID, // Attempt to set different tenant ID
		}

		userID := fixtures.User.Create(t.T(), &Body, token)

		userFound, httpStatus := fixtures.User.GetByID(t.T(), userID, token)

		t.Equal(http.StatusOK, httpStatus)
		t.Equal(userID, userFound.ID, "UserID should match the created user ID")
		t.Equal(tenant.ID, userFound.TenantID, "TenantID should match the token's tenant ID")
	})
}

func (t *TestSuite) Test_UserIntegration_GetByID() {
	t.Run("should be able to retrieve by id", func() {
		newUser := &tenant_case.CreateInput{
			Name:     "Name",
			UserName: "UserName",
			Phone:    "Phone",
			Email:    "email@email.com",
		}

		tenant := fixtures.Tenant.Create(t.T(), newUser)

		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.User.URI+tenant.UserID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().Object().
			ContainsSubset(map[string]any{
				"tenantId": tenant.ID,
				"name":     newUser.UserName,
				"email":    newUser.Email,
				"phone":    newUser.Phone,
			})
	})

	t.Run("should not be able to get a user from another tenant", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)

		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.User.URI+tenant.UserID).
			WithHeader("Authorization", "Bearer "+anotherToken).
			Expect().
			Status(http.StatusNotFound)
	})
}

func (t *TestSuite) Test_UserIntegration_Paginate() {
	t.Run("should be able to paginate", func() {
		defaultUser := user_case.CreateInput{
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
		tenant := fixtures.Tenant.Create(t.T(), &tenant_case.CreateInput{
			Name:     "Tenant Name",
			UserName: defaultUser.Name,
			Phone:    defaultUser.Phone,
			Email:    defaultUser.Email,
		})
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		userIDs := make(map[string]bool)
		for range 5 {
			userID := fixtures.User.Create(t.T(), &defaultUser, token)
			userIDs[userID] = true
		}

		response := new(g.PaginateOutput)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.User.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithQueryObject(database.PaginateInput{
				Page:     0,
				PageSize: 10,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Len(response.Data, 6) // 5 created + 1 tenant user
		t.Equal(1, response.MaxPages)

		for _, user := range response.Data {
			if user.ID == tenant.UserID {
				continue // Skip tenant user
			}
			t.True(userIDs[user.ID])
			t.Equal(defaultUser.Name, user.Name)
			t.Equal(defaultUser.Photo, user.Photo)
			t.Equal(defaultUser.State, user.State)
			t.Equal(defaultUser.City, user.City)
		}
	})

	t.Run("should be able to paginate users from another tenant if is a backoffice user", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil) // Create 1 users
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		for range 5 {
			fixtures.User.Create(t.T(), nil, token) // Create 5 users
		}

		anotherTenant := fixtures.Tenant.Create(t.T(), nil) // Create 1 user
		backoffToken := fixtures.Auth.BackofficeToken(t.T(), anotherTenant.UserID)
		EXPECTED_AMOUNT_OF_USERS := 7

		response := new(g.PaginateOutput)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.User.URI).
			WithHeader("Authorization", "Bearer "+backoffToken).
			WithQueryObject(database.PaginateInput{
				Page:     0,
				PageSize: 10,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Len(response.Data, EXPECTED_AMOUNT_OF_USERS)
		t.Equal(1, response.MaxPages)
	})

	t.Run("should not be able to paginate users from another tenant", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		for range 5 {
			fixtures.User.Create(t.T(), nil, token)
		}

		// Needed to created 1 user in another tenant
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		EXPECTED_AMOUNT_OF_USERS := 1
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)

		response := new(g.PaginateOutput)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.User.URI).
			WithHeader("Authorization", "Bearer "+anotherToken).
			WithQueryObject(database.PaginateInput{
				Page:     0,
				PageSize: 10,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Len(response.Data, EXPECTED_AMOUNT_OF_USERS)
		t.Equal(1, response.MaxPages)
	})
}

func (t *TestSuite) Test_UserIntegration_Patch() {
	t.Run("should be able to update", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := g.PatchValues{
			ZIP:        "NewZIP",
			Phone:      "NewPhone",
			Email:      "new-email@email.com",
			Street:     "NewStreet",
			SocialID:   "NewSocialID",
			TaxID:      "NewTaxID",
			City:       "NewCity",
			State:      "NewState",
			Complement: "NewComplement",
			Name:       "NewName",
			Photo:      "http://example.com/new-photo.png",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPatch, fixtures.User.URI+tenant.UserID).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusOK).
			Body().NotEmpty()

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.User.URI+tenant.UserID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().Object().
			ContainsSubset(map[string]any{
				"tenantId":   tenant.ID,
				"email":      "new-email@email.com",
				"name":       "NewName",
				"zip":        "NewZIP",
				"phone":      "NewPhone",
				"street":     "NewStreet",
				"socialId":   "NewSocialID",
				"taxId":      "NewTaxID",
				"city":       "NewCity",
				"state":      "NewState",
				"complement": "NewComplement",
				"photo":      "http://example.com/new-photo.png",
			})
	})

	t.Run("should return an error if photo is invalid", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := g.PatchValues{
			Photo: "invalid-url",
		}

		response := &middlewares.ValidationErrorResponse{}
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPatch, fixtures.User.URI+tenant.UserID).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Decode(response)

		t.Len(response.Errors, 1)
		t.Equal("Photo", response.Errors[0].Field)
		t.Equal("url", response.Errors[0].Tag)
	})

	t.Run("should be able to update another user for the same tenant if role is TenantManager", func() {
		t.T().Skip("Skipping until TenantManager role is implemented")
	})

	t.Run("should be able to update another user with backoffice role", func() {
		tenant1 := fixtures.Tenant.Create(t.T(), nil)
		tenant2 := fixtures.Tenant.Create(t.T(), nil)
		backofficeToken := fixtures.Auth.BackofficeToken(t.T(), tenant1.UserID)

		Body := g.PatchValues{
			Name: "NewName",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPatch, fixtures.User.URI+tenant2.UserID).
			WithHeader("Authorization", "Bearer "+backofficeToken).
			WithJSON(Body).
			Expect().
			Status(http.StatusOK).
			Body().NotEmpty()

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.User.URI+tenant2.UserID).
			WithHeader("Authorization", "Bearer "+backofficeToken).
			Expect().
			Status(http.StatusOK).
			JSON().Object().
			ContainsSubset(map[string]any{
				"name": "NewName",
			})
	})

	t.Run("should return not found if user does not exist", func() {
		NON_EXISTENT_ID := uuid.NewString()
		tenant := fixtures.Tenant.Create(t.T(), nil)
		backofficeToken := fixtures.Auth.BackofficeToken(t.T(), tenant.UserID)

		Body := g.PatchValues{
			Name: "NewName",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPatch, fixtures.User.URI+NON_EXISTENT_ID).
			WithHeader("Authorization", "Bearer "+backofficeToken).
			WithJSON(Body).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("should not be able to patch a user from another tenant", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		anotherTenant := fixtures.Tenant.Create(t.T(), nil)

		Body := g.PatchValues{
			Name: "New Name",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPatch, fixtures.User.URI+anotherTenant.UserID).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}

func (t *TestSuite) Test_UserIntegration_Delete() {
	t.Run("should be able to delete", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodDelete, fixtures.User.URI+tenant.UserID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNoContent).
			Body().IsEmpty()

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.User.URI+tenant.UserID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("should not be able to delete a user from another tenant", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		anotherTenant := fixtures.Tenant.Create(t.T(), nil)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodDelete, fixtures.User.URI+anotherTenant.UserID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}
