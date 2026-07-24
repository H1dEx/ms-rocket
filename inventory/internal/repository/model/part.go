package model

import (
	"time"
)

// Dimensions размеры детали Part
type Dimensions struct {
	// length длинна детали
	Length float64 `bson:"length"`
	// width ширина детали
	Width float64 `bson:"width"`
	// height высота детали
	Height float64 `bson:"height"`
	// weght вес детали
	Weight float64 `bson:"weight"`
}

// Manufacturer информация о производителе детали Part
type Manufacturer struct {
	// name название произвлдителя
	Name string `bson:"name"`
	// country страна производителя
	Country string `bson:"country"`
	// website ссылка на сайт производителя
	Website string `bson:"website"`
}

type MetadataValue struct {
	StringValue *string
	Int64Value  *int64
	DoubleValue *float64
	BoolValue   *bool
}

type Part struct {
	UUID          string                   `bson:"_id"`
	Name          string                   `bson:"name"`
	Description   string                   `bson:"description,omitempty"`
	Price         uint64                   `bson:"price"`
	StockQuantity int64                    `bson:"quantity"`
	Category      string                   `bson:"category"`
	Dimensions    Dimensions               `bson:"dimensions,omitempty"`
	Manufacturer  Manufacturer             `bson:"manufacturer,omitempty"`
	Tags          []string                 `bson:"tags,omitempty"`
	Metadata      map[string]MetadataValue `bson:"metadata,omitempty"`
	CreatedAt     time.Time                `bson:"created_at"`
	UpdatedAt     *time.Time               `bson:"updated_at,omitempty"`
}

type PartFilter struct {
	Uuids []string
	// Список имён. Пусто — не фильтруем по имени
	Names []string
	// Список категорий. Пусто — не фильтруем по категории
	Categories []string
	// Список стран производителей. Пусто — не фильтруем по стране
	ManufacturerCountries []string
	// Список тегов. Пусто — не фильтруем по тегам
	Tags []string
}
