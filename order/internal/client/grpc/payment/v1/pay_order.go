package v1

import (
	"context"

	"github.com/Medveddo/rocket-science/order/internal/client/converter"
	"github.com/Medveddo/rocket-science/order/internal/model"
	paymentV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/payment/v1"
)

func (c *paymentClient) PayOrder(ctx context.Context, request model.PayOrderRequest, userUUID string) (model.PayOrderResponse, error) {
	paymentMethod := converter.PaymentMethodToProto(request.PaymentMethod)

	response, err := c.client.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     request.OrderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return model.PayOrderResponse{}, err
	}

	return model.PayOrderResponse{
		TransactionUUID: response.GetTransactionUuid(),
	}, nil
}
