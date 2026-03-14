package order_controller

import (
	"rabi-food-core/libs/database"
	"rabi-food-core/libs/database/gateways/order_gateway"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Paginate godoc
// @Summary Paginate orders
// @Description Paginate orders using query filters
// @Tags orders
// @Produce json
// @Param Page query int false "Page number"
// @Param PageSize query int false "Page size"
// @Param userId query string false "User ID"
// @Param tenantId query string false "Tenant ID"
// @Param paymentStatus query string false "Payment status"
// @Param fulfillmentStatus query string false "Fulfillment status"
// @Param deliveryStatus query string false "Delivery status"
// @Param createdAtFrom query string false "Created at from (RFC3339)"
// @Param createdAtTo query string false "Created at to (RFC3339)"
// @Success 200 {object} order_gateway.PaginateOutput
// @Failure 500 {string} string "Internal server error"
// @Router /order/ [get].
func (c *OrderController) Paginate(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("Page", "0"))
	if err != nil {
		return err
	}

	pageSize, err := strconv.Atoi(ctx.Query("PageSize", "10"))
	if err != nil {
		return err
	}

	filter := order_gateway.PaginateFilter{}
	err = ctx.QueryParser(&filter)
	if err != nil {
		return err
	}

	paginate := database.PaginateInput{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := c.usecase.Paginate(ctx.UserContext(), filter, paginate)

	if err != nil {
		return err
	}

	return ctx.JSON(result)
}
