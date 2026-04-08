package usecases_test

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/Rabi-IT/rabi-food-core/domain/auth"
	"github.com/Rabi-IT/rabi-food-core/domain/payment_status"
	"github.com/Rabi-IT/rabi-food-core/features/order/gateway"
	"github.com/Rabi-IT/rabi-food-core/features/order/usecases"
	product_gateway "github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	"github.com/Rabi-IT/rabi-food-core/fixtures"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"

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

func (t *TestSuite) Test_OrderIntegration_Create() {
	t.Run("should be able to create an order successfully with valid products", func() {
		EXPECTED_PRODUCT_NAME := "Product Name"
		EXPECTED_PRODUCT_PRICE := 100
		EXPECTED_TOTAL_PRICE := 100

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		productID := fixtures.Product.Create(t.T(), &product_gateway.CreateInput{
			Name:       "Product Name",
			CategoryID: fixtures.Category.Create(t.T(), nil, token),
			Price:      100,
			IsActive:   true,
		}, token)

		Body := usecases.CreateInput{
			Notes: "Notes",
			Items: []usecases.OrderItem{
				{ProductID: productID, Quantity: 1},
			},
		}

		orderID := httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Order.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			Body().NotEmpty().Raw()

		orderFound, httpStatus := fixtures.Order.GetByID(t.T(), orderID, token)
		t.Equal(http.StatusOK, httpStatus)
		t.Equal("Notes", orderFound.Notes)

		t.Len(orderFound.Items, 1)
		t.Equal(productID, orderFound.Items[0].ProductID)
		t.Equal(EXPECTED_PRODUCT_NAME, orderFound.Items[0].ProductName)
		t.EqualValues(EXPECTED_PRODUCT_PRICE, orderFound.Items[0].UnitPrice)
		t.EqualValues(EXPECTED_TOTAL_PRICE, orderFound.TotalPrice)
	})

	t.Run("should correctly calculate total price from multiple items", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		productID1 := fixtures.Product.Create(t.T(), nil, token)
		productID2 := fixtures.Product.Create(t.T(), nil, token)

		Body := usecases.CreateInput{
			Notes: "Notes",
			Items: []usecases.OrderItem{
				{ProductID: productID1, Quantity: 1},
				{ProductID: productID2, Quantity: 2},
			},
		}

		orderID := httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Order.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			Body().NotEmpty().Raw()

		orderFound, httpStatus := fixtures.Order.GetByID(t.T(), orderID, token)
		t.Equal(http.StatusOK, httpStatus)
		t.Equal("Notes", orderFound.Notes)

		t.Len(orderFound.Items, 2)
		t.Equal(productID1, orderFound.Items[0].ProductID)
		t.EqualValues(300, orderFound.TotalPrice)

		t.EqualValues(1, orderFound.Items[0].Quantity)
		t.EqualValues(100, orderFound.Items[0].UnitPrice)
		t.EqualValues(100, orderFound.Items[0].Total)

		t.Equal(productID2, orderFound.Items[1].ProductID)
		t.EqualValues(2, orderFound.Items[1].Quantity)
		t.EqualValues(100, orderFound.Items[1].UnitPrice)
		t.EqualValues(200, orderFound.Items[1].Total)
	})

	t.Run("should fail when no products are found for given IDs", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		nonExistingProductID := uuid.NewString()

		Body := usecases.CreateInput{
			Notes: "Notes",
			Items: []usecases.OrderItem{
				{ProductID: nonExistingProductID, Quantity: 1},
			},
		}

		expectedError := errs.ProductNotFound(nonExistingProductID)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Order.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusNotFound).
			JSON().Object().
			Value("code").IsEqual(expectedError.Code)
	})

	t.Run("should fail when some product IDs are missing in the database", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		existingProductID := fixtures.Product.Create(t.T(), nil, token)

		nonExistingProductID := uuid.NewString()
		Body := usecases.CreateInput{
			Notes: "Notes",
			Items: []usecases.OrderItem{
				{ProductID: existingProductID, Quantity: 1},
				{ProductID: nonExistingProductID, Quantity: 1},
			},
		}

		expectedError := errs.ProductNotFound(nonExistingProductID)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Order.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusNotFound).
			JSON().Object().
			Value("code").IsEqual(expectedError.Code)
	})

	t.Run("should ignore provided tenantID and use token tenant when user is not backoffice", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		productID := fixtures.Product.Create(t.T(), nil, token)
		Body := map[string]any{
			"tenantId": anotherTenant.ID,
			"notes":    "Notes",
			"items": []map[string]any{
				{"productId": productID, "quantity": 1},
			},
		}

		orderID := httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Order.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			Body().NotEmpty().Raw()

		orderFound, httpStatus := fixtures.Order.GetByID(t.T(), orderID, token)
		t.Equal(http.StatusOK, httpStatus)
		t.Equal("Notes", orderFound.Notes)

		t.Len(orderFound.Items, 1)
		t.Equal(productID, orderFound.Items[0].ProductID)
	})

	t.Run("should not be able to create an order with inactive products", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		productID := fixtures.Product.Create(t.T(), &product_gateway.CreateInput{
			Name:       "Inactive Product",
			CategoryID: fixtures.Category.Create(t.T(), nil, token),
			Price:      100,
			IsActive:   false,
		}, token)

		Body := usecases.CreateInput{
			Notes: "Notes",
			Items: []usecases.OrderItem{
				{ProductID: productID, Quantity: 1},
			},
		}

		execpectedError := errs.ProductNotFound(productID)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Order.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusNotFound).
			JSON().Object().
			Value("code").IsEqual(execpectedError.Code)
	})
}

