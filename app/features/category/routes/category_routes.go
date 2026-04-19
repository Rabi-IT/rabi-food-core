package routes

import (
	"github.com/Rabi-IT/rabi-food-core/features/category/controller"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

	"github.com/gofiber/fiber/v2"
)

// CategoryProtected registers category routes that require a valid JWT.
func CategoryProtected(app *fiber.App, c *controller.CategoryController) {
	route := app.Group("/category")

	route.Get("/", middlewares.RequireTenantID, c.Paginate)
	route.Get("/:id", c.GetByID)

	route.Post("/", middlewares.RequireTenant, c.Create)
	route.Patch("/:id", middlewares.RequireTenant, c.Patch)
	route.Delete("/:id", middlewares.RequireTenant, c.Delete)

	backoffice := app.Group("/backoffice/category", middlewares.RequireBackoffice)
	backoffice.Get("/", c.BackofficePaginate)
}
