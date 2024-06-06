package users

import (
	"context"

	"github.com/Oloruntobi1/grey/internal/models"
)

type UserService struct {
	userRepo UserAdapter
}

func NewUserService(userRepo UserAdapter) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (string, error) {
	return s.userRepo.CreateUser(ctx, user)
}
