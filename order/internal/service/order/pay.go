package order

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/Medveddo/rocket-science/order/internal/model"
)

func (s *orderService) PayOrder(ctx context.Context, request model.PayOrderRequest) (model.PayOrderResponse, error) {
	order, err := s.repo.Get(ctx, request.OrderUUID)
	if err != nil {
		return model.PayOrderResponse{}, err
	}

	response, err := s.paymentClient.PayOrder(ctx, request, order.UserUUID.String())
	if err != nil {
		log.Println("failed to process payment:", err)
		return model.PayOrderResponse{}, err
	}

	transactionUUID, err := uuid.Parse(response.TransactionUUID)
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
