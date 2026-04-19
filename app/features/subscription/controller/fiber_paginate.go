package controller

import (
	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"

	"github.com/gofiber/fiber/v2"
)

// Paginate godoc
// @Summary Paginate subscriptions
// @Description Paginate subscriptions using query filters
// @Tags subscriptions
// @Produce json
// @Param Page query int false "Page number"
// @Param PageSize query int false "Page size"
// @Success 200 {object} g.PaginateOutput
// @Failure 500 {string} string "Internal server error"
// @Router /subscription/ [get].
func (c *SubscriptionController) Paginate(ctx *fiber.Ctx) error {
	uctx := ctx.UserContext()
	tenantID := ctx.Get("X-Tenant-ID")

	filter := g.PaginateFilter{
		TenantID: tenantID,
	}

	paginate := database.PaginateInput{
		Page:     ctx.QueryInt("Page", database.DefaultPage),
		PageSize: ctx.QueryInt("PageSize", database.DefaultPageSize),
	}

	result, err := c.usecase.Paginate(uctx, filter, paginate)
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}
