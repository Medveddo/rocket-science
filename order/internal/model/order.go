package model

import (
	"github.com/google/uuid"
)

type Order struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PartUUIDs       []string
	TotalPrice      float64
	TransactionUUID *uuid.UUID
	PaymentMethod   *PayOrderRequestPaymentMethod
	Status          OrderStatus
}

type PayOrderRequest struct {
	OrderUUID     string
	PaymentMethod PayOrderRequestPaymentMethod
}

type PayOrderResponse struct {
	TransactionUUID string
}

type PayOrderRequestPaymentMethod string

const (
	PayOrderRequestPaymentMethodCARD          PayOrderRequestPaymentMethod = "CARD"
	PayOrderRequestPaymentMethodSBP           PayOrderRequestPaymentMethod = "SBP"
	PayOrderRequestPaymentMethodCREDITCARD    PayOrderRequestPaymentMethod = "CREDIT_CARD"
	PayOrderRequestPaymentMethodINVESTORMONEY PayOrderRequestPaymentMethod = "INVESTOR_MONEY"
)

type OrderStatus string

const (
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"
)

type CreateOrderRequest struct {
	UserUUID  uuid.UUID
	PartUuids []string
}

type CreateOrderResponse struct {
	OrderUUID  uuid.UUID
	TotalPrice float64
}
