package converter

import (
	"github.com/Medveddo/rocket-science/order/internal/model"
	repoModel "github.com/Medveddo/rocket-science/order/internal/repository/model"
)

func OrderToRepoOrder(o model.Order) repoModel.Order {
	return repoModel.Order{
		OrderUUID:       o.OrderUUID,
		UserUUID:        o.UserUUID,
		PartUUIDs:       o.PartUUIDs,
		TotalPrice:      o.TotalPrice,
		TransactionUUID: o.TransactionUUID,
		PaymentMethod:   (*string)(o.PaymentMethod),
		Status:          string(o.Status),
	}
}

func RepoOrderToOrder(o repoModel.Order) model.Order {
	return model.Order{
		OrderUUID:       o.OrderUUID,
		UserUUID:        o.UserUUID,
		PartUUIDs:       o.PartUUIDs,
		TotalPrice:      o.TotalPrice,
		TransactionUUID: o.TransactionUUID,
		PaymentMethod:   (*model.PayOrderRequestPaymentMethod)(o.PaymentMethod),
		Status:          model.OrderStatus(o.Status),
	}
}
