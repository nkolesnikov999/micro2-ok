package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Part struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Uuid          string             `bson:"uuid"`
	Name          string             `bson:"name"`
	Description   string             `bson:"description"`
	Price         float64            `bson:"price"`
	StockQuantity int64              `bson:"stock_quantity"`
	Category      Category           `bson:"category"`
	Dimensions    *Dimensions        `bson:"dimensions,omitempty"`
	Manufacturer  *Manufacturer      `bson:"manufacturer,omitempty"`
	Tags          []string           `bson:"tags"`
	Metadata      map[string]*Value  `bson:"metadata"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}

type Category int32

const (
	CategoryUnspecified Category = 0
	CategoryEngine      Category = 1
	CategoryFuel        Category = 2
	CategoryPorthole    Category = 3
	CategoryWing        Category = 4
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
	StringValue string  `bson:"string_value,omitempty"`
	Int64Value  int64   `bson:"int64_value,omitempty"`
	DoubleValue float64 `bson:"double_value,omitempty"`
	BoolValue   bool    `bson:"bool_value,omitempty"`
}
