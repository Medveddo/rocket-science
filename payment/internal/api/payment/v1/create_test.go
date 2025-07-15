package v1

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/Medveddo/rocket-science/payment/internal/converter"
	"github.com/Medveddo/rocket-science/payment/internal/model"
	paymentV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/payment/v1"
)

func (s *APISuite) TestCreateSuccess() {
	var (
		expectedUUID = gofakeit.UUID()

		request = &paymentV1.PayOrderRequest{
			OrderUuid:     gofakeit.UUID(),
			PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		}

		expectedResponse = &paymentV1.PayOrderResponse{
			TransactionUuid: expectedUUID,
		}

		expectedModelInfo = converter.FromProtoPayOrderRequest(request)

		modelServiceResponse = &model.PayOrderResponse{TransactionUUID: expectedUUID}
	)

	s.paymentService.On("PayOrder", s.ctx, expectedModelInfo).Return(modelServiceResponse, nil)

	res, err := s.api.PayOrder(s.ctx, request)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedResponse, res)
}
