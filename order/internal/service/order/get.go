package order

import (
	"context"

	"github.com/Medveddo/rocket-science/order/internal/model"
)

func (s *orderService) GetOrder(ctx context.Context, orderUUID string) (model.Order, error) {
	return s.repo.Get(ctx, orderUUID)
}
