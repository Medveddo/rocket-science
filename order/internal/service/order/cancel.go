package order

import (
	"context"

	"github.com/Medveddo/rocket-science/order/internal/model"
)

func (s *orderService) CancelOrder(ctx context.Context, orderUUID string) (model.Order, error) {
	order, err := s.repo.Get(ctx, orderUUID)
	if err != nil {
		return model.Order{}, err
	}

	if order.Status == model.OrderStatusPAID {
		return model.Order{}, model.ErrOrderAlreadyPaid
	}

	order.Status = model.OrderStatusCANCELLED

	err = s.repo.Update(ctx, order)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
