package payment

// import (
// 	"github.com/brianvoe/gofakeit/v7"

// 	"github.com/Medveddo/rocket-science/payment/internal/model"
// )

// func (s *ServiceSuite) TestCreateSuccess() {
// 	var (
// 		expectedUUID = gofakeit.UUID()

// 		request = model.PayOrderRequest{
// 			OrderID:       gofakeit.UUID(),
// 			PaymentMethod: model.PaymentMethod_PAYMENT_METHOD_CARD,
// 		}

// 		expectedResponse = model.PayOrderResponse{
// 			TransactionUUID: expectedUUID,
// 		}
// 	)

// 	response, err := s.service.PayOrder(s.ctx, &request)
// 	s.Require().NoError(err)
// 	s.Require().NotNil(response)
// 	s.Require().Equal(response.TransactionUUID, expectedUUID)
// 	s.Require().Equal(response, expectedResponse)
// }
