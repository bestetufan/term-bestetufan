package entity

import (
	"time"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Email     string    `gorm:"size:255;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"-"`
	Roles     []*Role   `gorm:"many2many:user_roles;" json:"roles"`
	CreatedAt time.Time `gorm:"<-:create" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(email string, password string, roles []*Role) *User {
	return &User{
		Email:    email,
		Password: password,
		Roles:    roles,
	}
}

func (User) TableName() string {
	return "user"
}
