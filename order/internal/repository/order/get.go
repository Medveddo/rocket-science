package order

import (
	"context"

	"github.com/Medveddo/rocket-science/order/internal/model"
	"github.com/Medveddo/rocket-science/order/internal/repository/converter"
)

func (r *orderRepository) Get(ctx context.Context, orderUUID string) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, ok := r.orders[orderUUID]
	if !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	return converter.RepoOrderToOrder(*order), nil
}
