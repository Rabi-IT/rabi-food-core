package fiber_adapter

import (
	"rabi-food-core/config"
	"rabi-food-core/libs/http"
	"rabi-food-core/libs/http/controllers/category_controller"
	"rabi-food-core/libs/http/controllers/category_controller/docs"
	"rabi-food-core/libs/http/controllers/order_controller"
	"rabi-food-core/libs/http/controllers/product_controller"
	"rabi-food-core/libs/http/controllers/subscription_controller"
	"rabi-food-core/libs/http/controllers/tenant_controller"
	"rabi-food-core/libs/http/controllers/user_controller"
	"rabi-food-core/libs/http/fiber_adapter/middlewares"
	ui "rabi-food-core/libs/http/fiber_adapter/scalar"
	"rabi-food-core/libs/http/routes"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
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

	if !config.Env.IsProduction() {
		docApi := humafiber.New(app, huma.Config{
			OpenAPI: &huma.OpenAPI{
				OpenAPI: "3.1.0",
				Info: &huma.Info{
					Title:   "Products API",
					Version: "1.0.0",
				},
			},
			OpenAPIPath: "/openapi",
		})

		scalarUI := ui.NewScalarUI(docApi)
		scalarUI.RegisterRoutes(app)

		defer docs.RegisterCategory(docApi)
	}

	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.AuthSecret)},
	})

	requestIDMiddleware := requestid.New()

	app.
		Use(cors.New()).
		Use(requestIDMiddleware).
		Post("/tenant", tenantController.Create).
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
