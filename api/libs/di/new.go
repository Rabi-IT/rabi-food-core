package di

import (
	"errors"

	"github.com/Rabi-IT/rabi-food-core/config"
	auth_controller "github.com/Rabi-IT/rabi-food-core/features/auth/controller"
	auth_gateway "github.com/Rabi-IT/rabi-food-core/features/auth/gateway"
	auth_case "github.com/Rabi-IT/rabi-food-core/features/auth/usecases"
	category_controller "github.com/Rabi-IT/rabi-food-core/features/category/controller"
	category_gateway "github.com/Rabi-IT/rabi-food-core/features/category/gateway"
	category_case "github.com/Rabi-IT/rabi-food-core/features/category/usecases"
	order_controller "github.com/Rabi-IT/rabi-food-core/features/order/controller"
	order_gateway "github.com/Rabi-IT/rabi-food-core/features/order/gateway"
	"github.com/Rabi-IT/rabi-food-core/features/order/usecases"
	product_controller "github.com/Rabi-IT/rabi-food-core/features/product/controller"
	product_gateway "github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	product_case "github.com/Rabi-IT/rabi-food-core/features/product/usecases"
	payment_controller "github.com/Rabi-IT/rabi-food-core/features/payment/controller"
	payment_gateway "github.com/Rabi-IT/rabi-food-core/features/payment/gateway"
	payment_case "github.com/Rabi-IT/rabi-food-core/features/payment/usecases"
	subscription_controller "github.com/Rabi-IT/rabi-food-core/features/subscription/controller"
	subscription_gateway "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	subscription_case "github.com/Rabi-IT/rabi-food-core/features/subscription/usecases"
	tenant_controller "github.com/Rabi-IT/rabi-food-core/features/tenant/controller"
	tenant_gateway "github.com/Rabi-IT/rabi-food-core/features/tenant/gateway"
	tenant_case "github.com/Rabi-IT/rabi-food-core/features/tenant/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
	"github.com/Rabi-IT/rabi-food-core/libs/http"

	"github.com/samber/do"
)

var (
	ErrHTTPPortNotConfigured = errors.New("HTTP port is not configured")
)

