package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID          string       `gorm:"primary_key;" json:"id"`
	UserName    string       `gorm:"size:255;not null;" json:"user_name"`
	Status      string       `gorm:"size:100;not null;" json:"status"`
	Name        string       `gorm:"size:100;not null;" json:"name"`
	Address     string       `gorm:"size:255;not null;" json:"address"`
	PhoneNumber string       `gorm:"size:100;not null;" json:"phone_number"`
	CardNumber  string       `gorm:"size:255;not null;" json:"card_number"`
	CardExp     string       `gorm:"size:100;not null;" json:"card_exp"`
	CardCVV     int          `gorm:"not null;" json:"card_cvv"`
	Items       []*OrderItem `gorm:"foreignkey:OrderID;" json:"items"`
	CreatedAt   time.Time    `gorm:"<-:create" json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type OrderItem struct {
	OrderID   string    `gorm:"primary_key" json:"order_id"`
	ProductID uint32    `gorm:"primary_key" json:"product_id"`
	Product   *Product  `gorm:"foreignkey:ProductID;references:ID" json:"product"`
	CreatedAt time.Time `gorm:"<-:create" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewOrder(userName string, name string, address string, phoneNumber string,
	cardNumber string, cardExp string, cardCVV int) (*Order, error) {
	if len(userName) == 0 {
		return nil, fmt.Errorf("userName field is required")
	}
	return &Order{
		UserName:    userName,
		Status:      "incomplete",
		Name:        name,
		Address:     address,
		PhoneNumber: phoneNumber,
		CardNumber:  cardNumber,
		CardExp:     cardExp,
		CardCVV:     cardCVV,
		Items:       nil,
	}, nil
}

func NewOrderItem(orderId string, productId uint32) (*OrderItem, error) {
	return &OrderItem{
		OrderID:   orderId,
		ProductID: productId,
	}, nil
}

func (Order) TableName() string {
	return "order"
}

func (b *Order) BeforeCreate(db *gorm.DB) (err error) {
	b.ID = uuid.NewString()
	return
}

func (OrderItem) TableName() string {
	return "order_item"
}

func (b *Order) AddItem(item *OrderItem) error {
	b.Items = append(b.Items, item)
	return nil
}
