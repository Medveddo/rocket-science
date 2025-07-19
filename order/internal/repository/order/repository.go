package order

import (
	sq "github.com/Masterminds/squirrel"
	// repoModel "github.com/Medveddo/rocket-science/order/internal/repository/model"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Medveddo/rocket-science/order/internal/repository"
)

var _ repository.OrderRepository = (*orderRepository)(nil)

type orderRepository struct {
	pool      *pgxpool.Pool
	psql      sq.StatementBuilderType
	tableName string
}

func NewOrderRepository(pool *pgxpool.Pool) *orderRepository {
	return &orderRepository{
		pool:      pool,
		tableName: "orders",
		psql:      sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
