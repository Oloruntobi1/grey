package repositories

import (
	"context"
	"errors"
	"fmt"

	db "github.com/Oloruntobi1/grey/internal/db/sqlc"
	"github.com/Oloruntobi1/grey/internal/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserRepository struct {
	db     db.Querier
	tracer trace.Tracer
}

func NewUserRepository(db db.Querier) *UserRepository {
	return &UserRepository{
		db:     db,
		tracer: otel.Tracer("userRepository"),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, userModel *models.User) (string, error) {
	ctx, span := r.tracer.Start(ctx, "userRepo.Create")
	defer span.End()

	dbUser, err := r.toDb(userModel)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return "", fmt.Errorf("mapping failed: err %v", err)
	}

	userDB, err := r.db.CreateUser(ctx, dbUser)
	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.UniqueViolation {
			err = ErrUserAlreadyExists
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err, trace.WithAttributes(attribute.String("email", userModel.Email)))
			return "", err
		}
		err = fmt.Errorf("failed to add user in db: %w", err)
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return "", err
	}

	return userDB.ID.String(), nil
}

func (u *UserRepository) toDb(userModel *models.User) (db.CreateUserParams, error) {
	user := db.CreateUserParams{
		Name:  userModel.Name,
		Email: userModel.Email,
	}

	return user, nil
}
