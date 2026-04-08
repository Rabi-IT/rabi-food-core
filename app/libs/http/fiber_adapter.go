package http

import (
	"github.com/Rabi-IT/rabi-food-core/config"
	category_controller "github.com/Rabi-IT/rabi-food-core/features/category/controller"
	order_controller "github.com/Rabi-IT/rabi-food-core/features/order/controller"
	product_controller "github.com/Rabi-IT/rabi-food-core/features/product/controller"
	subscription_controller "github.com/Rabi-IT/rabi-food-core/features/subscription/controller"
	tenant_controller "github.com/Rabi-IT/rabi-food-core/features/tenant/controller"
	user_controller "github.com/Rabi-IT/rabi-food-core/features/user/controller"

	category_routes "github.com/Rabi-IT/rabi-food-core/features/category/routes"
	order_routes "github.com/Rabi-IT/rabi-food-core/features/order/routes"
	product_routes "github.com/Rabi-IT/rabi-food-core/features/product/routes"
	subscription_routes "github.com/Rabi-IT/rabi-food-core/features/subscription/routes"
	tenant_routes "github.com/Rabi-IT/rabi-food-core/features/tenant/routes"
	user_routes "github.com/Rabi-IT/rabi-food-core/features/user/routes"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	jwtware "github.com/gofiber/contrib/jwt"
)

type fiberAdapter struct {
	port string
	app  *fiber.App
}

func New(
	port string,
	tenantController *tenant_controller.TenantController,
	tenantBackofficeController *tenant_controller.TenantBackofficeController,
	userController *user_controller.UserController,
	userBackofficeController *user_controller.UserBackofficeController,
	productController *product_controller.ProductController,
	productBackofficeController *product_controller.ProductBackofficeController,
	categoryController *category_controller.CategoryController,
	categoryBackofficeController *category_controller.CategoryBackofficeController,
	orderController *order_controller.OrderController,
	orderBackofficeController *order_controller.OrderBackofficeController,
	subscriptionController *subscription_controller.SubscriptionController,
	subscriptionBackofficeController *subscription_controller.SubscriptionBackofficeController,
) HTTPServer {
	app := fiber.New(fiber.Config{
		Immutable:    true,
		ErrorHandler: middlewares.ErrorHandler,
	})

	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.AuthSecret)},
	})

	app.
		Use(cors.New()).
		Use(requestid.New())

	if !config.Env.IsProduction() {
		app.Static("/docs", "./libs/docs")
	}

	app.
		Use(middlewares.Logging()).
		Post("/tenant", tenantController.Create).
		Use(jwtMiddleware).
		Use(middlewares.Session)

	user_routes.User(app, userController, userBackofficeController)
	tenant_routes.Tenant(app, tenantController, tenantBackofficeController)
	product_routes.Product(app, productController, productBackofficeController)
	category_routes.Category(app, categoryController, categoryBackofficeController)
	order_routes.Order(app, orderController, orderBackofficeController)
	subscription_routes.Subscription(app, subscriptionController, subscriptionBackofficeController)

	return &fiberAdapter{
		app:  app,
		port: port,
	}
}

func (f *fiberAdapter) Start() error {
	return f.app.Listen(":" + f.port)
}

func (f *fiberAdapter) Stop() error {
	return f.app.Shutdown()
}
