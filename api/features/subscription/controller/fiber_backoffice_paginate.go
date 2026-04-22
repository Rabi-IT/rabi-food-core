package controller

import (
	"github.com/Rabi-IT/rabi-food-core/features/subscription"
	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"

	"github.com/gofiber/fiber/v2"
)

// Paginate godoc
// @Summary Paginate subscriptions (backoffice)
// @Description Paginate subscriptions without tenant restriction; TenantID filter is optional
// @Tags subscriptions backoffice
// @Produce json
// @Param Page query int false "Page number"
// @Param PageSize query int false "Page size"
// @Param tenantId query string false "Optional tenant ID filter"
// @Param userId query string false "Optional user ID filter"
// @Param status query string false "Optional status filter"
// @Success 200 {object} g.PaginateOutput
// @Failure 500 {string} string "Internal server error"
// @Router /backoffice/subscription/ [get].
func (c *SubscriptionController) BackofficePaginate(ctx *fiber.Ctx) error {
	filter := g.PaginateFilter{
		TenantID: ctx.Query("tenantId"),
		UserID:   ctx.Query("userId"),
		Status:   subscription.Status(ctx.Query("status")),
	}

	uctx := ctx.UserContext()

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
