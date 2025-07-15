package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Medveddo/rocket-science/order/internal/service/mocks"
)

type APISuite struct {
	suite.Suite

	ctx context.Context //nolint:containedctx

	orderService *mocks.OrderService

	api *orderApi
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	s.orderService = mocks.NewOrderService(s.T())

	s.api = NewOrderAPI(s.orderService)
}

func (s *APISuite) TearDownTest() {}

func TestPaymentAPI(t *testing.T) {
	suite.Run(t, new(APISuite))
}
