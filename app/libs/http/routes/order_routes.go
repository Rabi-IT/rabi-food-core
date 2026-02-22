package routes

import (
	"rabi-food-core/libs/http/controllers/order_controller"

	"github.com/gofiber/fiber/v2"
)

func Order(app *fiber.App, c *order_controller.OrderController) {
	route := app.Group("/order")
	route.Post("/", c.Create)
	route.Delete("/:id", c.Delete)
	route.Get("/:id", c.GetByID)
	route.Get("/", c.Paginate)
	route.Post("/:id/payments/confirm", c.ConfirmPayment)
}
