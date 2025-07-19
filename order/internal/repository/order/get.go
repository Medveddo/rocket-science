package order

import (
	"context"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/Medveddo/rocket-science/order/internal/model"
	"github.com/Medveddo/rocket-science/order/internal/repository/converter"
	repoModel "github.com/Medveddo/rocket-science/order/internal/repository/model"
)

func (r *orderRepository) Get(ctx context.Context, orderUUID string) (model.Order, error) {
	query, args, err := r.psql.
		Select("order_uuid", "user_uuid", "part_uuids", "total_price", "transaction_uuid", "payment_method", "status").
		From(r.tableName).
		Where(sq.Eq{"order_uuid": orderUUID}).
		ToSql()
	if err != nil {
		log.Println(err)
		return model.Order{}, err
	}

	var o repoModel.Order
	// var partUuids []string
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&o.OrderUUID,
		&o.UserUUID,
		&o.PartUUIDs,
		&o.TotalPrice,
		&o.TransactionUUID,
		&o.PaymentMethod,
		&o.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Order{}, model.ErrOrderNotFound
		}
		return model.Order{}, err
	}
	return converter.RepoOrderToOrder(o), nil
}
