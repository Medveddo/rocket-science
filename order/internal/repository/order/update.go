package order

import (
	"context"

	"github.com/Medveddo/rocket-science/order/internal/model"
)

func (r *orderRepository) Update(ctx context.Context, order model.Order) error {
	suffix := `ON CONFLICT (order_uuid) DO UPDATE
		SET user_uuid = EXCLUDED.user_uuid,
		part_uuids = EXCLUDED.part_uuids,
		total_price = EXCLUDED.total_price,
		transaction_uuid = EXCLUDED.transaction_uuid,
		payment_method = EXCLUDED.payment_method,
		status = EXCLUDED.status`

	query, args, err := r.psql.Insert(r.tableName).
		Columns("order_uuid", "user_uuid", "part_uuids", "total_price", "transaction_uuid", "payment_method", "status").
		Values(
			order.OrderUUID,
			order.UserUUID,
			order.PartUUIDs,
			order.TotalPrice,
			order.TransactionUUID,
			order.PaymentMethod,
			order.Status,
		).Suffix(suffix).ToSql()
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, query, args...)
	return err
}
