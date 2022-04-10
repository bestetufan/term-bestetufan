package migration

import (
	"fmt"

	"github.com/bestetufan/beste-store/internal/domain/entity"
	"gorm.io/gorm"
)

var roles = []entity.Role{
	{Name: "admin"},
	{Name: "customer"},
}

var users = []entity.User{
	{Email: "admin@bestestore.com", Password: "admin1234", Roles: []*entity.Role{&roles[0]}},
}

var categories = []entity.Category{
	{Name: "Cep Telefonu", IsActive: true},
}

func Execute(db *gorm.DB) error {
	// Check if migration done
	tables, err := db.Migrator().GetTables()
	if err != nil {
		return fmt.Errorf("seeder - Load - db.Migrator.GetTables: %w", err)
	}

	if len(tables) == 0 {
		// Auto create tables
		err = db.AutoMigrate(
			&entity.User{},
			&entity.Role{},
			&entity.Basket{},
			&entity.BasketItem{},
			&entity.Category{},
			&entity.Product{},
			&entity.Order{},
			&entity.OrderItem{},
		)

		if err != nil {
			return fmt.Errorf("seeder - Load - db.AutoMigrate: %w", err)
		}

		for i := range categories {
			err := db.Model(&entity.Category{}).Create(&categories[i]).Error
			if err != nil {
				return fmt.Errorf("seeder - Load - Model(Category).Create: %w", err)
			}
		}

		for i := range roles {
			err := db.Model(&entity.Role{}).Create(&roles[i]).Error
			if err != nil {
				return fmt.Errorf("seeder - Load - Model(Role).Create: %w", err)
			}
		}

		for i := range users {
			err := db.Model(&entity.User{}).Create(&users[i]).Error
			if err != nil {
				return fmt.Errorf("seeder - Load - Model(User).Create: %w", err)
			}
		}
	}

	return nil
}
