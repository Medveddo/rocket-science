package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/Medveddo/rocket-science/payment/internal/model"
)

func (s *paymentService) PayOrder(ctx context.Context, req *model.PayOrderRequest) (*model.PayOrderResponse, error) {
	if req.PaymentMethod == model.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return nil, model.ErrInvalidPaymentMethod
	}
	transactionUUID := uuid.New()
	return &model.PayOrderResponse{
		TransactionUUID: transactionUUID.String(),
	}, nil
}
