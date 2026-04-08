package controller

import (
	"strconv"

	"github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"

	"github.com/gofiber/fiber/v2"
)

// Paginate godoc
// @Summary Paginate products (backoffice)
// @Description Paginate products without tenant restriction; TenantID filter is optional
// @Tags products backoffice
// @Produce json
// @Param Page query int false "Page number"
// @Param PageSize query int false "Page size"
// @Param tenantId query string false "Optional tenant ID filter"
// @Param name query string false "Product name"
// @Param categoryId query string false "Category ID"
// @Param isActive query bool false "Filter by active products"
// @Success 200 {object} product_gateway.PaginateOutput
// @Failure 500 {string} string "Internal server error"
// @Router /backoffice/product/ [get].
func (c *ProductBackofficeController) Paginate(ctx *fiber.Ctx) error {
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
	logger.GetWideEvent(uctx).Event = "backoffice-paginate-products"

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
