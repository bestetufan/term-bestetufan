package repo

import (
	"errors"

	"github.com/bestetufan/beste-store/internal/domain/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetAll() ([]entity.User, int) {
	var users []entity.User
	var count int64

	r.db.Find(&users).Count(&count)

	return users, int(count)
}

func (r *UserRepository) GetById(id int) *entity.User {
	var user entity.User
	result := r.db.First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return &user
}

func (r *UserRepository) GetByIdWithRoles(id int) *entity.User {
	var user entity.User
	result := r.db.Preload("Roles").First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return &user
}

func (r *UserRepository) GetByEmail(email string) *entity.User {
	var user entity.User
	result := r.db.Where(entity.User{Email: email}).Attrs(entity.User{}).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return &user
}

func (r *UserRepository) GetByEmailPassword(email string, password string) *entity.User {
	var user entity.User
	result := r.db.Where(entity.User{Email: email, Password: password}).Attrs(entity.User{}).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return &user
}

func (r *UserRepository) Create(c *entity.User) error {
	result := r.db.Create(&c)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *UserRepository) Update(c *entity.User) error {
	result := r.db.Save(&c)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *UserRepository) DeleteById(id int) error {
	result := r.db.Delete(&entity.User{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
