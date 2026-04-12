package http

import (
	"github.com/Rabi-IT/rabi-food-core/config"
	auth_controller "github.com/Rabi-IT/rabi-food-core/features/auth/controller"
	category_controller "github.com/Rabi-IT/rabi-food-core/features/category/controller"
	order_controller "github.com/Rabi-IT/rabi-food-core/features/order/controller"
	product_controller "github.com/Rabi-IT/rabi-food-core/features/product/controller"
	subscription_controller "github.com/Rabi-IT/rabi-food-core/features/subscription/controller"
	tenant_controller "github.com/Rabi-IT/rabi-food-core/features/tenant/controller"

	auth_routes "github.com/Rabi-IT/rabi-food-core/features/auth/routes"
	category_routes "github.com/Rabi-IT/rabi-food-core/features/category/routes"
	order_routes "github.com/Rabi-IT/rabi-food-core/features/order/routes"
	product_routes "github.com/Rabi-IT/rabi-food-core/features/product/routes"
	subscription_routes "github.com/Rabi-IT/rabi-food-core/features/subscription/routes"
	tenant_routes "github.com/Rabi-IT/rabi-food-core/features/tenant/routes"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

	fiberprometheus "github.com/ansrivas/fiberprometheus/v2"
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
	authController *auth_controller.AuthController,
	tenantController *tenant_controller.TenantController,
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

	prom := fiberprometheus.New("rabi-food-core")
	prom.RegisterAt(app, "/metrics")

	app.
		Use(cors.New()).
		Use(requestid.New()).
		Use(prom.Middleware)

	if !config.Env.IsProduction() {
		app.Static("/docs", "./libs/docs")
	}

	app.Use(middlewares.Logging())

	// Public routes — no JWT required
	auth_routes.Auth(app, authController)
	app.Post("/tenant", tenantController.Create)

	// Protected routes — JWT required
	app.
		Use(jwtMiddleware).
		Use(middlewares.Session)

	auth_routes.AuthProtected(app, authController)
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
