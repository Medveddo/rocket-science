package part

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/Medveddo/rocket-science/inventory/internal/model"
	u "github.com/Medveddo/rocket-science/inventory/internal/utils"
)

func (s *ServiceSuite) TestGetSuccess() {
	var (
		partUuid = gofakeit.UUID()

		expectedPart = model.Part{
			UUID:          gofakeit.UUID(),
			Name:          gofakeit.ProductName(),
			Description:   gofakeit.Product().Description,
			Price:         gofakeit.Product().Price,
			StockQuantity: int64(gofakeit.Int8()),
			Category:      model.Category_CATEGORY_ENGINE,
			Dimensions: &model.Dimensions{
				Length: float64(gofakeit.Int8()),
				Width:  float64(gofakeit.Int8()),
				Height: float64(gofakeit.Int8()),
				Weight: float64(gofakeit.Int8()),
			},
			Manufacturer: &model.Manufacturer{
				Name:    gofakeit.Company(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			},
			Tags: []string{"engine"},
			Metadata: map[string]*model.Value{
				"power_output":    {DoubleValue: u.ToPtr(9.5)},
				"is_experimental": {BoolValue: u.ToPtr(true)},
			},
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		}
	)

	s.partRepo.On("GetPart", s.ctx, partUuid).Return(expectedPart, nil)

	part, err := s.service.GetPart(s.ctx, partUuid)
	s.Require().NoError(err)
	s.Require().Equal(part, expectedPart)
}

func (s *ServiceSuite) TestGetPartNotFound() {
	var (
		partUuid = gofakeit.UUID()

		expectedErr = model.ErrPartDoesNotExist
	)

	s.partRepo.On("GetPart", s.ctx, partUuid).Return(model.Part{}, model.ErrPartDoesNotExist)

	part, err := s.service.GetPart(s.ctx, partUuid)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
	s.Require().Empty(part)
}
