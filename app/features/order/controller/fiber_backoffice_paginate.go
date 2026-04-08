package controller

import (
	"strconv"

	"github.com/Rabi-IT/rabi-food-core/features/order/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"

	"github.com/gofiber/fiber/v2"
)

// Paginate godoc
// @Summary Paginate orders (backoffice)
// @Description Paginate orders without tenant restriction; TenantID filter is optional
// @Tags orders backoffice
// @Produce json
// @Param Page query int false "Page number"
// @Param PageSize query int false "Page size"
// @Param tenantId query string false "Optional tenant ID filter"
// @Param userId query string false "User ID"
// @Param paymentStatus query string false "Payment status"
// @Param fulfillmentStatus query string false "Fulfillment status"
// @Param deliveryStatus query string false "Delivery status"
// @Param createdAtFrom query string false "Created at from (RFC3339)"
// @Param createdAtTo query string false "Created at to (RFC3339)"
// @Success 200 {object} gateway.PaginateOutput
// @Failure 500 {string} string "Internal server error"
// @Router /backoffice/order/ [get].
func (c *OrderController) BackofficePaginate(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("Page", "0"))
	if err != nil {
		return err
	}

	pageSize, err := strconv.Atoi(ctx.Query("PageSize", "10"))
	if err != nil {
		return err
	}

	filter := gateway.PaginateFilter{}
	if err = ctx.QueryParser(&filter); err != nil {
		return err
	}

	uctx := ctx.UserContext()

	paginate := database.PaginateInput{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := c.usecase.Paginate(uctx, filter, paginate)
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}
