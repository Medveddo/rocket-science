package client

import (
	"context"

	"github.com/Medveddo/rocket-science/order/internal/model"
)

type PaymentClient interface {
	PayOrder(ctx context.Context, request model.PayOrderRequest, userUUID string) (model.PayOrderResponse, error)
}

type InventoryClient interface {
	GetPart(ctx context.Context, uuid string) (model.Part, error)
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}
