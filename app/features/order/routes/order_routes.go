package routes

import (
	"rabi-food-core/features/order/controller"

	"github.com/gofiber/fiber/v2"
)

func Order(app *fiber.App, c *controller.OrderController) {
	route := app.Group("/order")
	route.Post("/", c.Create)
	route.Delete("/:id", c.Delete)
	route.Get("/:id", c.GetByID)
	route.Get("/", c.Paginate)
	route.Post("/:id/payments/confirm", c.ConfirmPayment)
}
