package routes

import (
	"github.com/Rabi-IT/rabi-food-core/domain/auth"
	"github.com/Rabi-IT/rabi-food-core/features/product/controller"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Product(app *fiber.App, c *controller.ProductController) {
	route := app.Group("/product")
	route.Post("/", c.Create)
	route.Patch("/:id", c.Patch)
	route.Delete("/:id", c.Delete)
	route.Get("/:id", c.GetByID)
	route.Get("/", c.Paginate)

	backoffice := app.Group("/backoffice/product", middlewares.RequireRole(auth.Backoffice))
	backoffice.Get("/", c.BackofficePaginate)
}
