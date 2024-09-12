package services

import (
	"p2pbot/internal/db/models"
	"p2pbot/internal/db/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) CreateUser(user *models.User) (int, error) {
	return s.repo.Save(user)
}

func (s *UserService) GetUserByChatID(id int64) (*models.User, error) {
	return s.repo.GetByChatID(id)
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.repo.GetByEmail(email)
}
