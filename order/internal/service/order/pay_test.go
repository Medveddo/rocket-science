package order

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/Medveddo/rocket-science/order/internal/model"
)

func (s *ServiceSuite) TestPayOrder_Success() {
	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := "550e8400-e29b-41d4-a716-446655440000"

	existingOrder := model.Order{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.OrderStatusPENDINGPAYMENT,
	}

	s.mockRepo.EXPECT().Get(s.ctx, orderUUID.String()).Return(existingOrder, nil)
	s.mockRepo.EXPECT().Update(s.ctx, mock.AnythingOfType("model.Order")).Return(nil)

	s.mockPaymentClient.EXPECT().PayOrder(s.ctx, mock.AnythingOfType("model.PayOrderRequest"), userUUID.String()).Return(model.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil)

	request := model.PayOrderRequest{
		OrderUUID:     orderUUID.String(),
		PaymentMethod: model.PayOrderRequestPaymentMethodCARD,
	}

	response, err := s.service.PayOrder(s.ctx, request)

	s.Require().NoError(err)
	s.Require().Equal(transactionUUID, response.TransactionUUID)
}

func (s *ServiceSuite) TestPayOrder_OrderNotFound() {
	orderUUID := uuid.New()

	s.mockRepo.EXPECT().Get(s.ctx, orderUUID.String()).Return(model.Order{}, model.ErrOrderNotFound)

	request := model.PayOrderRequest{
		OrderUUID:     orderUUID.String(),
		PaymentMethod: model.PayOrderRequestPaymentMethodCARD,
	}

	response, err := s.service.PayOrder(s.ctx, request)

	s.Require().Error(err)
	s.Require().Equal(model.ErrOrderNotFound, err)
	s.Require().Empty(response.TransactionUUID)
}

func (s *ServiceSuite) TestPayOrder_PaymentError() {
	orderUUID := uuid.New()
	userUUID := uuid.New()

	existingOrder := model.Order{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.OrderStatusPENDINGPAYMENT,
	}

	s.mockRepo.EXPECT().Get(s.ctx, orderUUID.String()).Return(existingOrder, nil)

	s.mockPaymentClient.EXPECT().PayOrder(s.ctx, mock.AnythingOfType("model.PayOrderRequest"), userUUID.String()).Return(model.PayOrderResponse{}, model.ErrFailedToFetchInventory)

	request := model.PayOrderRequest{
		OrderUUID:     orderUUID.String(),
		PaymentMethod: model.PayOrderRequestPaymentMethodCARD,
	}

	response, err := s.service.PayOrder(s.ctx, request)

	s.Require().Error(err)
	s.Require().Equal(model.ErrFailedToFetchInventory, err)
	s.Require().Empty(response.TransactionUUID)
}

func (s *ServiceSuite) TestPayOrder_InvalidTransactionUUID() {
	orderUUID := uuid.New()
	userUUID := uuid.New()

	existingOrder := model.Order{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.OrderStatusPENDINGPAYMENT,
	}

	s.mockRepo.EXPECT().Get(s.ctx, orderUUID.String()).Return(existingOrder, nil)

	s.mockPaymentClient.EXPECT().PayOrder(s.ctx, mock.AnythingOfType("model.PayOrderRequest"), userUUID.String()).Return(model.PayOrderResponse{
		TransactionUUID: "invalid-uuid",
	}, nil)

	request := model.PayOrderRequest{
		OrderUUID:     orderUUID.String(),
		PaymentMethod: model.PayOrderRequestPaymentMethodCARD,
	}

	response, err := s.service.PayOrder(s.ctx, request)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "invalid UUID length")
	s.Require().Empty(response.TransactionUUID)
}

func (s *ServiceSuite) TestPayOrder_UpdateError() {
	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := "550e8400-e29b-41d4-a716-446655440000"

	existingOrder := model.Order{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.OrderStatusPENDINGPAYMENT,
	}

	s.mockRepo.EXPECT().Get(s.ctx, orderUUID.String()).Return(existingOrder, nil)
	s.mockRepo.EXPECT().Update(s.ctx, mock.AnythingOfType("model.Order")).Return(model.ErrOrderNotFound)

	s.mockPaymentClient.EXPECT().PayOrder(s.ctx, mock.AnythingOfType("model.PayOrderRequest"), userUUID.String()).Return(model.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil)

	request := model.PayOrderRequest{
		OrderUUID:     orderUUID.String(),
		PaymentMethod: model.PayOrderRequestPaymentMethodCARD,
	}

	response, err := s.service.PayOrder(s.ctx, request)

	s.Require().Error(err)
	s.Require().Equal(model.ErrOrderNotFound, err)
	s.Require().Empty(response.TransactionUUID)
}
