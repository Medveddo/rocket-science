package payment

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Medveddo/rocket-science/payment/internal/model"
	"github.com/Medveddo/rocket-science/platform/pkg/logger"
)

func (s *paymentService) PayOrder(ctx context.Context, req *model.PayOrderRequest) (*model.PayOrderResponse, error) {
	if req.PaymentMethod == model.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		logger.Warn(ctx, "Invalid payment method", zap.String("order_id", req.OrderID))
		return nil, model.ErrInvalidPaymentMethod
	}

	// Fake payment due to the fact that we don't have a real payment system

	transactionUUID := uuid.New()
	logger.Info(ctx, "Payment successful", zap.String("order_id", req.OrderID), zap.String("transaction_uuid", transactionUUID.String()))
	return &model.PayOrderResponse{
		TransactionUUID: transactionUUID.String(),
	}, nil
}
