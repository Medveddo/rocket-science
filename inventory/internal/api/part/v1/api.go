package v1

import (
	"github.com/Medveddo/rocket-science/inventory/internal/service"
	inventoryV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/inventory/v1"
)

type partAPI struct {
	inventoryV1.UnimplementedInventoryServiceServer

	partService service.PartService
}

func NewPartAPI(partService service.PartService) *partAPI {
	return &partAPI{
		partService: partService,
	}
}
