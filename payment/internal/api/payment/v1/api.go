package v1

import (
	"github.com/Medveddo/rocket-science/payment/internal/service"
	paymentV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/payment/v1"
)

type paymentAPI struct {
	paymentV1.UnimplementedPaymentServiceServer
	Service service.PaymentService
}

func NewPaymentAPI(svc service.PaymentService) *paymentAPI {
	return &paymentAPI{Service: svc}
}
