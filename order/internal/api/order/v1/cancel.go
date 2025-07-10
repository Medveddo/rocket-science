package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/Medveddo/rocket-science/order/internal/model"
	orderV1 "github.com/Medveddo/rocket-science/shared/pkg/openapi/order/v1"
)

// CancelOrder implements order_v1.Handler.
func (h *orderApi) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	_, err := h.orderService.CancelOrder(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Order with UUID: '" + params.OrderUUID + "' not found",
			}, nil
		}
		if errors.Is(err, model.ErrOrderAlreadyPaid) {
			return &orderV1.NotFoundError{
				Code:    http.StatusConflict,
				Message: "Order already paid and cannot be cancelled",
			}, nil
		}
		return nil, h.NewError(ctx, err)
	}

	return &orderV1.CancelOrderNoContent{}, nil
}
