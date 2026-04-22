package routes

import (
	"github.com/Rabi-IT/rabi-food-core/features/order/controller"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

	"github.com/gofiber/fiber/v2"
)

// OrderProtected registers order routes that require a valid JWT.
func OrderProtected(app *fiber.App, c *controller.OrderController) {
	route := app.Group("/order")
	route.Post("/", middlewares.RequireTenantID, c.Create)
	route.Delete("/:id", c.Delete)
	route.Get("/:id", c.GetByID)
	route.Get("/", c.Paginate)
	route.Post("/:id/payments/confirm", c.ConfirmPayment)

	backoffice := app.Group("/backoffice/order", middlewares.RequireBackoffice)
	backoffice.Get("/", c.BackofficePaginate)
}
