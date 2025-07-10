package order

import (
	"github.com/Medveddo/rocket-science/order/internal/repository"
	"github.com/Medveddo/rocket-science/order/internal/service"
	inventoryV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/payment/v1"
)

var _ service.OrderService = (*orderService)(nil)

type orderService struct {
	repo            repository.OrderRepository
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

func NewOrderService(
	repo repository.OrderRepository,
	inventoryClient inventoryV1.InventoryServiceClient,
	paymentClient paymentV1.PaymentServiceClient,
) *orderService {
	return &orderService{repo: repo, inventoryClient: inventoryClient, paymentClient: paymentClient}
}
