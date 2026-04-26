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
	api  *fiber.App
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
	api := fiber.New(fiber.Config{
		Immutable:    true,
		ErrorHandler: middlewares.ErrorHandler,
	})

	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.AuthSecret)},
	})

	prom := fiberprometheus.New("rabi-food-core")
	prom.RegisterAt(api, "/metrics")

	api.
		Use(cors.New()).
		Use(requestid.New()).
		Use(prom.Middleware).
		Use(middlewares.Logging()).
		Get("/health", func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

	if !config.Env.IsProduction() {
		api.Static("/docs", "./libs/docs")
	}

	api.Static("/templates", "./libs/templates")

	// Public routes — no JWT required
	auth_routes.Auth(api, authController)
	tenant_routes.Tenant(api, tenantController)
	product_routes.Product(api, productController)
	category_routes.Category(api, categoryController)
	subscription_routes.Subscription(api, subscriptionController)

	// Protected routes — JWT required
	api.
		Use(jwtMiddleware).
		Use(middlewares.Session)

	auth_routes.AuthProtected(api, authController)
	tenant_routes.TenantProtected(api, tenantController)
	product_routes.ProductProtected(api, productController)
	category_routes.CategoryProtected(api, categoryController)
	order_routes.OrderProtected(api, orderController)
	subscription_routes.SubscriptionProtected(api, subscriptionController)

	return &fiberAdapter{
		api:  api,
		port: port,
	}
}

func (f *fiberAdapter) Start() error {
	return f.api.Listen(":" + f.port)
}

func (f *fiberAdapter) Stop() error {
	return f.api.Shutdown()
}
