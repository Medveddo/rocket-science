package payment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context //nolint:containedctx

	service *paymentService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.service = NewService()
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
