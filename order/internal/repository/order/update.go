package order

import (
	"context"

	"github.com/Medveddo/rocket-science/order/internal/model"
	"github.com/Medveddo/rocket-science/order/internal/repository/converter"
)

func (r *orderRepository) Update(ctx context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	repoOrder := converter.OrderToRepoOrder(order)

	r.orders[repoOrder.OrderUUID.String()] = &repoOrder
	return nil
}
