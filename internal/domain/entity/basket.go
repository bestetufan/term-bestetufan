package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Basket struct {
	ID        string        `gorm:"primary_key;" json:"id"`
	UserName  string        `json:"user_name"`
	Items     []*BasketItem `gorm:"foreignkey:BasketID;" json:"items"`
	CreatedAt time.Time     `gorm:"<-:create" json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type BasketItem struct {
	BasketID  string    `gorm:"primary_key" json:"basket_id"`
	ProductID uint32    `gorm:"primary_key" json:"product_id"`
	Product   *Product  `gorm:"foreignkey:ProductID;references:ID" json:"product"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	CreatedAt time.Time `gorm:"<-:create" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewBasket(userName string) (*Basket, error) {
	if len(userName) == 0 {
		return nil, fmt.Errorf("userName field is required")
	}
	return &Basket{
		UserName: userName,
		Items:    nil,
	}, nil
}

func NewBasketItem(id string, productId uint32, quantity int) (*BasketItem, error) {
	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than zero")
	}
	return &BasketItem{
		BasketID:  id,
		ProductID: productId,
		Quantity:  quantity,
	}, nil
}

func (Basket) TableName() string {
	return "basket"
}

func (b *Basket) BeforeCreate(db *gorm.DB) (err error) {
	b.ID = uuid.NewString()
	return
}

func (BasketItem) TableName() string {
	return "basket_item"
}

func (b *Basket) SearchItemByProductId(productId uint32) (int, *BasketItem) {
	for i, n := range b.Items {
		if n.ProductID == productId {
			return i, n
		}
	}
	return -1, nil
}
