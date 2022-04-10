package entity

import (
	"time"
)

type Category struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:50;not null;" json:"name"`
	IsActive  bool      `gorm:"not null;" json:"is_active"`
	Products  []Product `gorm:"foreignkey:CategoryID" json:"products"`
	CreatedAt time.Time `gorm:"<-:create" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCategory(name string, active bool) *Category {
	return &Category{
		Name:     name,
		IsActive: active,
	}
}

func (Category) TableName() string {
	return "category"
}
