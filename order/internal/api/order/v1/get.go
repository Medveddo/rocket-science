package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/Medveddo/rocket-science/order/internal/converter"
	"github.com/Medveddo/rocket-science/order/internal/model"
	orderV1 "github.com/Medveddo/rocket-science/shared/pkg/openapi/order/v1"
)

// GetOrder implements order_v1.Handler.
func (h *orderApi) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	order, err := h.orderService.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Order with UUID: '" + params.OrderUUID + "' not found",
			}, nil
		}
		return nil, h.NewError(ctx, err)
	}

	apiOrder := converter.OrderToAPI(order)

	return &apiOrder, nil
}
