package repo

import (
	"github.com/bestetufan/beste-store/internal/domain/entity"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Get(userName string, id string) *entity.Order {
	var order entity.Order
	r.db.Where(&entity.Order{ID: id, UserName: userName}).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		First(&order)

	return &order
}

func (r *OrderRepository) GetAll(userName string) []entity.Order {
	var orders []entity.Order
	r.db.Where(&entity.Order{UserName: userName}).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Find(&orders)

	return orders
}

func (r *OrderRepository) Create(c *entity.Order) error {
	result := r.db.Create(&c)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *OrderRepository) Update(c *entity.Order) error {
	result := r.db.Save(&c)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
