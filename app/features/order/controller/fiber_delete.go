package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/features/order/gateway"
	parser "github.com/Rabi-IT/rabi-food-core/libs/http/parser"

	"github.com/gofiber/fiber/v2"
)

// Delete godoc
// @Summary Delete order
// @Description Delete an order by ID
// @Tags orders
// @Accept json
// @Produce text/plain
// @Param id path string true "Order ID"
// @Param order body gateway.DeleteFilter false "Optional delete filter data"
// @Success 204 {string} string "No content"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /order/{id} [delete].
func (c *OrderController) Delete(ctx *fiber.Ctx) error {
	filter := gateway.DeleteFilter{}

	err := parser.ParseBody(ctx, &filter)
	if err != nil {
		return ctx.JSON(err)
	}

	filter.ID = ctx.Params("id")
	deleted, err := c.usecase.Delete(ctx.UserContext(), filter)

	if err != nil {
		return err
	}

	if deleted {
		return ctx.SendStatus(http.StatusNoContent)
	}

	return ctx.SendStatus(http.StatusNotFound)
}
