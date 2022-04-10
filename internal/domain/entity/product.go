package entity

import (
	"time"
)

type Product struct {
	ID         uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name       string    `gorm:"size:255;not null;" json:"name"`
	Sku        string    `gorm:"size:100;not null;unique" json:"sku"`
	UnitPrice  float64   `json:"unit_price"`
	Quantity   int       `json:"quantity"`
	CategoryID uint32    `json:"category_id"`
	Category   Category  `json:"category"`
	CreatedAt  time.Time `gorm:"<-:create" json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewProduct(name string, sku string, unitPrice float64, quantity int, categoryID uint32) *Product {
	return &Product{
		Name:       name,
		Sku:        sku,
		UnitPrice:  unitPrice,
		Quantity:   quantity,
		CategoryID: categoryID,
	}
}

func (Product) TableName() string {
	return "product"
}
