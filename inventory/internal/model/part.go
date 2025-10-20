package model

import "time"

type Part struct {
	Uuid string
	// Название детали
	Name string
	// Описание детали
	Description string
	// Цена за единицу
	Price float64
	// Количество на складе
	StockQuantity int64
	// Категория
	Category Category
	// Размеры детали
	Dimensions *Dimensions
	// Информация о производителе
	Manufacturer *Manufacturer
	// Теги для быстрого поиска
	Tags []string
	// Гибкие метаданные
	Metadata map[string]*Value
	// Дата создания
	CreatedAt time.Time
	// Дата обновления
	UpdatedAt time.Time
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
	StringValue string
	Int64Value  int64
	DoubleValue float64
	BoolValue   bool
}

type PartsFilter struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
