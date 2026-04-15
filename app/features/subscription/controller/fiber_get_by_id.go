package controller

import (
	"net/http"

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
	filter := gateway.GetByIDFilter{
		ID:       ctx.Params("id"),
		TenantID: ctx.Get("X-Tenant-ID"),
	}

	data, err := c.usecase.GetByID(ctx.UserContext(), filter)
	if err != nil {
		return err
	}

	if data == nil {
		return ctx.SendStatus(http.StatusNotFound)
	}

	return ctx.JSON(data)
}
