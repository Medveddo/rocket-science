package app

import (
	"context"

	paymentV1API "github.com/Medveddo/rocket-science/payment/internal/api/payment/v1"
	"github.com/Medveddo/rocket-science/payment/internal/client"
	"github.com/Medveddo/rocket-science/payment/internal/service"
	paymentService "github.com/Medveddo/rocket-science/payment/internal/service/payment"
	paymentV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1API paymentV1.PaymentServiceServer

	paymentService service.PaymentService
	paymentClient  client.PaymentClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1API(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = paymentV1API.NewPaymentAPI(d.PaymentService(ctx))
	}

	return d.paymentV1API
}

func (d *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService(d.PaymentClient(ctx))
	}

	return d.paymentService
}

func (d *diContainer) PaymentClient(ctx context.Context) client.PaymentClient {
	if d.paymentClient == nil {
		d.paymentClient = client.NewRandomPaymentClient()
	}

	return d.paymentClient
}
