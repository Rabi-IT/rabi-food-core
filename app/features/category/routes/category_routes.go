package routes

import (
	"github.com/Rabi-IT/rabi-food-core/domain/auth"
	"github.com/Rabi-IT/rabi-food-core/features/category/controller"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Category(app *fiber.App, c *controller.CategoryController) {
	route := app.Group("/category", middlewares.RequireTenantID)
	route.Post("/", c.Create)
	route.Patch("/:id", c.Patch)
	route.Delete("/:id", c.Delete)
	route.Get("/:id", c.GetByID)
	route.Get("/", c.Paginate)

	backoffice := app.Group("/backoffice/category", middlewares.RequireRole(auth.Backoffice))
	backoffice.Get("/", c.BackofficePaginate)
}
