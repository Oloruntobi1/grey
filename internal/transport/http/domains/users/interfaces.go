package users

import (
	"context"

	"github.com/Oloruntobi1/grey/internal/models"
)

type UserAdapter interface {
	CreateUser(ctx context.Context, userModel *models.User) (string, error)
}
