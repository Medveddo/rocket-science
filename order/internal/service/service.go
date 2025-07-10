package service

import (
	"context"

	"github.com/Medveddo/rocket-science/order/internal/model"
)

type OrderService interface {
	CancelOrder(ctx context.Context, orderUUID string) (model.Order, error)
	PayOrder(ctx context.Context, request model.PayOrderRequest) (model.PayOrderResponse, error)
	CreateOrder(ctx context.Context, order model.CreateOrderRequest) (model.CreateOrderResponse, error)
	GetOrder(ctx context.Context, orderUUID string) (model.Order, error)
}
