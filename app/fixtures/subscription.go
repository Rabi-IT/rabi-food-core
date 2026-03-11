package fixtures

import (
	"net/http"
	"rabi-food-core/libs/database/gateways/subscription_gateway"
	"rabi-food-core/usecases/subscription_case"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

type subscriptionFixture struct {
	URI string
}

var Subscription = subscriptionFixture{"/subscription/"}

func (subscriptionFixture) Create(t *testing.T, input *subscription_case.CreateInput, token string) string {
	t.Helper()
	Body := input
	if Body == nil {
		Body = &subscription_case.CreateInput{
			Items: []subscription_case.SubscriptionItemInput{
				{
					ProductID: Product.Create(t, nil, token),
					Quantity:  1,
				},
			},
			DeliveryDays: []subscription_gateway.DeliveryDay{
				{
					Weekday:   0,
					StartHour: 10, //nolint:mnd
					EndHour:   12, //nolint:mnd
				},
			},
			TotalCycles: 10, //nolint:mnd
			AutoRenew:   true,
			Notes:       "Notes",
		}
	}

	id := ""
	httpexpect.Default(t, AppURL).
		Request(http.MethodPost, Subscription.URI).
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(Body).
		Expect().
		Status(http.StatusCreated).
		Body().Decode(&id)

	return id
}

func (subscriptionFixture) GetByID(t *testing.T, id string, token string) (subscription_gateway.GetByIDOutput, int) {
	t.Helper()
	found := subscription_gateway.GetByIDOutput{}

	obj := httpexpect.Default(t, AppURL).
		Request(http.MethodGet, Subscription.URI+id).
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK)

	response := obj.Raw()

	obj.JSON().Object().Decode(&found)

	err := response.Body.Close()
	require.NoError(t, err)

	return found, response.StatusCode
}