func newInjector(dbConfig *config.DatabaseConfig) *do.Injector {
	injector := do.New()

	// Database
	do.Provide(injector, func(_ *do.Injector) (*database.PgxAdapter, error) {
		db, ok := database.NewPgx(dbConfig).(*database.PgxAdapter)
		if !ok {
			//nolint:err113
			return nil, errors.New("failed to cast database adapter")
		}

		return db, nil
	})

	// Database as interface
	do.Provide(injector, func(i *do.Injector) (database.Database, error) {
		return do.MustInvoke[*database.PgxAdapter](i), nil
	})

	// HTTP Server
	do.Provide(injector, func(i *do.Injector) (http.HTTPServer, error) {
		authController := do.MustInvoke[*auth_controller.AuthController](i)
		tenantController := do.MustInvoke[*tenant_controller.TenantController](i)
		productController := do.MustInvoke[*product_controller.ProductController](i)
		categoryController := do.MustInvoke[*category_controller.CategoryController](i)
		orderController := do.MustInvoke[*order_controller.OrderController](i)
		subscriptionController := do.MustInvoke[*subscription_controller.SubscriptionController](i)
		paymentController := do.MustInvoke[*payment_controller.PaymentController](i)

		if config.ApiPort == "" {
			return nil, ErrHTTPPortNotConfigured
		}

		return http.New(
			config.ApiPort,
			authController,
			tenantController,
			productController,
			categoryController,
			orderController,
			subscriptionController,
			paymentController,
		), nil
	})

	// Auth dependencies
	do.Provide(injector, func(i *do.Injector) (auth_gateway.AuthGateway, error) {
		db := do.MustInvoke[*database.PgxAdapter](i)

		return auth_gateway.NewGoTrue(config.GoTrueURL, config.GoTrueServiceKey, db), nil
	})

	// Tenant dependencies (registered before AuthCase — AuthCase depends on TenantCase)
	do.Provide(injector, func(i *do.Injector) (tenant_gateway.TenantGateway, error) {
		db := do.MustInvoke[*database.PgxAdapter](i)

		return &tenant_gateway.PgxTenantGatewayAdapter{DB: db}, nil
	})

	do.Provide(injector, func(i *do.Injector) (*tenant_case.TenantCase, error) {
		gw := do.MustInvoke[tenant_gateway.TenantGateway](i)

		return tenant_case.New(gw), nil
	})

	do.Provide(injector, func(i *do.Injector) (*auth_case.AuthCase, error) {
		gw := do.MustInvoke[auth_gateway.AuthGateway](i)
		tc := do.MustInvoke[*tenant_case.TenantCase](i)

		return auth_case.New(gw, tc), nil
	})

	do.Provide(injector, func(i *do.Injector) (*auth_controller.AuthController, error) {
		c := do.MustInvoke[*auth_case.AuthCase](i)

		return auth_controller.New(c), nil
	})

	do.Provide(injector, func(i *do.Injector) (*tenant_controller.TenantController, error) {
		c := do.MustInvoke[*tenant_case.TenantCase](i)

		return tenant_controller.New(c), nil
	})

	// Product dependencies
	do.Provide(injector, func(i *do.Injector) (*product_case.ProductCase, error) {
		gw := do.MustInvoke[product_gateway.ProductGateway](i)

		return product_case.New(gw), nil
	})

	do.Provide(injector, func(i *do.Injector) (product_gateway.ProductGateway, error) {
		db := do.MustInvoke[*database.PgxAdapter](i)

		return &product_gateway.PgxProductGatewayAdapter{DB: db}, nil
	})

	do.Provide(injector, func(i *do.Injector) (*product_controller.ProductController, error) {
		c := do.MustInvoke[*product_case.ProductCase](i)

		return product_controller.New(c), nil
	})

	// Category dependencies
	do.Provide(injector, func(i *do.Injector) (*category_controller.CategoryController, error) {
		c := do.MustInvoke[*category_case.CategoryCase](i)

		return category_controller.New(c), nil
	})

	do.Provide(injector, func(i *do.Injector) (*category_case.CategoryCase, error) {
		gw := do.MustInvoke[category_gateway.CategoryGateway](i)

		return category_case.New(gw), nil
	})

	do.Provide(injector, func(i *do.Injector) (category_gateway.CategoryGateway, error) {
		db := do.MustInvoke[*database.PgxAdapter](i)

		return &category_gateway.PgxCategoryGatewayAdapter{DB: db}, nil
	})

	// Order dependencies
	do.Provide(injector, func(i *do.Injector) (*order_controller.OrderController, error) {
		c := do.MustInvoke[*usecases.OrderCase](i)

		return order_controller.New(c), nil
	})

	do.Provide(injector, func(i *do.Injector) (order_gateway.OrderGateway, error) {
		db := do.MustInvoke[*database.PgxAdapter](i)

		return &order_gateway.PgxOrderGatewayAdapter{DB: db}, nil
	})

	do.Provide(injector, func(i *do.Injector) (*usecases.OrderCase, error) {
		gw := do.MustInvoke[order_gateway.OrderGateway](i)
		productCase := do.MustInvoke[*product_case.ProductCase](i)

		return usecases.New(gw, productCase), nil
	})

	// Payment dependencies
	do.Provide(injector, func(i *do.Injector) (payment_gateway.PaymentGateway, error) {
		auth := do.MustInvoke[auth_gateway.AuthGateway](i)
		return payment_gateway.NewStripe(auth), nil
	})

	do.Provide(injector, func(i *do.Injector) (*payment_case.PaymentCase, error) {
		payment := do.MustInvoke[payment_gateway.PaymentGateway](i)
		subCase := do.MustInvoke[*subscription_case.SubscriptionCase](i)

		return payment_case.New(payment, subCase), nil
	})

	do.Provide(injector, func(i *do.Injector) (*payment_controller.PaymentController, error) {
		c := do.MustInvoke[*payment_case.PaymentCase](i)

		return payment_controller.New(c), nil
	})

	// Subscription dependencies
	do.Provide(injector, func(i *do.Injector) (*subscription_controller.SubscriptionController, error) {
		c := do.MustInvoke[*subscription_case.SubscriptionCase](i)

		return subscription_controller.New(c), nil
	})

	do.Provide(injector, func(i *do.Injector) (subscription_gateway.SubscriptionGateway, error) {
		db := do.MustInvoke[*database.PgxAdapter](i)

		return &subscription_gateway.PgxSubscriptionGatewayAdapter{DB: db}, nil
	})

	do.Provide(injector, func(i *do.Injector) (*subscription_case.SubscriptionCase, error) {
		gw := do.MustInvoke[subscription_gateway.SubscriptionGateway](i)
		productCase := do.MustInvoke[*product_case.ProductCase](i)
		payment := do.MustInvoke[payment_gateway.PaymentGateway](i)

		return subscription_case.New(gw, productCase, payment), nil
	})

	return injector
}

func NewTest() *do.Injector {
	return newInjector(config.TestDatabase)
}

func Injector() *do.Injector {
	if config.Env.IsProduction() {
		return newInjector(config.ProductionDatabase)
	}

	return NewTest()
}
