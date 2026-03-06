package ui

import (
	"encoding/json"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
)

type ScalarUI struct {
	api huma.API
}

func NewScalarUI(api huma.API) *ScalarUI {
	return &ScalarUI{api: api}
}

func (s *ScalarUI) RegisterRoutes(router fiber.Router) {
	docs := router.Group("/docs")
	docs.Get("/openapi.json", s.serveSchema)
	docs.Get("/", s.serveUI)
}

func (s *ScalarUI) serveSchema(c *fiber.Ctx) error {
	schema, err := json.MarshalIndent(s.api.OpenAPI(), "", "  ")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("failed to generate schema")
	}

	c.Set("Content-Type", "application/json")
	return c.Send(schema)
}

func (s *ScalarUI) serveUI(c *fiber.Ctx) error {
	html := `
	<!doctype html>
	<html>
	<head>
		<title>Scalar API Reference</title>
		<meta charset="utf-8" />
		<meta
		name="viewport"
		content="width=device-width, initial-scale=1" />
	</head>
	<body>
		<div id="app"></div>
		<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
		<script>
		Scalar.createApiReference('#app', {
			url: '/docs/openapi.json',
		})
		</script>
	</body>
	</html>`

	c.Set("Content-Type", "text/html")
	c.Set("Cache-Control", "no-store, no-cache, must-revalidate")
	c.Set("Pragma", "no-cache")
	return c.SendString(html)
}
