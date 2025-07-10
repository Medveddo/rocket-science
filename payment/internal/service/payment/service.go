package payment

import (
	"github.com/Medveddo/rocket-science/payment/internal/service"
)

var _ service.PaymentService = (*paymentService)(nil)

type paymentService struct{}

func NewService() *paymentService {
	return &paymentService{}
}
