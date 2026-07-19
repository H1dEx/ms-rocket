package model

import "time"

type Category string

const (
	CategoryUnknown  Category = "UNKNOWN"
	CategoryEngine   Category = "ENGINE"
	CategoryFuel     Category = "FUEL"
	CategoryPorthole Category = "PORTHOLE"
	CategoryWing     Category = "WING"
)

// Dimensions размеры детали Part
type Dimensions struct {
	// length длинна детали
	Length float64
	// width ширина детали
	Width float64
	// height высота детали
	Height float64
	// weght вес детали
	Weight float64
}

// Manufacturer информация о производителе детали Part
type Manufacturer struct {
	// name название произвлдителя
	Name string
	// country страна производителя
	Country string
	// website ссылка на сайт производителя
	Website string
}

type MetadataValue struct {
	StringValue *string
	Int64Value  *int64
	DoubleValue *float64
	BoolValue   *bool
}

type Part struct {
	// uuid Уникальный идентификатор детали
	Uuid string
	// name Название детали
	Name string
	// description Описание детали
	Description string
	// price Цена за единицу
	Price uint64
	// stock_quantity Количество на складе
	StockQuantity int64
	// category Категория
	Category Category
	// dimensions Размеры детали
	Dimensions *Dimensions
	// manufacturer Информация о производителе
	Manufacturer *Manufacturer
	// tags для быстрого поиска
	Tags []string
	// metadata Гибкие метаданные
	Metadata map[string]*MetadataValue
	// created_at Дата создания
	CreatedAt *time.Time
	// updated_at Дата обновления
	UpdatedAt *time.Time
}

type PartFilter struct {
	Uuids []string
	// Список имён. Пусто — не фильтруем по имени
	Names []string
	// Список категорий. Пусто — не фильтруем по категории
	Categories []Category
	// Список стран производителей. Пусто — не фильтруем по стране
	ManufacturerCountries []string
	// Список тегов. Пусто — не фильтруем по тегам
	Tags []string
}
