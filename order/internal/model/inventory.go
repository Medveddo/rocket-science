package model

import (
	"time"
)

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]Value
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Category string

const (
	CategoryUnspecified Category = "CATEGORY_UNSPECIFIED"
	CategoryEngine      Category = "CATEGORY_ENGINE"
	CategoryFuel        Category = "CATEGORY_FUEL"
	CategoryPorthole    Category = "CATEGORY_PORTHOLE"
	CategoryWing        Category = "CATEGORY_WING"
	CategoryShield      Category = "CATEGORY_SHIELD"
)

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Value struct {
	StringValue *string
	Int64Value  *int64
	DoubleValue *float64
	BoolValue   *bool
}

type PartsFilter struct {
	UUIDs                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
