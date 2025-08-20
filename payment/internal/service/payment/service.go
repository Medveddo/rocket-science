package payment

import (
	"github.com/Medveddo/rocket-science/payment/internal/client"
	"github.com/Medveddo/rocket-science/payment/internal/service"
)

var _ service.PaymentService = (*paymentService)(nil)

type paymentService struct {
	paymentClient client.PaymentClient
}

func NewService(paymentClient client.PaymentClient) *paymentService {
	return &paymentService{
		paymentClient: paymentClient,
	}
}
