package converter

import (
	"github.com/Medveddo/rocket-science/order/internal/model"
	inventoryV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/inventory/v1"
)

func PartFromProto(protoPart *inventoryV1.Part) model.Part {
	return model.Part{
		UUID:          protoPart.GetUuid(),
		Name:          protoPart.GetName(),
		Description:   protoPart.GetDescription(),
		Price:         protoPart.GetPrice(),
		StockQuantity: protoPart.GetStockQuantity(),
		Category:      CategoryFromProto(protoPart.GetCategory()),
		Dimensions:    DimensionsFromProto(protoPart.GetDimensions()),
		Manufacturer:  ManufacturerFromProto(protoPart.GetManufacturer()),
		Tags:          protoPart.GetTags(),
		Metadata:      MetadataFromProto(protoPart.GetMetadata()),
		CreatedAt:     protoPart.GetCreatedAt().AsTime(),
		UpdatedAt:     protoPart.GetUpdatedAt().AsTime(),
	}
}

func PartsFromProto(protoParts []*inventoryV1.Part) []model.Part {
	parts := make([]model.Part, len(protoParts))
	for i, protoPart := range protoParts {
		parts[i] = PartFromProto(protoPart)
	}
	return parts
}

func CategoryFromProto(protoCategory inventoryV1.Category) model.Category {
	switch protoCategory {
	case inventoryV1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventoryV1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventoryV1.Category_CATEGORY_WING:
		return model.CategoryWing
	case inventoryV1.Category_CATEGORY_SHIELD:
		return model.CategoryShield
	default:
		return model.CategoryUnspecified
	}
}

func CategoryToProto(category model.Category) inventoryV1.Category {
	switch category {
	case model.CategoryEngine:
		return inventoryV1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryV1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryV1.Category_CATEGORY_WING
	case model.CategoryShield:
		return inventoryV1.Category_CATEGORY_SHIELD
	default:
		return inventoryV1.Category_CATEGORY_UNSPECIFIED
	}
}

func DimensionsFromProto(protoDimensions *inventoryV1.Dimensions) *model.Dimensions {
	if protoDimensions == nil {
		return nil
	}
	return &model.Dimensions{
		Length: protoDimensions.GetLength(),
		Width:  protoDimensions.GetWidth(),
		Height: protoDimensions.GetHeight(),
		Weight: protoDimensions.GetWeight(),
	}
}

func ManufacturerFromProto(protoManufacturer *inventoryV1.Manufacturer) *model.Manufacturer {
	if protoManufacturer == nil {
		return nil
	}
	return &model.Manufacturer{
		Name:    protoManufacturer.GetName(),
		Country: protoManufacturer.GetCountry(),
		Website: protoManufacturer.GetWebsite(),
	}
}

func MetadataFromProto(protoMetadata map[string]*inventoryV1.Value) map[string]model.Value {
	if protoMetadata == nil {
		return nil
	}
	metadata := make(map[string]model.Value)
	for key, protoValue := range protoMetadata {
		metadata[key] = ValueFromProto(protoValue)
	}
	return metadata
}

func ValueFromProto(protoValue *inventoryV1.Value) model.Value {
	if protoValue == nil {
		return model.Value{}
	}

	switch kind := protoValue.GetKind().(type) {
	case *inventoryV1.Value_StringValue:
		return model.Value{StringValue: &kind.StringValue}
	case *inventoryV1.Value_Int64Value:
		return model.Value{Int64Value: &kind.Int64Value}
	case *inventoryV1.Value_DoubleValue:
		return model.Value{DoubleValue: &kind.DoubleValue}
	case *inventoryV1.Value_BoolValue:
		return model.Value{BoolValue: &kind.BoolValue}
	default:
		return model.Value{}
	}
}
