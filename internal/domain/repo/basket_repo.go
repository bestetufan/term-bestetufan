package repo

import (
	"errors"

	"github.com/bestetufan/beste-store/internal/domain/entity"
	"gorm.io/gorm"
)

type BasketRepository struct {
	db *gorm.DB
}

func NewBasketRepository(db *gorm.DB) *BasketRepository {
	return &BasketRepository{
		db: db,
	}
}

func (r *BasketRepository) Get(userName string) *entity.Basket {
	var basket entity.Basket
	result := r.db.Preload("Items").Preload("Items.Product").Preload("Items.Product.Category").FirstOrCreate(&basket, entity.Basket{UserName: userName})

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return &basket
}

func (r *BasketRepository) Create(c *entity.Basket) error {
	result := r.db.Create(&c)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *BasketRepository) Update(c *entity.Basket) error {
	result := r.db.Save(&c)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *BasketRepository) Delete(userName string) error {
	result := r.db.Where(&entity.Basket{UserName: userName}).Delete(&entity.Basket{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *BasketRepository) AddItem(c *entity.BasketItem) error {
	result := r.db.Create(&c)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *BasketRepository) UpdateItem(c *entity.BasketItem) error {
	result := r.db.Save(&c)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *BasketRepository) DeleteItem(c *entity.BasketItem) error {
	result := r.db.Delete(&c)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *BasketRepository) DeleteItemsByBasketId(basketId string) error {
	result := r.db.Where(&entity.BasketItem{BasketID: basketId}).Delete(&entity.BasketItem{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
