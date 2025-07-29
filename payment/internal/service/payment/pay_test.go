package payment

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/Medveddo/rocket-science/payment/internal/model"
)

func (s *ServiceSuite) TestPaymentService_PayOrder_Success() {
	expectedOrderUUID := "test-order-123"
	expectedTransactionUUID := "test-transaction-456"

	request := &model.PayOrderRequest{
		OrderID:       expectedOrderUUID,
		PaymentMethod: model.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	s.paymentClientMock.EXPECT().Pay(s.ctx, expectedOrderUUID).Return(expectedTransactionUUID, nil)

	response, err := s.service.PayOrder(s.ctx, request)

	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().Equal(expectedTransactionUUID, response.TransactionUUID)
}

func (s *ServiceSuite) TestPaymentService_PayOrder_InvalidPaymentMethod() {
	expectedOrderUUID := "test-order-123"

	request := &model.PayOrderRequest{
		OrderID:       expectedOrderUUID,
		PaymentMethod: model.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
	}

	response, err := s.service.PayOrder(s.ctx, request)

	s.Require().Error(err)
	s.Require().Nil(response)
	s.Require().Equal(model.ErrInvalidPaymentMethod, err)
}

func (s *ServiceSuite) TestPaymentService_PayOrder_PaymentClientError() {
	expectedOrderUUID := "test-order-123"

	request := &model.PayOrderRequest{
		OrderID:       expectedOrderUUID,
		PaymentMethod: model.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	s.paymentClientMock.EXPECT().Pay(s.ctx, expectedOrderUUID).Return("", gofakeit.Error())

	response, err := s.service.PayOrder(s.ctx, request)

	s.Require().Error(err)
	s.Require().Nil(response)
}
