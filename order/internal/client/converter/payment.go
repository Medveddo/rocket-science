package converter

import (
	"github.com/Medveddo/rocket-science/order/internal/model"
	paymentV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/payment/v1"
)

func PaymentMethodToProto(method model.PayOrderRequestPaymentMethod) paymentV1.PaymentMethod {
	switch method {
	case model.PayOrderRequestPaymentMethodCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PayOrderRequestPaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PayOrderRequestPaymentMethodCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PayOrderRequestPaymentMethodINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
