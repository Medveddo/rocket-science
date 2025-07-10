package order

import (
	"sync"

	"github.com/Medveddo/rocket-science/order/internal/repository"
	repoModel "github.com/Medveddo/rocket-science/order/internal/repository/model"
)

var _ repository.OrderRepository = (*orderRepository)(nil)

type orderRepository struct {
	mu     sync.RWMutex
	orders map[string]*repoModel.Order
}

func NewOrderRepository() *orderRepository {
	return &orderRepository{
		orders: make(map[string]*repoModel.Order),
	}
}
