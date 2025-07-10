package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Medveddo/rocket-science/inventory/internal/model"
	u "github.com/Medveddo/rocket-science/inventory/internal/utils"
	inventoryV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/inventory/v1"
)

// --- gRPC → Service Model ---

func PartFromProto(p *inventoryV1.Part) model.Part {
	return model.Part{
		UUID:          p.Uuid,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      model.Category(p.Category),
		Dimensions:    DimensionsFromProto(p.Dimensions),
		Manufacturer:  ManufacturerFromProto(p.Manufacturer),
		Tags:          p.Tags,
		Metadata:      MetadataFromProto(p.Metadata),
		CreatedAt:     p.CreatedAt.AsTime(),
		UpdatedAt:     p.UpdatedAt.AsTime(),
	}
}

func DimensionsFromProto(d *inventoryV1.Dimensions) *model.Dimensions {
	if d == nil {
		return nil
	}
	return &model.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func ManufacturerFromProto(m *inventoryV1.Manufacturer) *model.Manufacturer {
	if m == nil {
		return nil
	}
	return &model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func MetadataFromProto(meta map[string]*inventoryV1.Value) map[string]*model.Value {
	result := make(map[string]*model.Value, len(meta))
	for k, v := range meta {
		result[k] = ValueFromProto(v)
	}
	return result
}

func ValueFromProto(v *inventoryV1.Value) *model.Value {
	if v == nil {
		return nil
	}
	mv := &model.Value{}
	switch x := v.Kind.(type) {
	case *inventoryV1.Value_DoubleValue:
		mv.DoubleValue = u.ToPtr(x.DoubleValue)
	case *inventoryV1.Value_Int64Value:
		mv.Int64Value = u.ToPtr(x.Int64Value)
	case *inventoryV1.Value_BoolValue:
		mv.BoolValue = u.ToPtr(x.BoolValue)
	case *inventoryV1.Value_StringValue:
		mv.StringValue = u.ToPtr(x.StringValue)
	}
	return mv
}

// --- Service Model → gRPC ---

func PartToProto(p model.Part) *inventoryV1.Part {
	return &inventoryV1.Part{
		Uuid:          p.UUID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      inventoryV1.Category(p.Category),
		Dimensions:    DimensionsToProto(p.Dimensions),
		Manufacturer:  ManufacturerToProto(p.Manufacturer),
		Tags:          p.Tags,
		Metadata:      MetadataToProto(p.Metadata),
		CreatedAt:     timestamppb.New(p.CreatedAt),
		UpdatedAt:     timestamppb.New(p.UpdatedAt),
	}
}

func DimensionsToProto(d *model.Dimensions) *inventoryV1.Dimensions {
	if d == nil {
		return nil
	}
	return &inventoryV1.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func ManufacturerToProto(m *model.Manufacturer) *inventoryV1.Manufacturer {
	if m == nil {
		return nil
	}
	return &inventoryV1.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func MetadataToProto(meta map[string]*model.Value) map[string]*inventoryV1.Value {
	result := make(map[string]*inventoryV1.Value, len(meta))
	for k, v := range meta {
		result[k] = ValueToProto(v)
	}
	return result
}

func ValueToProto(v *model.Value) *inventoryV1.Value {
	if v == nil {
		return nil
	}
	switch {
	case v.DoubleValue != nil:
		return &inventoryV1.Value{Kind: &inventoryV1.Value_DoubleValue{DoubleValue: *v.DoubleValue}}
	case v.Int64Value != nil:
		return &inventoryV1.Value{Kind: &inventoryV1.Value_Int64Value{Int64Value: *v.Int64Value}}
	case v.BoolValue != nil:
		return &inventoryV1.Value{Kind: &inventoryV1.Value_BoolValue{BoolValue: *v.BoolValue}}
	case v.StringValue != nil:
		return &inventoryV1.Value{Kind: &inventoryV1.Value_StringValue{StringValue: *v.StringValue}}
	default:
		return &inventoryV1.Value{}
	}
}

// ListPartsFilterFromProto converts a proto filter to the internal service filter.
func ListPartsFilterFromProto(f *inventoryV1.PartsFilter) *model.PartsFilter {
	if f == nil {
		return nil
	}
	categories := make([]model.Category, 0, len(f.Categories))
	for _, c := range f.Categories {
		categories = append(categories, model.Category(c))
	}
	return &model.PartsFilter{
		UUIDs:                 f.Uuids,
		Names:                 f.Names,
		Categories:            categories,
		ManufacturerCountries: f.ManufacturerCountries,
		Tags:                  f.Tags,
	}
}
