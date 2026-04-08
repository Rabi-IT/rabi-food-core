package fixtures

import (
	"net/http"
	"testing"

	"github.com/Rabi-IT/rabi-food-core/features/order"
	g "github.com/Rabi-IT/rabi-food-core/features/order/gateway"
	c "github.com/Rabi-IT/rabi-food-core/features/order/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

type orderFixture struct {
	URI          string
	BackofficeURI string
}

var Order = orderFixture{URI: "/order/", BackofficeURI: "/backoffice/order/"}

func (orderFixture) Create(t *testing.T, input *c.CreateInput, token string) string {
	t.Helper()
	Body := input
	if Body == nil {
		Body = &c.CreateInput{
			Items: []c.OrderItem{
				{
					ProductID: Product.Create(t, nil, token),
					Quantity:  1,
				},
			},
			Notes: "Notes",
		}
	}

	id := ""
	httpexpect.Default(t, AppURL).
		Request(http.MethodPost, Order.URI).
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(Body).
		Expect().
		Status(http.StatusCreated).
		Body().Decode(&id)

	return id
}

func (orderFixture) GetByID(t *testing.T, id string, token string) (g.GetByIDOutput, int) {
	t.Helper()
	require.NotEmpty(t, id)

	found := g.GetByIDOutput{}

	obj := httpexpect.Default(t, AppURL).
		Request(http.MethodGet, Order.URI+id).
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK)

	response := obj.Raw()

	obj.JSON().Object().Decode(&found)

	err := response.Body.Close()
	require.NoError(t, err)

	return found, response.StatusCode
}

func (orderFixture) ExpectFulfillmentStatus(
	t *testing.T,
	id string,
	expectedStatus order.FulfillmentStatus,
	token string,
) {
	t.Helper()
	found, httpStatus := Order.GetByID(t, id, token)
	require.Equal(t, http.StatusOK, httpStatus)
	require.Equal(t, expectedStatus, found.FulfillmentStatus)
}

func (orderFixture) Patch(t *testing.T, id string, input *g.PatchValues, token string) *errs.AppError {
	t.Helper()
	require.NotEmpty(t, id)

	obj := httpexpect.Default(t, AppURL).
		Request(http.MethodPatch, Order.URI+id).
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(input).
		Expect()

	raw := obj.Raw()

	if raw.StatusCode != http.StatusOK && raw.StatusCode != http.StatusNotFound {
		appErr := errs.AppError{}
		obj.JSON().Object().Decode(&appErr)
		err := raw.Body.Close()
		require.NoError(t, err)

		return &appErr
	}

	return nil
}

func (orderFixture) ConfirmPayment(
	t *testing.T,
	orderID string,
	input *c.ConfirmPaymentInput,
	token string,
) *errs.AppError {
	t.Helper()

	obj := httpexpect.Default(t, AppURL).
		Request(http.MethodPost, Order.URI+orderID+"/payments/confirm").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(input).
		Expect()

	raw := obj.Raw()

	if raw.StatusCode == http.StatusNotFound {
		return &errs.AppError{Status: raw.StatusCode}
	}

	if raw.StatusCode != http.StatusOK {
		appErr := errs.AppError{}
		obj.JSON().Object().Decode(&appErr)
		err := raw.Body.Close()
		require.NoError(t, err)

		return &appErr
	}

	return nil
}
