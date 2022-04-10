package service

import (
	"errors"

	"github.com/bestetufan/beste-store/internal/domain/entity"
	"github.com/bestetufan/beste-store/internal/domain/repo"
)

type UserService struct {
	repo repo.UserRepository
}

func NewUserService(r repo.UserRepository) *UserService {
	return &UserService{
		repo: r,
	}
}

func (s *UserService) GetUser(email string, password string) *entity.User {
	return s.repo.GetByEmailPassword(email, password)
}

func (s *UserService) UserHasRole(id int, role string) bool {
	user := s.repo.GetByIdWithRoles(id)
	if user == nil {
		return false
	}

	for i := range user.Roles {
		if user.Roles[i].Name == role {
			return true
		}
	}

	return false
}

func (s *UserService) CreateUser(user *entity.User) error {
	userExists := s.repo.GetByEmail(user.Email)
	if userExists != nil {
		return errors.New("user with same email already exist in database")
	}

	err := s.repo.Create(user)
	if err != nil {
		return errors.New("an unknown error occurred during operation")
	}

	return nil
}
