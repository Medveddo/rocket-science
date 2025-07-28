package v1

import (
	paymentV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/payment/v1"
)

type paymentClient struct {
	client paymentV1.PaymentServiceClient
}

func NewPaymentClientV1(client paymentV1.PaymentServiceClient) *paymentClient {
	return &paymentClient{client: client}
}
