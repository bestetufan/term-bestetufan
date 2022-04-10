package repo

import (
	"errors"

	"github.com/bestetufan/beste-store/internal/domain/entity"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) GetAll(pageIndex, pageSize int) ([]entity.Category, int) {
	var categories []entity.Category
	var count int64

	r.db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&categories).Count(&count)

	return categories, int(count)
}

func (r *CategoryRepository) GetAllActives(pageIndex, pageSize int) ([]entity.Category, int) {
	var categories []entity.Category
	var count int64

	r.db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Where(&entity.Category{IsActive: true}).Find(&categories).Count(&count)

	return categories, int(count)
}

func (r *CategoryRepository) GetById(id uint32) *entity.Category {
	var product entity.Category
	result := r.db.First(&product, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return &product
}

func (r *CategoryRepository) GetByName(name string) *entity.Category {
	var category entity.Category
	result := r.db.Where(entity.Category{Name: name}).Attrs(entity.Category{}).First(&category)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return &category
}

func (r *CategoryRepository) Create(c *entity.Category) error {
	result := r.db.Create(&c)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *CategoryRepository) Update(c entity.Category) error {
	result := r.db.Save(&c)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *CategoryRepository) DeleteById(id uint32) error {
	result := r.db.Delete(&entity.Category{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
