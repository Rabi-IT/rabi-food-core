package docs

import (
	"rabi-food-core/libs/database/gateways/category_gateway"
)

// Tipos de request/response da documentação OpenAPI para category.

type RequestCreate[T any] struct {
	Body T
}

type ResponseCreate struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type ResponseStatus struct {
	StatusCode int `json:"statusCode"`
}

type RequestPatch struct {
	ID   string `path:"id"`
	Body category_gateway.PatchValues
}

type RequestDelete struct {
	ID string `path:"id"`
}

type RequestGetByID struct {
	ID string `path:"id"`
}

type RequestPaginate struct {
	Page     int    `query:"Page" default:"0"`
	PageSize int    `query:"PageSize" default:"10"`
	TenantID string `query:"tenantId"`
	Name     string `query:"name"`
}
