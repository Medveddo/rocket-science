package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

type Value struct {
	DoubleValue *float64 `bson:"double_value,omitempty"`
	Int64Value  *int64   `bson:"int64_value,omitempty"`
	BoolValue   *bool    `bson:"bool_value,omitempty"`
	StringValue *string  `bson:"string_value,omitempty"`
}

type Part struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string             `bson:"name"`
	Description   string             `bson:"description"`
	Price         float64            `bson:"price"`
	StockQuantity int64              `bson:"stock_quantity"`
	Category      Category           `bson:"category"`
	Dimensions    *Dimensions        `bson:"dimensions"`
	Manufacturer  *Manufacturer      `bson:"manufacturer"`
	Tags          []string           `bson:"tags"`
	Metadata      map[string]*Value  `bson:"metadata"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}

// Категория детали
type Category int32

const (
	Category_CATEGORY_UNSPECIFIED Category = 0
	Category_CATEGORY_ENGINE      Category = 1
	Category_CATEGORY_FUEL        Category = 2
	Category_CATEGORY_PORTHOLE    Category = 3
	Category_CATEGORY_WING        Category = 4
	Category_CATEGORY_SHIELD      Category = 5
)
