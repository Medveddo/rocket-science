package order

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/Medveddo/rocket-science/order/internal/model"
	paymentV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/payment/v1"
)

func (s *orderService) PayOrder(ctx context.Context, request model.PayOrderRequest) (model.PayOrderResponse, error) {
	order, err := s.repo.Get(ctx, request.OrderUUID)
	if err != nil {
		return model.PayOrderResponse{}, err
	}

	paymentMethodsMap := map[model.PayOrderRequestPaymentMethod]paymentV1.PaymentMethod{
		model.PayOrderRequestPaymentMethodCARD:          paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		model.PayOrderRequestPaymentMethodSBP:           paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		model.PayOrderRequestPaymentMethodCREDITCARD:    paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		model.PayOrderRequestPaymentMethodINVESTORMONEY: paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
	}

	paymentMethod, ok := paymentMethodsMap[request.PaymentMethod]
	if !ok {
		return model.PayOrderResponse{}, model.ErrPaymentMethodIsNotSupported
	}

	response, err := s.paymentClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		UserUuid:      order.UserUUID.String(),
		OrderUuid:     order.OrderUUID.String(),
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		log.Println("failed to process payment:", err)
		return model.PayOrderResponse{}, err
	}

	transactionUUID, err := uuid.Parse(response.TransactionUuid)
	if err != nil {
		log.Println("failed to parse transaction UUID:", err)
		return model.PayOrderResponse{}, err
	}

	order.Status = model.OrderStatusPAID
	order.PaymentMethod = &request.PaymentMethod
	order.TransactionUUID = &transactionUUID

	err = s.repo.Update(ctx, order)
	if err != nil {
		return model.PayOrderResponse{}, err
	}

	return model.PayOrderResponse{
		TransactionUUID: transactionUUID.String(),
	}, nil
}
