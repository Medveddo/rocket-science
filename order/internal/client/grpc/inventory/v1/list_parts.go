package v1

import (
	"context"

	"github.com/Medveddo/rocket-science/order/internal/client/converter"
	"github.com/Medveddo/rocket-science/order/internal/model"
	inventoryV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/inventory/v1"
)

func (c *inventoryClient) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	protoFilter := &inventoryV1.PartsFilter{
		Uuids:                 filter.UUIDs,
		Names:                 filter.Names,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}

	if len(filter.Categories) > 0 {
		protoCategories := make([]inventoryV1.Category, len(filter.Categories))
		for i, category := range filter.Categories {
			protoCategories[i] = converter.CategoryToProto(category)
		}
		protoFilter.Categories = protoCategories
	}

	response, err := c.client.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: protoFilter,
	})
	if err != nil {
		return nil, err
	}

	return converter.PartsFromProto(response.GetParts()), nil
}
