package client

import (
	"context"

	"github.com/google/uuid"
)

type PaymentClient interface {
	Pay(ctx context.Context, orderUUID string) (string, error)
}

type RandomPaymentClient struct{}

func NewRandomPaymentClient() *RandomPaymentClient {
	return &RandomPaymentClient{}
}

func (c *RandomPaymentClient) Pay(ctx context.Context, orderUUID string) (string, error) {
	transactionUUID := uuid.New()
	return transactionUUID.String(), nil
}
