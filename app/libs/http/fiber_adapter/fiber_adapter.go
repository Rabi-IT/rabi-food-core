package fiber_adapter

import (
	"rabi-food-core/config"
	"rabi-food-core/libs/docs"
	"rabi-food-core/libs/http"
	"rabi-food-core/libs/http/controllers/category_controller"
	"rabi-food-core/libs/http/controllers/order_controller"
	"rabi-food-core/libs/http/controllers/product_controller"
	"rabi-food-core/libs/http/controllers/subscription_controller"
	"rabi-food-core/libs/http/controllers/tenant_controller"
	"rabi-food-core/libs/http/controllers/user_controller"
	"rabi-food-core/libs/http/fiber_adapter/middlewares"
	"rabi-food-core/libs/http/routes"

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
) http.HTTPServer {
	app := fiber.New(fiber.Config{
		Immutable:    true,
		ErrorHandler: middlewares.ErrorHandler,
	})

	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(config.AuthSecret)},
		ErrorHandler: middlewares.ErrorHandler,
	})

	app.
		Use(cors.New()).
		Use(requestid.New())

	if !config.Env.IsProduction() {
		app.
			Static("/docs", "./libs/docs").
			Get("/docs", docs.ServeUI)
	}

	app.Post("/tenant", tenantController.Create).
		Use(jwtMiddleware).
		Use(middlewares.Session).
		Use(middlewares.Logging())

	routes.User(app, userController)
	routes.Tenant(app, tenantController)
	routes.Product(app, productController)
	routes.Category(app, categoryController)
	routes.Order(app, orderController)
	routes.Subscription(app, subscriptionController)

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
