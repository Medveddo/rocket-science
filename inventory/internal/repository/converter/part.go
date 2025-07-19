package converter

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Medveddo/rocket-science/inventory/internal/model"
	repoModel "github.com/Medveddo/rocket-science/inventory/internal/repository/model"
)

// Converts domain Part to repository Part
// nolint
func PartToRepoPart(p model.Part) (repoModel.Part, error) {
	id_, err := primitive.ObjectIDFromHex(p.UUID)
	if err != nil {
		return repoModel.Part{}, err
	}

	repoPart := repoModel.Part{
		ID:            id_,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      repoModel.Category(p.Category),
		Tags:          p.Tags,
		Metadata:      make(map[string]*repoModel.Value),
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
	if p.Dimensions != nil {
		repoPart.Dimensions = &repoModel.Dimensions{
			Length: p.Dimensions.Length,
			Width:  p.Dimensions.Width,
			Height: p.Dimensions.Height,
			Weight: p.Dimensions.Weight,
		}
	}
	if p.Manufacturer != nil {
		repoPart.Manufacturer = &repoModel.Manufacturer{
			Name:    p.Manufacturer.Name,
			Country: p.Manufacturer.Country,
			Website: p.Manufacturer.Website,
		}
	}
	for k, v := range p.Metadata {
		value := PartValueToRepoValue(*v)
		repoPart.Metadata[k] = &value
	}
	return repoPart, nil
}

// Converts repository Part to domain Part
// nolint
func RepoPartToPart(r repoModel.Part) model.Part {
	part := model.Part{
		UUID:          r.ID.Hex(),
		Name:          r.Name,
		Description:   r.Description,
		Price:         r.Price,
		StockQuantity: r.StockQuantity,
		Category:      model.Category(r.Category),
		Tags:          r.Tags,
		Metadata:      make(map[string]*model.Value),
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
	if r.Dimensions != nil {
		part.Dimensions = &model.Dimensions{
			Length: r.Dimensions.Length,
			Width:  r.Dimensions.Width,
			Height: r.Dimensions.Height,
			Weight: r.Dimensions.Weight,
		}
	}
	if r.Manufacturer != nil {
		part.Manufacturer = &model.Manufacturer{
			Name:    r.Manufacturer.Name,
			Country: r.Manufacturer.Country,
			Website: r.Manufacturer.Website,
		}
	}
	for k, v := range r.Metadata {
		value := RepoValueToPartValue(*v)
		part.Metadata[k] = &value
	}
	return part
}

// Value converters
func PartValueToRepoValue(v model.Value) repoModel.Value {
	return repoModel.Value{
		DoubleValue: v.DoubleValue,
		Int64Value:  v.Int64Value,
		BoolValue:   v.BoolValue,
		StringValue: v.StringValue,
	}
}

func RepoValueToPartValue(v repoModel.Value) model.Value {
	return model.Value{
		DoubleValue: v.DoubleValue,
		Int64Value:  v.Int64Value,
		BoolValue:   v.BoolValue,
		StringValue: v.StringValue,
	}
}