func (t *TestSuite) Test_OrderIntegration_GetByID() {
	t.Run("should be able to get by id", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		orderID := fixtures.Order.Create(t.T(), nil, token)

		found, status := fixtures.Order.GetByID(t.T(), orderID, token)

		t.Equal(http.StatusOK, status)
		t.Equal(orderID, found.ID)
		t.Equal("Notes", found.Notes)
	})

	t.Run("should return NotFound when get by id not found", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		_ = fixtures.Order.Create(t.T(), nil, token)

		NON_EXISTING_ID := uuid.New().String()

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Order.URI+NON_EXISTING_ID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})

	t.Run("should not be able to get a order from another tenant", func() {
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)
		anotherOrderID := fixtures.Order.Create(t.T(), nil, anotherToken)

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Order.URI+anotherOrderID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}

func (t *TestSuite) Test_OrderIntegration_Paginate() {
	t.Run("should be able to paginate", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		for range 15 {
			fixtures.Order.Create(t.T(), nil, token)
		}

		response := new(gateway.PaginateOutput)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Order.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithQueryObject(database.PaginateInput{
				Page:     0,
				PageSize: 10,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Len(response.Data, 10)
		t.Equal(2, response.MaxPages)

		for i := range 10 {
			t.NotEmpty(response.Data[i].ID)
			t.Equal("Notes", response.Data[i].Notes)
		}
	})

	t.Run("should not be able to paginate orders from another tenant", func() {
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)

		for range 5 {
			fixtures.Order.Create(t.T(), nil, anotherToken)
		}

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		response := new(gateway.PaginateOutput)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Order.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithQueryObject(database.PaginateInput{
				Page:     0,
				PageSize: 10,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Empty(response.Data)
		t.Equal(0, response.MaxPages)
	})
}

func (t *TestSuite) Test_OrderIntegration_ConfirmPayment() {
	t.Run("should confirm payment when order is pending", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		orderID := fixtures.Order.Create(t.T(), nil, token)
		systemToken := fixtures.Auth.SystemToken(t.T())

		body := usecases.ConfirmPaymentInput{
			OrderID:           orderID,
			ExternalPaymentID: "any-id",
			Provider:          "any-provider",
			PaidAt:            time.Now().Truncate(time.Microsecond),
		}

		err := fixtures.Order.ConfirmPayment(t.T(), orderID, &body, systemToken)
		t.Require().Nil(err)

		orderFound, httpStatus := fixtures.Order.GetByID(t.T(), orderID, token)
		t.Require().Equal(http.StatusOK, httpStatus)
		t.Require().Equal(payment_status.StatusPaid, orderFound.PaymentStatus)
		t.Require().NotNil(orderFound.PaidAt)
		t.Require().Equal(body.PaidAt.UnixMicro(), orderFound.PaidAt.UnixMicro())
		t.Require().NotNil(orderFound.ExternalPaymentID)
		t.Require().Equal(body.ExternalPaymentID, *orderFound.ExternalPaymentID)
	})

	t.Run("should return ok when confirming payment again with same external payment id (idempotency)", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		orderID := fixtures.Order.Create(t.T(), nil, token)
		systemToken := fixtures.Auth.SystemToken(t.T())

		body := usecases.ConfirmPaymentInput{
			OrderID:           orderID,
			ExternalPaymentID: "any-id",
			Provider:          "any-provider",
			PaidAt:            time.Now().Truncate(time.Microsecond),
		}

		n := 5
		var wg sync.WaitGroup
		errCh := make(chan *errs.AppError, n)
		start := make(chan struct{})
		for range n {
			wg.Go(func() {
				<-start
				errCh <- fixtures.Order.ConfirmPayment(t.T(), orderID, &body, systemToken)
			})
		}

		close(start)
		wg.Wait()
		close(errCh)

		for err := range errCh {
			t.Require().Nil(err)
		}

		orderFound, httpStatus := fixtures.Order.GetByID(t.T(), orderID, token)
		t.Require().Equal(http.StatusOK, httpStatus)
		t.Require().Equal(payment_status.StatusPaid, orderFound.PaymentStatus)
		t.Require().NotNil(orderFound.PaidAt)
		t.Require().Equal(body.PaidAt.UnixMicro(), orderFound.PaidAt.UnixMicro())
		t.Require().NotNil(orderFound.ExternalPaymentID)
		t.Require().Equal(body.ExternalPaymentID, *orderFound.ExternalPaymentID)
	})

	t.Run("should return conflict when confirming payment with different external payment id", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		orderID := fixtures.Order.Create(t.T(), nil, token)
		systemToken := fixtures.Auth.SystemToken(t.T())
		body := usecases.ConfirmPaymentInput{
			OrderID:           orderID,
			ExternalPaymentID: "any-id",
			Provider:          "any-provider",
			PaidAt:            time.Now(),
		}

		err := fixtures.Order.ConfirmPayment(t.T(), orderID, &body, systemToken)
		t.Require().Nil(err)
		// Confirm payment again with different external payment id
		body.ExternalPaymentID = "different-id"
		err = fixtures.Order.ConfirmPayment(t.T(), orderID, &body, systemToken)
		t.Require().NotNil(err)
		t.Require().Equal(errs.ErrPaymentExternalIDConflict.Code, err.Code)
	})

	t.Run("should return conflict when external payment id is already used by another order", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		orderID1 := fixtures.Order.Create(t.T(), nil, token)
		orderID2 := fixtures.Order.Create(t.T(), nil, token)
		systemToken := fixtures.Auth.SystemToken(t.T())

		body1 := usecases.ConfirmPaymentInput{
			OrderID:           orderID1,
			ExternalPaymentID: "any-id",
			Provider:          "any-provider",
			PaidAt:            time.Now(),
		}

		err := fixtures.Order.ConfirmPayment(t.T(), orderID1, &body1, systemToken)
		t.Require().Nil(err)

		// Confirm payment for another order with same external payment id
		body2 := usecases.ConfirmPaymentInput{
			OrderID:           orderID2,
			ExternalPaymentID: "any-id",
			Provider:          "any-provider",
			PaidAt:            time.Now(),
		}

		err = fixtures.Order.ConfirmPayment(t.T(), orderID2, &body2, systemToken)
		t.Require().NotNil(err)
		t.Require().Equal(errs.ErrConflict.Code, err.Code)
	})

	t.Run("should return not found when order does not exist", func() {
		systemToken := fixtures.Auth.SystemToken(t.T())
		NON_EXISTING_ORDER_ID := uuid.NewString()

		body := usecases.ConfirmPaymentInput{
			OrderID:           NON_EXISTING_ORDER_ID,
			ExternalPaymentID: "any-id",
			Provider:          "any-provider",
			PaidAt:            time.Now(),
		}

		err := fixtures.Order.ConfirmPayment(t.T(), NON_EXISTING_ORDER_ID, &body, systemToken)
		t.Require().NotNil(err)
		t.Require().Equal(http.StatusNotFound, err.Status)
	})

	t.Run("should return forbidden when user role is "+string(auth.User), func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userToken := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		orderID := fixtures.Order.Create(t.T(), nil, userToken)

		body := usecases.ConfirmPaymentInput{
			OrderID:           orderID,
			ExternalPaymentID: "any-id",
			Provider:          "any-provider",
			PaidAt:            time.Now(),
		}

		err := fixtures.Order.ConfirmPayment(t.T(), orderID, &body, userToken)
		t.Require().NotNil(err)
		t.Require().Equal(errs.ErrForbidden.Code, err.Code)
	})

	t.Run("should return forbidden when user role is staff", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		staffToken := fixtures.Auth.StaffToken(t.T(), tenant.UserID)
		orderID := fixtures.Order.Create(t.T(), nil, staffToken)

		body := usecases.ConfirmPaymentInput{
			OrderID:           orderID,
			ExternalPaymentID: "any-id",
			Provider:          "any-provider",
			PaidAt:            time.Now(),
		}

		err := fixtures.Order.ConfirmPayment(t.T(), orderID, &body, staffToken)
		t.Require().NotNil(err)
		t.Require().Equal(errs.ErrForbidden.Code, err.Code)
	})

	t.Run("should return invalid transition when confirming payment for canceled order", func() {
		t.T().Skipf("Need to implement cancellation flow first")
	})

	t.Run("should return invalid transition when confirming payment for refunded order", func() {
		t.T().Skipf("Need to implement refund flow first")
	})
}

