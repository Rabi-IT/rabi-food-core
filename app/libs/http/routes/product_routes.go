package routes

import (
	"rabi-food-core/libs/http/controllers/product_controller"

	"github.com/gofiber/fiber/v2"
)

func Product(app *fiber.App, c *product_controller.ProductController) {
	route := app.Group("/product")
	route.Post("/", c.Create)
	route.Patch("/:id", c.Patch)
	route.Delete("/:id", c.Delete)
	route.Get("/:id", c.GetByID)
	route.Get("/", c.Paginate)
}
