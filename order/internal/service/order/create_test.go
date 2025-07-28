package order

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/Medveddo/rocket-science/order/internal/model"
)

func (s *ServiceSuite) TestCreateOrder_Success() {
	userUUID := uuid.New()
	partUUIDs := []string{"part-1", "part-2"}

	s.mockInventoryClient.EXPECT().ListParts(s.ctx, model.PartsFilter{UUIDs: partUUIDs}).Return([]model.Part{
		{UUID: "part-1", Price: 100.0},
		{UUID: "part-2", Price: 200.0},
	}, nil)

	s.mockRepo.EXPECT().Update(s.ctx, mock.AnythingOfType("model.Order")).Return(nil)

	request := model.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: partUUIDs,
	}

	response, err := s.service.CreateOrder(s.ctx, request)

	s.Require().NoError(err)
	s.Require().NotEmpty(response.OrderUUID)
	s.Require().Equal(300.0, response.TotalPrice)
}

func (s *ServiceSuite) TestCreateOrder_InventoryError() {
	userUUID := uuid.New()
	partUUIDs := []string{"part-1", "part-2"}

	s.mockInventoryClient.EXPECT().ListParts(s.ctx, model.PartsFilter{UUIDs: partUUIDs}).Return(nil, model.ErrFailedToFetchInventory)

	request := model.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: partUUIDs,
	}

	response, err := s.service.CreateOrder(s.ctx, request)

	s.Require().Error(err)
	s.Require().Equal(model.ErrFailedToFetchInventory, err)
	s.Require().Empty(response.OrderUUID)
	s.Require().Equal(0.0, response.TotalPrice)
}

func (s *ServiceSuite) TestCreateOrder_MissingParts() {
	userUUID := uuid.New()
	partUUIDs := []string{"part-1", "part-2", "part-3"}

	s.mockInventoryClient.EXPECT().ListParts(s.ctx, model.PartsFilter{UUIDs: partUUIDs}).Return([]model.Part{
		{UUID: "part-1", Price: 100.0},
		{UUID: "part-2", Price: 200.0},
		// part-3 is missing
	}, nil)

	request := model.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: partUUIDs,
	}

	response, err := s.service.CreateOrder(s.ctx, request)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "part-3")
	s.Require().Contains(err.Error(), "part does not exist")
	s.Require().Empty(response.OrderUUID)
	s.Require().Equal(0.0, response.TotalPrice)
}

func (s *ServiceSuite) TestCreateOrder_RepositoryError() {
	userUUID := uuid.New()
	partUUIDs := []string{"part-1", "part-2"}

	s.mockInventoryClient.EXPECT().ListParts(s.ctx, model.PartsFilter{UUIDs: partUUIDs}).Return([]model.Part{
		{UUID: "part-1", Price: 100.0},
		{UUID: "part-2", Price: 200.0},
	}, nil)

	s.mockRepo.EXPECT().Update(s.ctx, mock.AnythingOfType("model.Order")).Return(model.ErrOrderNotFound)

	request := model.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: partUUIDs,
	}

	response, err := s.service.CreateOrder(s.ctx, request)

	s.Require().Error(err)
	s.Require().Equal(model.ErrOrderNotFound, err)
	s.Require().Empty(response.OrderUUID)
	s.Require().Equal(0.0, response.TotalPrice)
}
