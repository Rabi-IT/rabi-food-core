package routes

import (
	"github.com/Rabi-IT/rabi-food-core/features/product/controller"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

	"github.com/gofiber/fiber/v2"
)

// ProductProtected registers product routes that require a valid JWT.
func ProductProtected(app *fiber.App, c *controller.ProductController) {
	route := app.Group("/product")
	route.Get("/", middlewares.RequireTenantID, c.Paginate)
	route.Get("/:id", c.GetByID)

	route.Post("/", middlewares.RequireTenant, c.Create)
	route.Patch("/:id", middlewares.RequireTenant, c.Patch)
	route.Delete("/:id", middlewares.RequireTenant, c.Delete)

	backoffice := app.Group("/backoffice/product", middlewares.RequireBackoffice)
	backoffice.Get("/", c.BackofficePaginate)
}
