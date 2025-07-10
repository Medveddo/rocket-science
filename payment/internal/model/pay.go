package model

type PayOrderRequest struct {
	OrderID       string
	PaymentMethod PaymentMethod
}

type PayOrderResponse struct {
	TransactionUUID string
}

type PaymentMethod int32

const (
	PaymentMethod_PAYMENT_METHOD_UNSPECIFIED    PaymentMethod = 0
	PaymentMethod_PAYMENT_METHOD_CARD           PaymentMethod = 1
	PaymentMethod_PAYMENT_METHOD_SBP            PaymentMethod = 2
	PaymentMethod_PAYMENT_METHOD_CREDIT_CARD    PaymentMethod = 3
	PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY PaymentMethod = 4
)
