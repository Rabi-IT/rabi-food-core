package controller

import (
	g "github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	"github.com/Rabi-IT/rabi-food-core/features/product/usecases"

	"github.com/gofiber/fiber/v2"
)

type catalogCategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type catalogProductResponse struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	Photo         string           `json:"photo"`
	CategoryID    string           `json:"categoryId"`
	Unit          string           `json:"unit"`
	Price         uint             `json:"price"`
	DiscountRules []g.DiscountRule `json:"discountRules"`
}

type catalogResponse struct {
	Products   []catalogProductResponse  `json:"products"`
	Categories []catalogCategoryResponse `json:"categories"`
}

func (c *ProductController) Catalog(ctx *fiber.Ctx) error {
	uctx := ctx.UserContext()

	products, err := c.usecase.Catalog(uctx, usecases.CatalogInput{
		TenantID: ctx.Get("X-Tenant-ID"),
	})
	if err != nil {
		return err
	}

	seen := make(map[string]struct{}, len(products))
	categories := make([]catalogCategoryResponse, 0, len(products))
	catalogProducts := make([]catalogProductResponse, len(products))

	for i, p := range products {
		catalogProducts[i] = catalogProductResponse{
			ID:            p.ID,
			Name:          p.Name,
			Description:   p.Description,
			Photo:         p.Photo,
			CategoryID:    p.CategoryID,
			Unit:          p.Unit,
			Price:         p.Price,
			DiscountRules: p.DiscountRules,
		}
		if _, ok := seen[p.CategoryID]; ok {
			continue
		}
		seen[p.CategoryID] = struct{}{}
		categories = append(categories, catalogCategoryResponse{ID: p.CategoryID, Name: p.CategoryName})
	}

	return ctx.JSON(catalogResponse{Products: catalogProducts, Categories: categories})
}
