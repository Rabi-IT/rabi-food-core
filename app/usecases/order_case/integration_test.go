package order_case_test

import (
	"net/http"
	"rabi-food-core/domain"
	"rabi-food-core/domain/order"
	"rabi-food-core/fixtures"
	"rabi-food-core/libs/database"
	"rabi-food-core/libs/database/gateways/order_gateway"
	"rabi-food-core/libs/database/gateways/product_gateway"
	"rabi-food-core/libs/errs"
	"rabi-food-core/usecases/order_case"
	"sync"
	"testing"
	"time"

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

		Body := order_case.CreateInput{
			Notes: "Notes",
			Items: []order_case.OrderItem{
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

		Body := order_case.CreateInput{
			Notes: "Notes",
			Items: []order_case.OrderItem{
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

		Body := order_case.CreateInput{
			Notes: "Notes",
			Items: []order_case.OrderItem{
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
		Body := order_case.CreateInput{
			Notes: "Notes",
			Items: []order_case.OrderItem{
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

		Body := order_case.CreateInput{
			Notes: "Notes",
			Items: []order_case.OrderItem{
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

		response := new(order_gateway.PaginateOutput)
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

		response := new(order_gateway.PaginateOutput)
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

		body := order_case.ConfirmPaymentInput{
			OrderID:           orderID,
			ExternalPaymentID: "any-id",
			Provider:          "any-provider",
			PaidAt:            time.Now().Truncate(time.Microsecond),
		}

		err := fixtures.Order.ConfirmPayment(t.T(), orderID, &body, systemToken)
		t.Require().NoError(err)

		orderFound, httpStatus := fixtures.Order.GetByID(t.T(), orderID, token)
		t.Require().Equal(http.StatusOK, httpStatus)
		t.Require().Equal(order.PaymentPaid, orderFound.PaymentStatus)
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

		body := order_case.ConfirmPaymentInput{
			OrderID:           orderID,
			ExternalPaymentID: "any-id",
			Provider:          "any-provider",
			PaidAt:            time.Now().Truncate(time.Microsecond),
		}

		n := 5
		var wg sync.WaitGroup
		errCh := make(chan error, n)
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
			t.Require().NoError(err)
		}

		orderFound, httpStatus := fixtures.Order.GetByID(t.T(), orderID, token)
		t.Require().Equal(http.StatusOK, httpStatus)
		t.Require().Equal(order.PaymentPaid, orderFound.PaymentStatus)
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
		body := order_case.ConfirmPaymentInput{
			OrderID:           orderID,
			ExternalPaymentID: "any-id",
			Provider:          "any-provider",
			PaidAt:            time.Now(),
		}

		err := fixtures.Order.ConfirmPayment(t.T(), orderID, &body, systemToken)
		t.Require().NoError(err)
		// Confirm payment again with different external payment id
		body.ExternalPaymentID = "different-id"
		err = fixtures.Order.ConfirmPayment(t.T(), orderID, &body, systemToken)
		t.Require().NotNil(err)
		t.Require().Equal(errs.ErrConflict.Code, err.Code)
	})

	t.Run("should return conflict when external payment id is already used by another order", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		orderID1 := fixtures.Order.Create(t.T(), nil, token)
		orderID2 := fixtures.Order.Create(t.T(), nil, token)
		systemToken := fixtures.Auth.SystemToken(t.T())

		body1 := order_case.ConfirmPaymentInput{
			OrderID:           orderID1,
			ExternalPaymentID: "any-id",
			Provider:          "any-provider",
			PaidAt:            time.Now(),
		}

		err := fixtures.Order.ConfirmPayment(t.T(), orderID1, &body1, systemToken)
		t.Require().NoError(err)

		// Confirm payment for another order with same external payment id
		body2 := order_case.ConfirmPaymentInput{
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

		body := order_case.ConfirmPaymentInput{
			OrderID:           NON_EXISTING_ORDER_ID,
			ExternalPaymentID: "any-id",
			Provider:          "any-provider",
			PaidAt:            time.Now(),
		}

		err := fixtures.Order.ConfirmPayment(t.T(), NON_EXISTING_ORDER_ID, &body, systemToken)
		t.Require().NotNil(err)
		t.Require().Equal(http.StatusNotFound, err.Status)
	})

	t.Run("should return forbidden when user role is "+string(domain.UserRole), func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userToken := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		orderID := fixtures.Order.Create(t.T(), nil, userToken)

		body := order_case.ConfirmPaymentInput{
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

		body := order_case.ConfirmPaymentInput{
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

func (t *TestSuite) Test_OrderIntegration_Patch_status() {
	tests := []struct {
		name string
		flow []order.FulfillmentStatus
		err  *errs.AppError
	}{
		{
			name: "Pending > Confirmed > In Production > Ready",
			flow: []order.FulfillmentStatus{
				order.FulfillmentConfirmed,
				order.FulfillmentInProduction,
				order.FulfillmentReady,
			},
		},
		{
			name: "Pending > In Production > error",
			flow: []order.FulfillmentStatus{
				order.FulfillmentInProduction,
			},
			err: errs.InvalidTranstion(order.FulfillmentPending, order.FulfillmentInProduction),
		},
		{
			name: "Pending > Ready > error",
			flow: []order.FulfillmentStatus{
				order.FulfillmentReady,
			},
			err: errs.InvalidTranstion(order.FulfillmentPending, order.FulfillmentReady),
		},
		{
			name: "Pending > Confirmed > Ready > error",
			flow: []order.FulfillmentStatus{
				order.FulfillmentConfirmed,
				order.FulfillmentReady,
			},
			err: errs.InvalidTranstion(order.FulfillmentConfirmed, order.FulfillmentReady),
		},
		{
			name: "Pending > Confirmed > In Production > Cancelled",
			flow: []order.FulfillmentStatus{
				order.FulfillmentConfirmed,
				order.FulfillmentInProduction,
				order.FulfillmentCancelled,
			},
		},
		{
			name: "Pending > Cancelled",
			flow: []order.FulfillmentStatus{
				order.FulfillmentCancelled,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			tenant := fixtures.Tenant.Create(t.T(), nil)
			token := fixtures.Auth.StaffToken(t.T(), tenant.UserID)
			orderID := fixtures.Order.Create(t.T(), nil, token)

			for i, statusToSet := range tt.flow {
				Body := order_gateway.PatchValues{
					FulfillmentStatus: statusToSet,
				}

				appErr := fixtures.Order.Patch(t.T(), orderID, &Body, token)
				isLast := i == len(tt.flow)-1

				if isLast && tt.err != nil {
					t.Require().NotNil(appErr)
					t.Require().Equal(tt.err.Code, appErr.Code)

					return
				}

				t.Nil(appErr)
				fixtures.Order.ExpectFulfillmentStatus(t.T(), orderID, statusToSet, token)
			}
		})
	}
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
