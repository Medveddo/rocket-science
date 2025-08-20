package payment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Medveddo/rocket-science/payment/internal/client/mocks"
	"github.com/Medveddo/rocket-science/platform/pkg/logger"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context //nolint:containedctx

	service           *paymentService
	paymentClientMock *mocks.PaymentClient
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	logger.SetNopLogger()

	s.paymentClientMock = mocks.NewPaymentClient(s.T())
	s.service = NewService(s.paymentClientMock)
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
