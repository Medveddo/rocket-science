package service

import (
	"context"

	"github.com/Medveddo/rocket-science/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, req *model.PayOrderRequest) (*model.PayOrderResponse, error)
}
