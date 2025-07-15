package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Medveddo/rocket-science/payment/internal/service/mocks"
)

type APISuite struct {
	suite.Suite

	ctx context.Context //nolint:containedctx

	paymentService *mocks.PaymentService

	api *paymentAPI
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	s.paymentService = mocks.NewPaymentService(s.T())

	s.api = NewPaymentAPI(s.paymentService)
}

func (s *APISuite) TearDownTest() {}

func TestPaymentAPI(t *testing.T) {
	suite.Run(t, new(APISuite))
}
