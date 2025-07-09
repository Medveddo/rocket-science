package converter

import (
	"github.com/Medveddo/rocket-science/payment/internal/model"
	paymentV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/payment/v1"
)

func FromProtoPayOrderRequest(req *paymentV1.PayOrderRequest) *model.PayOrderRequest {
	return &model.PayOrderRequest{
		OrderID:       req.GetOrderUuid(),
		PaymentMethod: model.PaymentMethod((req.GetPaymentMethod())),
	}
}

func ToProtoPayOrderResponse(res *model.PayOrderResponse) *paymentV1.PayOrderResponse {
	return &paymentV1.PayOrderResponse{
		TransactionUuid: res.TransactionUUID,
	}
}
