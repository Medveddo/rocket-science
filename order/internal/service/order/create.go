package order

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/Medveddo/rocket-science/order/internal/model"
)

func (s *orderService) CreateOrder(ctx context.Context, request model.CreateOrderRequest) (model.CreateOrderResponse, error) {
	parts, err := s.inventoryClient.ListParts(ctx, model.PartsFilter{
		UUIDs: request.PartUuids,
	})
	if err != nil {
		log.Printf("error while fetching inventory: %v\n", err)
		return model.CreateOrderResponse{}, model.ErrFailedToFetchInventory
	}

	returned := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		returned[part.UUID] = struct{}{}
	}

	var missing []string
	for _, reqUUID := range request.PartUuids {
		if _, ok := returned[reqUUID]; !ok {
			missing = append(missing, reqUUID)
		}
	}
	if len(missing) > 0 {
		err = fmt.Errorf("the following partUuid(s) do not exist: %v: %w", missing, model.ErrPartDoesNotExist)
		return model.CreateOrderResponse{}, err
	}

	var totalPrice float64
	for _, part := range parts {
		totalPrice += part.Price
	}

	orderUuid := uuid.New()

	order := model.Order{
		UserUUID:   request.UserUUID,
		OrderUUID:  orderUuid,
		PartUUIDs:  request.PartUuids,
		Status:     model.OrderStatusPENDINGPAYMENT,
		TotalPrice: totalPrice,
	}
	err = s.repo.Update(ctx, order)
	if err != nil {
		return model.CreateOrderResponse{}, err
	}

	response := model.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}
	return response, nil
}
