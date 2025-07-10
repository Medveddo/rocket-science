package v1

import (
	"context"
	"net/http"

	"github.com/Medveddo/rocket-science/order/internal/service"
	orderV1 "github.com/Medveddo/rocket-science/shared/pkg/openapi/order/v1"
)

var _ orderV1.Handler = (*orderApi)(nil)

type orderApi struct {
	orderService service.OrderService
}

func NewOrderAPI(orderService service.OrderService) *orderApi {
	return &orderApi{
		orderService: orderService,
	}
}

// NewError создает новую ошибку в формате GenericError
func (a *orderApi) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}
