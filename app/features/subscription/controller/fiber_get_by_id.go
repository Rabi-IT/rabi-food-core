package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	gateway "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"

	"github.com/gofiber/fiber/v2"
)

// keep import for swagger documentation.
var _ *gateway.GetByIDOutput

// GetByID godoc
// @Summary Get subscription by ID
// @Description Get a subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} gateway.GetByIDOutput
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /subscription/{id} [get].
func (c *SubscriptionController) GetByID(ctx *fiber.Ctx) error {
	uctx := ctx.UserContext()
	session := app_context.GetSession(uctx)

	filter := gateway.GetByIDFilter{ID: ctx.Params("id")}
	switch {
	case session.Role.IsUser():
		filter.UserID = session.UserID
	case session.Role.IsTenant():
		filter.TenantID = session.TenantID
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
