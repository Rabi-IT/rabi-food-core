package docs

import (
	"context"
	"net/http"

	"rabi-food-core/libs/database"
	"rabi-food-core/libs/database/gateways/category_gateway"
	"rabi-food-core/usecases/category_case"

	"github.com/danielgtaylor/huma/v2"
)

// Handlers de documentação: recebem o usecase e retornam as funções usadas pelo docproxy.

func adaptCreate[T any, R string](
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

func createDoc(uc *category_case.CategoryCase) func(context.Context, *RequestCreate[category_gateway.CreateInput]) (*ResponseCreate, error) {
	return adaptCreate(uc.Create)
}

func adaptPatch(uc *category_case.CategoryCase) func(context.Context, *RequestPatch) (*ResponseStatus, error) {
	return func(ctx context.Context, input *RequestPatch) (*ResponseStatus, error) {
		updated, err := uc.Patch(ctx, category_gateway.PatchFilter{ID: input.ID}, input.Body)
		if err != nil {
			return nil, err
		}
		if !updated {
			return nil, huma.NewError(http.StatusNotFound, "category not found")
		}
		return &ResponseStatus{StatusCode: http.StatusOK}, nil
	}
}

func adaptDelete(uc *category_case.CategoryCase) func(context.Context, *RequestDelete) (*ResponseStatus, error) {
	return func(ctx context.Context, input *RequestDelete) (*ResponseStatus, error) {
		deleted, err := uc.Delete(ctx, category_gateway.DeleteFilter{ID: input.ID})
		if err != nil {
			return nil, err
		}
		if !deleted {
			return nil, huma.NewError(http.StatusNotFound, "category not found")
		}
		return &ResponseStatus{StatusCode: http.StatusNoContent}, nil
	}
}

func adaptGetByID(uc *category_case.CategoryCase) func(context.Context, *RequestGetByID) (*category_gateway.GetByIDOutput, error) {
	return func(ctx context.Context, input *RequestGetByID) (*category_gateway.GetByIDOutput, error) {
		out, err := uc.GetByID(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		if out == nil {
			return nil, huma.NewError(http.StatusNotFound, "category not found")
		}
		return out, nil
	}
}

func adaptPaginate(uc *category_case.CategoryCase) func(context.Context, *RequestPaginate) (*category_gateway.PaginateOutput, error) {
	return func(ctx context.Context, input *RequestPaginate) (*category_gateway.PaginateOutput, error) {
		var tenantID, name *string
		if input.TenantID != "" {
			tenantID = &input.TenantID
		}
		if input.Name != "" {
			name = &input.Name
		}
		filter := category_gateway.PaginateFilter{TenantID: tenantID, Name: name}
		paginate := database.PaginateInput{Page: input.Page, PageSize: input.PageSize}
		out, err := uc.Paginate(ctx, filter, paginate)
		if err != nil {
			return nil, err
		}
		return &out, nil
	}
}
