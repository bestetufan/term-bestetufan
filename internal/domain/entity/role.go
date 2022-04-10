package entity

import (
	"time"
)

type Role struct {
	Name      string    `gorm:"primary_key;size:100;not null" json:"name"`
	CreatedAt time.Time `gorm:"<-:create" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewRole(name string) *Role {
	return &Role{
		Name: name,
	}
}

func (Role) TableName() string {
	return "role"
}
