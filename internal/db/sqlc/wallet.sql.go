// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: wallet.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createWallet = `-- name: CreateWallet :one
INSERT INTO wallets(
user_id,
balance
) VALUES (
    $1, $2
) RETURNING id, user_id, balance, created_at, updated_at, deleted_at, is_deleted
`

type CreateWalletParams struct {
	UserID  uuid.UUID      `json:"user_id"`
	Balance pgtype.Numeric `json:"balance"`
}

func (q *Queries) CreateWallet(ctx context.Context, arg CreateWalletParams) (Wallet, error) {
	row := q.db.QueryRow(ctx, createWallet, arg.UserID, arg.Balance)
	var i Wallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Balance,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.IsDeleted,
	)
	return i, err
}