package order

import (
	grpcClient "github.com/Medveddo/rocket-science/order/internal/client/grpc"
	"github.com/Medveddo/rocket-science/order/internal/repository"
	"github.com/Medveddo/rocket-science/order/internal/service"
)

var _ service.OrderService = (*orderService)(nil)

type orderService struct {
	repo            repository.OrderRepository
	inventoryClient grpcClient.InventoryClient
	paymentClient   grpcClient.PaymentClient
}

func NewOrderService(
	repo repository.OrderRepository,
	inventoryClient grpcClient.InventoryClient,
	paymentClient grpcClient.PaymentClient,
) *orderService {
	return &orderService{repo: repo, inventoryClient: inventoryClient, paymentClient: paymentClient}
}
