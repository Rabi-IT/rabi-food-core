package controller

import (
	"github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
	"github.com/Rabi-IT/rabi-food-core/libs/database/filter"

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
// @Param isActive query int false "Filter by active status: 1=active, -1=inactive"
// @Success 200 {object} gateway.PaginateOutput
// @Failure 500 {string} string "Internal server error"
// @Router /backoffice/product/ [get].
func (c *ProductController) BackofficePaginate(ctx *fiber.Ctx) error {
	uctx := ctx.UserContext()

	f := gateway.PaginateFilter{
		TenantID:   ctx.Query("tenantId"),
		Name:       ctx.Query("name"),
		CategoryID: ctx.Query("categoryId"),
		IsActive:   filter.Bool(ctx.QueryInt("isActive", 0)),
	}

	paginate := database.PaginateInput{
		Page:     ctx.QueryInt("Page", database.DefaultPage),
		PageSize: ctx.QueryInt("PageSize", database.DefaultPageSize),
	}

	result, err := c.usecase.Paginate(uctx, f, paginate)
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}
