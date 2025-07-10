package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Medveddo/rocket-science/payment/internal/converter"
	"github.com/Medveddo/rocket-science/payment/internal/model"
	paymentV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/payment/v1"
)

func (a *PaymentAPI) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	internalReq := converter.FromProtoPayOrderRequest(req)
	res, err := a.Service.PayOrder(ctx, internalReq)
	if err != nil {
		if errors.Is(err, model.ErrInvalidPaymentMethod) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return converter.ToProtoPayOrderResponse(res), nil
}
