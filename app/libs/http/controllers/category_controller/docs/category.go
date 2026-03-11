package docs

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"rabi-food-core/config"
	"rabi-food-core/libs/database/gateways/category_gateway"
	"rabi-food-core/libs/http/docproxy"
	"rabi-food-core/usecases/category_case"

	"github.com/danielgtaylor/huma/v2"
)

var (
	categoryCaseForDoc *category_case.CategoryCase
)

const baseURL = "http://localhost:"

// SetCategoryCase define o usecase de category usado pela documentação (chamado pelo DI).
func SetCategoryCase(uc *category_case.CategoryCase) {
	categoryCaseForDoc = uc
}

// RegisterCategory registra as operações de documentação (OpenAPI) para category.
func RegisterCategory(api huma.API) {
	uc := categoryCaseForDoc
	docBaseURL := baseURL
	if config.AppPort != "" {
		docBaseURL += config.AppPort
	}

	humaErr := func(statusCode int, body []byte) error { return huma.NewError(statusCode, string(body)) }

	// POST /category
	createOp := docproxy.Op[*RequestCreate[category_gateway.CreateInput], *ResponseCreate]{
		Method: http.MethodPost,
		Path:   "/category",
		Body: func(input *RequestCreate[category_gateway.CreateInput]) ([]byte, error) {
			return json.Marshal(input.Body)
		},
		BuildOutput: func(statusCode int, body []byte) (*ResponseCreate, error) {
			return &ResponseCreate{StatusCode: statusCode, Body: string(body)}, nil
		},
		BuildError: humaErr,
	}
	huma.Register(api, huma.Operation{
		OperationID: "create-category",
		Method:      "POST",
		Path:        "/category",
		Summary:     "Criar categoria",
		Tags:        []string{"Category"},
	}, docproxy.Wrap(docBaseURL, createOp, createDoc(uc)))

	// PATCH /category/{id}
	patchOp := docproxy.Op[*RequestPatch, *ResponseStatus]{
		Method:   http.MethodPatch,
		PathFunc: func(input *RequestPatch) string { return "/category/" + input.ID },
		Body:     func(input *RequestPatch) ([]byte, error) { return json.Marshal(input.Body) },
		BuildOutput: func(statusCode int, body []byte) (*ResponseStatus, error) {
			return &ResponseStatus{StatusCode: statusCode}, nil
		},
		BuildError: humaErr,
	}
	huma.Register(api, huma.Operation{
		OperationID: "patch-category",
		Method:      "PATCH",
		Path:        "/category/{id}",
		Summary:     "Atualizar categoria",
		Tags:        []string{"Category"},
	}, docproxy.Wrap(docBaseURL, patchOp, adaptPatch(uc)))

	// DELETE /category/{id}
	deleteOp := docproxy.Op[*RequestDelete, *ResponseStatus]{
		Method:   http.MethodDelete,
		PathFunc: func(input *RequestDelete) string { return "/category/" + input.ID },
		BuildOutput: func(statusCode int, body []byte) (*ResponseStatus, error) {
			return &ResponseStatus{StatusCode: statusCode}, nil
		},
		BuildError: humaErr,
	}
	huma.Register(api, huma.Operation{
		OperationID: "delete-category",
		Method:      "DELETE",
		Path:        "/category/{id}",
		Summary:     "Excluir categoria",
		Tags:        []string{"Category"},
	}, docproxy.Wrap(docBaseURL, deleteOp, adaptDelete(uc)))

	// GET /category/{id}
	getByIDOp := docproxy.Op[*RequestGetByID, *category_gateway.GetByIDOutput]{
		Method:   http.MethodGet,
		PathFunc: func(input *RequestGetByID) string { return "/category/" + input.ID },
		BuildOutput: func(statusCode int, body []byte) (*category_gateway.GetByIDOutput, error) {
			if statusCode != http.StatusOK {
				return nil, huma.NewError(statusCode, string(body))
			}
			var out category_gateway.GetByIDOutput
			if err := json.Unmarshal(body, &out); err != nil {
				return nil, err
			}
			return &out, nil
		},
		BuildError: humaErr,
	}
	huma.Register(api, huma.Operation{
		OperationID: "get-category-by-id",
		Method:      "GET",
		Path:        "/category/{id}",
		Summary:     "Buscar categoria por ID",
		Tags:        []string{"Category"},
	}, docproxy.Wrap(docBaseURL, getByIDOp, adaptGetByID(uc)))

	// GET /category (paginate)
	paginatePathWithQuery := func(input *RequestPaginate) string {
		path := "/category?"
		path += "Page=" + strconv.Itoa(input.Page) + "&PageSize=" + strconv.Itoa(input.PageSize)
		if input.TenantID != "" {
			path += "&tenantId=" + url.QueryEscape(input.TenantID)
		}
		if input.Name != "" {
			path += "&name=" + url.QueryEscape(input.Name)
		}
		return path
	}
	paginateOp := docproxy.Op[*RequestPaginate, *category_gateway.PaginateOutput]{
		Method:   http.MethodGet,
		PathFunc: paginatePathWithQuery,
		BuildOutput: func(statusCode int, body []byte) (*category_gateway.PaginateOutput, error) {
			if statusCode != http.StatusOK {
				return nil, huma.NewError(statusCode, string(body))
			}
			var out category_gateway.PaginateOutput
			if err := json.Unmarshal(body, &out); err != nil {
				return nil, err
			}
			return &out, nil
		},
		BuildError: humaErr,
	}
	huma.Register(api, huma.Operation{
		OperationID: "list-categories",
		Method:      "GET",
		Path:        "/category",
		Summary:     "Listar categorias (paginado)",
		Tags:        []string{"Category"},
	}, docproxy.Wrap(docBaseURL, paginateOp, adaptPaginate(uc)))
}
