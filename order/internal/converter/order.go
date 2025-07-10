package converter

import (
	"github.com/google/uuid"

	"github.com/Medveddo/rocket-science/order/internal/model"
	orderV1 "github.com/Medveddo/rocket-science/shared/pkg/openapi/order/v1"
)

func CreateOrderRequestFromAPI(request orderV1.CreateOrderRequest) model.CreateOrderRequest {
	return model.CreateOrderRequest{
		UserUUID:  request.UserUUID,
		PartUuids: request.GetPartUuids(),
	}
}

func CreateOrderResponseToAPI(response model.CreateOrderResponse) orderV1.CreateOrderResponse {
	return orderV1.CreateOrderResponse{
		OrderUUID:  response.OrderUUID,
		TotalPrice: float32(response.TotalPrice),
	}
}

func OrderToAPI(order model.Order) orderV1.OrderDto {
	dto := orderV1.OrderDto{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUUIDs,
		TotalPrice: order.TotalPrice,
		Status:     orderV1.OrderStatus(order.Status),
	}
	if order.PaymentMethod != nil {
		dto.PaymentMethod = orderV1.NewOptPaymentMethod(
			orderV1.PaymentMethod(string(*order.PaymentMethod)),
		)
	}
	if order.TransactionUUID != nil {
		dto.TransactionUUID = orderV1.NewOptNilUUID(*order.TransactionUUID)
	}
	return dto
}

func PayOrderRequestFromAPI(request orderV1.PayOrderRequest, params orderV1.PayOrderParams) model.PayOrderRequest {
	return model.PayOrderRequest{
		OrderUUID:     params.OrderUUID,
		PaymentMethod: model.PayOrderRequestPaymentMethod(request.PaymentMethod),
	}
}

func PayOrderResponseToAPI(response model.PayOrderResponse) orderV1.PayOrderResponse {
	return orderV1.PayOrderResponse{
		TransactionUUID: uuid.MustParse(response.TransactionUUID),
	}
}
