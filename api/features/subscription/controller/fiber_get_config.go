package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	gateway "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

// keep import for swagger documentation.
var _ *gateway.GetConfigOutput

// GetConfig godoc
// @Summary Get subscription config
// @Description Get subscription configuration for the current tenant
// @Tags subscriptions
// @Produce json
// @Success 200 {object} gateway.GetConfigOutput
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /subscription/config [get].
func (c *SubscriptionController) GetConfig(ctx *fiber.Ctx) error {
	tenantID := ctx.Get("X-Tenant-ID")

	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "get-subscription-config"

	data, err := c.usecase.GetConfig(uctx, tenantID)
	if err != nil {
		return err
	}

	if data == nil {
		return ctx.SendStatus(http.StatusNotFound)
	}

	return ctx.JSON(data)
}
