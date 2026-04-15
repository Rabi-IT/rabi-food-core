package controller

import (
	"strconv"

	"github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"

	"github.com/gofiber/fiber/v2"
)

// Paginate godoc
// @Summary Paginate products
// @Description Paginate products using query filters
// @Tags products
// @Produce json
// @Param Page query int false "Page number"
// @Param PageSize query int false "Page size"
// @Param tenantId query string false "Tenant ID"
// @Param name query string false "Product name"
// @Param categoryId query string false "Category ID"
// @Param isActive query bool false "Filter by active products"
// @Success 200 {object} gateway.PaginateOutput
// @Failure 500 {string} string "Internal server error"
// @Router /product/ [get].
func (c *ProductController) Paginate(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("Page", "0"))
	if err != nil {
		return err
	}

	pageSize, err := strconv.Atoi(ctx.Query("PageSize", "10"))
	if err != nil {
		return err
	}

	filter := gateway.PaginateFilter{}
	err = ctx.QueryParser(&filter)
	if err != nil {
		return err
	}

	uctx := ctx.UserContext()
	tenantID := ctx.Get("X-Tenant-ID")
	filter.TenantID = &tenantID

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
