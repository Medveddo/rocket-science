package v1

import (
	"errors"
	"net/http"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/Medveddo/rocket-science/order/internal/model"
	orderV1 "github.com/Medveddo/rocket-science/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestCancelOrder_Success() {
	orderUUID := gofakeit.UUID()
	s.orderService.On("CancelOrder", s.ctx, orderUUID).Return(model.Order{}, nil)

	res, err := s.api.CancelOrder(s.ctx, orderV1.CancelOrderParams{OrderUUID: orderUUID})

	s.Require().NoError(err)
	_, ok := res.(*orderV1.CancelOrderNoContent)
	s.Require().True(ok)
}

func (s *APISuite) TestCancelOrder_NotFound() {
	orderUUID := gofakeit.UUID()
	s.orderService.On("CancelOrder", s.ctx, orderUUID).Return(model.Order{}, model.ErrOrderNotFound)

	res, err := s.api.CancelOrder(s.ctx, orderV1.CancelOrderParams{OrderUUID: orderUUID})
	s.Require().NoError(err)

	notFoundErr, ok := res.(*orderV1.NotFoundError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusNotFound, notFoundErr.Code)
	s.Require().Contains(notFoundErr.Message, orderUUID)
}

func (s *APISuite) TestCancelOrder_AlreadyPaid() {
	orderUUID := gofakeit.UUID()
	s.orderService.On("CancelOrder", s.ctx, orderUUID).Return(model.Order{}, model.ErrOrderAlreadyPaid)

	res, err := s.api.CancelOrder(s.ctx, orderV1.CancelOrderParams{OrderUUID: orderUUID})
	s.Require().NoError(err)

	conflictErr, ok := res.(*orderV1.NotFoundError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusConflict, conflictErr.Code)
	s.Require().Contains(conflictErr.Message, "already paid")
}

func (s *APISuite) TestCancelOrder_InternalError() {
	orderUUID := gofakeit.UUID()
	testErr := gofakeit.Error()
	s.orderService.On("CancelOrder", s.ctx, orderUUID).Return(model.Order{}, testErr)

	res, err := s.api.CancelOrder(s.ctx, orderV1.CancelOrderParams{OrderUUID: orderUUID})
	s.Require().Nil(res)
	s.Require().Error(err)

	var genericErr *orderV1.GenericErrorStatusCode
	s.Require().True(errors.As(err, &genericErr))
	s.Require().Equal(http.StatusInternalServerError, genericErr.StatusCode)
	s.Require().Contains(genericErr.Response.Message.Value, testErr.Error())
}