func (t *TestSuite) Test_OrderIntegration_BackofficePaginate() {
	t.Run("should see orders from all tenants", func() {
		tenant1 := fixtures.Tenant.Create(t.T(), nil)
		token1 := fixtures.Auth.UserToken(t.T(), tenant1.UserID)
		tenant2 := fixtures.Tenant.Create(t.T(), nil)
		token2 := fixtures.Auth.UserToken(t.T(), tenant2.UserID)

		for range 3 {
			fixtures.Order.Create(t.T(), nil, token1)
		}
		for range 3 {
			fixtures.Order.Create(t.T(), nil, token2)
		}

		backofficeToken := fixtures.Auth.BackofficeToken(t.T(), tenant1.UserID)

		response := new(gateway.PaginateOutput)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Order.BackofficeURI).
			WithHeader("Authorization", "Bearer "+backofficeToken).
			WithQueryObject(database.PaginateInput{Page: 0, PageSize: 10}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Len(response.Data, 6)
	})

	t.Run("should filter by tenantId when provided", func() {
		tenant1 := fixtures.Tenant.Create(t.T(), nil)
		token1 := fixtures.Auth.UserToken(t.T(), tenant1.UserID)
		tenant2 := fixtures.Tenant.Create(t.T(), nil)
		token2 := fixtures.Auth.UserToken(t.T(), tenant2.UserID)

		for range 3 {
			fixtures.Order.Create(t.T(), nil, token1)
		}
		for range 2 {
			fixtures.Order.Create(t.T(), nil, token2)
		}

		backofficeToken := fixtures.Auth.BackofficeToken(t.T(), tenant1.UserID)

		response := new(gateway.PaginateOutput)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Order.BackofficeURI).
			WithHeader("Authorization", "Bearer "+backofficeToken).
			WithQueryObject(map[string]any{
				"Page":     0,
				"PageSize": 10,
				"tenantId": tenant1.ID,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Len(response.Data, 3)
	})
}

func (t *TestSuite) Test_OrderIntegration_Delete() {
	t.Run("should be able to delete", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		orderID := fixtures.Order.Create(t.T(), nil, token)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodDelete, fixtures.Order.URI+orderID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNoContent)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Order.URI+orderID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("should not be able to delete a order from another tenant", func() {
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)
		anotherOrderID := fixtures.Order.Create(t.T(), nil, anotherToken)

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodDelete, fixtures.Order.URI+anotherOrderID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}
