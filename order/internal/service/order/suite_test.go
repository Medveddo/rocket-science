package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	inventoryMocks "github.com/Medveddo/rocket-science/order/internal/client/grpc/inventory/v1/mocks"
	paymentMocks "github.com/Medveddo/rocket-science/order/internal/client/grpc/payment/v1/mocks"
	repositoryMocks "github.com/Medveddo/rocket-science/order/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context //nolint:containedctx

	service             *orderService
	mockRepo            *repositoryMocks.OrderRepository
	mockInventoryClient *inventoryMocks.Client
	mockPaymentClient   *paymentMocks.Client
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.mockRepo = repositoryMocks.NewOrderRepository(s.T())
	s.mockInventoryClient = inventoryMocks.NewClient(s.T())
	s.mockPaymentClient = paymentMocks.NewClient(s.T())

	s.service = NewOrderService(s.mockRepo, s.mockInventoryClient, s.mockPaymentClient)
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
