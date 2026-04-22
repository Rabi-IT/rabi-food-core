package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/features/order/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"

	"github.com/gofiber/fiber/v2"
)

// keep import for swagger documentation.
var _ *gateway.GetByIDOutput

// GetByID godoc
// @Summary Get order by ID
// @Description Get an order by ID
// @Tags orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} gateway.GetByIDOutput
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /order/{id} [get].
func (c *OrderController) GetByID(ctx *fiber.Ctx) error {
	uctx := ctx.UserContext()
	session := app_context.GetSession(uctx)

	filter := gateway.GetByIDFilter{ID: ctx.Params("id")}
	if session.IsTenant() {
		filter.TenantID = session.TenantID
	} else if session.Role.IsUser() {
		filter.UserID = session.UserID
	} else {
		return errs.ErrForbidden
	}

	data, err := c.usecase.GetByID(uctx, filter)

	if err != nil {
		return err
	}

	if data == nil {
		return ctx.SendStatus(http.StatusNotFound)
	}

	return ctx.JSON(data)
}
