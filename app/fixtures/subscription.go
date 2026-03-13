package fixtures

import (
	"net/http"
	"rabi-food-core/libs/database/gateways/subscription_gateway"
	"rabi-food-core/usecases/subscription_case"
	"testing"

	"github.com/stretchr/testify/require"
)

type subscriptionFixture struct {
	URI string
}

var (
	Subscription          = subscriptionFixture{"/subscription/"}
	DEFAULT_DISCOUNT_RULE = subscription_gateway.DiscountRule{
		CyclesThreshold: 5,
		Discount:        10,
	}
	DEFAULT_CUTOFF_OFFSET_MINUTES  = uint16(0)
	DEFAULT_MAX_ATTEMPTS_PER_ORDER = uint8(1)
)

func (subscriptionFixture) UpsertConfig(t *testing.T, input *subscription_gateway.UpsertConfigInput, token string) {
	t.Helper()
	Body := input
	if Body == nil {
		Body = &subscription_gateway.UpsertConfigInput{
			MaxAttemptsPerOrder: DEFAULT_MAX_ATTEMPTS_PER_ORDER,
			DiscountRules: []subscription_gateway.DiscountRule{
				DEFAULT_DISCOUNT_RULE,
			},
			CutoffOffsetMinutes: DEFAULT_CUTOFF_OFFSET_MINUTES,
		}
	}

	DefaultHTTP(t).
		Request(http.MethodPut, Subscription.URI+"config").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(Body).
		Expect().
		StatusList(http.StatusOK, http.StatusCreated)
}

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
	DefaultHTTP(t).
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

	obj := DefaultHTTP(t).
		Request(http.MethodGet, Subscription.URI+id).
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK)

	response := obj.Raw()

	obj.JSON().Object().Decode(&found)

	err := response.Body.Close()
	require.NoError(t, err)

	return found, response.StatusCode
}
