package routes

import (
	"rabi-food-core/features/category/controller"

	"github.com/gofiber/fiber/v2"
)

func Category(app *fiber.App, c *controller.CategoryController) {
	route := app.Group("/category")
	route.Post("/", c.Create)
	route.Patch("/:id", c.Patch)
	route.Delete("/:id", c.Delete)
	route.Get("/:id", c.GetByID)
	route.Get("/", c.Paginate)
}
