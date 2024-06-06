package db

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ForeignKeyViolation       = "23503"
	UniqueViolation           = "23505"
	InvalidTextRepresentation = "22P02"
	CheckViolation            = "23514" // when the CHECK constraint have been used in the DDL
)

var ErrRecordNotFound = pgx.ErrNoRows

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

func ErrorMessage(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Message
	}
	return ""
}
