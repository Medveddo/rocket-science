package payment

import (
	"context"

	"go.uber.org/zap"

	"github.com/Medveddo/rocket-science/payment/internal/model"
	"github.com/Medveddo/rocket-science/platform/pkg/logger"
)

func (s *paymentService) PayOrder(ctx context.Context, req *model.PayOrderRequest) (*model.PayOrderResponse, error) {
	if req.PaymentMethod == model.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		logger.Warn(ctx, "Invalid payment method", zap.String("order_id", req.OrderID))
		return nil, model.ErrInvalidPaymentMethod
	}

	// Process payment using the payment client
	transactionUUID, err := s.paymentClient.Pay(ctx, req.OrderID)
	if err != nil {
		logger.Error(ctx, "Payment failed", zap.String("order_id", req.OrderID), zap.Error(err))
		return nil, err
	}

	logger.Info(ctx, "Payment successful", zap.String("order_id", req.OrderID), zap.String("transaction_uuid", transactionUUID))
	return &model.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
