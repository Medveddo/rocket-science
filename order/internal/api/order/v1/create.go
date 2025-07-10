package v1

import (
	"context"
	"errors"

	"github.com/Medveddo/rocket-science/order/internal/converter"
	"github.com/Medveddo/rocket-science/order/internal/model"
	orderV1 "github.com/Medveddo/rocket-science/shared/pkg/openapi/order/v1"
)

func (a *orderApi) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	modelRequest := converter.CreateOrderRequestFromAPI(*req)

	response, err := a.orderService.CreateOrder(ctx, modelRequest)
	if err != nil {
		if errors.Is(err, model.ErrPartDoesNotExist) {
			return &orderV1.BadRequestError{
				Code:    400,
				Message: err.Error(),
			}, nil
		}
		return nil, a.NewError(ctx, err)
	}

	apiResponse := converter.CreateOrderResponseToAPI(response)

	return &apiResponse, nil
}
