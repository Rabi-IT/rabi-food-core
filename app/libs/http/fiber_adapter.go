package http

import (
	"rabi-food-core/config"
	category_controller "rabi-food-core/features/category/controller"
	order_controller "rabi-food-core/features/order/controller"
	product_controller "rabi-food-core/features/product/controller"
	subscription_controller "rabi-food-core/features/subscription/controller"
	tenant_controller "rabi-food-core/features/tenant/controller"
	user_controller "rabi-food-core/features/user/controller"

	category_routes "rabi-food-core/features/category/routes"
	order_routes "rabi-food-core/features/order/routes"
	product_routes "rabi-food-core/features/product/routes"
	subscription_routes "rabi-food-core/features/subscription/routes"
	tenant_routes "rabi-food-core/features/tenant/routes"
	user_routes "rabi-food-core/features/user/routes"
	"rabi-food-core/libs/http/middlewares"

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
	userController *user_controller.UserController,
	productController *product_controller.ProductController,
	categoryController *category_controller.CategoryController,
	orderController *order_controller.OrderController,
	subscriptionController *subscription_controller.SubscriptionController,
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

	app.Post("/tenant", tenantController.Create).
		Use(jwtMiddleware).
		Use(middlewares.Session).
		Use(middlewares.Logging())

	user_routes.User(app, userController)
	tenant_routes.Tenant(app, tenantController)
	product_routes.Product(app, productController)
	category_routes.Category(app, categoryController)
	order_routes.Order(app, orderController)
	subscription_routes.Subscription(app, subscriptionController)

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
