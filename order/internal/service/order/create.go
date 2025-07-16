package order

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/Medveddo/rocket-science/order/internal/model"
	inventoryV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/inventory/v1"
)

func (s *orderService) CreateOrder(ctx context.Context, request model.CreateOrderRequest) (model.CreateOrderResponse, error) {
	listPartsResponse, err := s.inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: request.PartUuids,
		},
	})
	if err != nil {
		log.Printf("error while fetching inventory: %v\n", err)
		return model.CreateOrderResponse{}, model.ErrFailedToFetchInventory
	}

	parts := listPartsResponse.GetParts()

	returned := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		returned[part.GetUuid()] = struct{}{}
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
		totalPrice += part.GetPrice()
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
