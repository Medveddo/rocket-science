package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Medveddo/rocket-science/order/internal/converter"
	"github.com/Medveddo/rocket-science/order/internal/model"
	orderV1 "github.com/Medveddo/rocket-science/shared/pkg/openapi/order/v1"
)

// PayOrder implements order_v1.Handler.
func (a *orderApi) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	request := converter.PayOrderRequestFromAPI(*req, params)
	response, err := a.orderService.PayOrder(ctx, request)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Order with UUID: '" + params.OrderUUID + "' not found",
			}, nil
		}
		if errors.Is(err, model.ErrPaymentMethodIsNotSupported) {
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("Payment method %v is not supported", req.PaymentMethod),
			}, nil
		}

		return &orderV1.InternalServerError{
			Code:    500,
			Message: err.Error(),
		}, err
	}

	apiResponse := converter.PayOrderResponseToAPI(response)
	return &apiResponse, nil
}
