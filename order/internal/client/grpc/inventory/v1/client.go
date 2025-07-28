package v1

import (
	"context"

	"github.com/Medveddo/rocket-science/order/internal/client/converter"
	"github.com/Medveddo/rocket-science/order/internal/model"
	inventoryV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/inventory/v1"
)

type inventoryClient struct {
	client inventoryV1.InventoryServiceClient
}

func NewInventoryClientV1(client inventoryV1.InventoryServiceClient) *inventoryClient {
	return &inventoryClient{client: client}
}

func (c *inventoryClient) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	response, err := c.client.GetPart(ctx, &inventoryV1.GetPartRequest{
		Uuid: uuid,
	})
	if err != nil {
		return model.Part{}, err
	}

	return converter.PartFromProto(response.GetPart()), nil
}
