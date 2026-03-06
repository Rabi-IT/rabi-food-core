package category_controller

import (
	"context"
	"net/http"
	"rabi-food-core/usecases/category_case"

	"github.com/danielgtaylor/huma/v2"
)

type CategoryController struct {
	usecase *category_case.CategoryCase
}

func New(usecase *category_case.CategoryCase) *CategoryController {
	return &CategoryController{usecase}
}

type RequestCreate[T any] struct {
	Body T
}

type ResponseCreate struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

func AdaptCreate[T any, R string](
	handler func(ctx context.Context, input T) (R, error),
) func(ctx context.Context, input *RequestCreate[T]) (*ResponseCreate, error) {
	return func(ctx context.Context, input *RequestCreate[T]) (*ResponseCreate, error) {
		result, err := handler(ctx, input.Body)
		if err != nil {
			return nil, err
		}
		return &ResponseCreate{
			StatusCode: http.StatusCreated,
			Body:       string(result),
		}, nil
	}
}

func (ctrl *CategoryController) AddDocRoute(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "create-category",
		Method:      "POST",
		Path:        "/category",
		Summary:     "Criar categoria",
		Tags:        []string{"Category"},
	}, AdaptCreate(ctrl.usecase.Create)) // TODO: Fazer chamar o endpoint, e não o usecase.
}
