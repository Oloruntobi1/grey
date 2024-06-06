// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users(
name,
email
) VALUES (
    $1, $2
) RETURNING id, name, email, created_at, updated_at, deleted_at, is_deleted
`

type CreateUserParams struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Name, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.IsDeleted,
	)
	return i, err
}
