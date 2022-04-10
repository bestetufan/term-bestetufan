package repo

import (
	"errors"

	"github.com/bestetufan/beste-store/internal/domain/entity"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) GetAll(pageIndex, pageSize int) ([]entity.Product, int) {
	var products []entity.Product
	var count int64

	r.db.Preload("Category").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&products).Count(&count)

	return products, int(count)
}

func (r *ProductRepository) GetAllByCategoryId(categoryId uint32) ([]entity.Product, int) {
	var products []entity.Product
	var count int64
	r.db.Preload("Category").Find(&products).Where(&entity.Product{CategoryID: categoryId}).Count(&count)

	return products, int(count)
}

func (r *ProductRepository) GetById(id uint32) *entity.Product {
	var product entity.Product
	result := r.db.Preload("Category").First(&product, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return &product
}

func (r *ProductRepository) GetBySKU(sku string) *entity.Product {
	var product entity.Product
	result := r.db.Preload("Category").Where(&entity.Product{Sku: sku}).First(&product)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return &product
}

func (r *ProductRepository) Search(query string) []entity.Product {
	var products []entity.Product
	r.db.Preload("Category").Where("name LIKE ? OR sku LIKE ?", "%"+query+"%", "%"+query+"%").Find(&products)

	return products
}

func (r *ProductRepository) Create(c *entity.Product) error {
	result := r.db.Create(&c)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *ProductRepository) Update(c *entity.Product) error {
	result := r.db.Save(&c)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *ProductRepository) DeleteById(id uint32) error {
	result := r.db.Delete(&entity.Product{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
