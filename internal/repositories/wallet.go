package repositories

import (
	"context"
	"errors"
	"fmt"

	db "github.com/Oloruntobi1/grey/internal/db/sqlc"
	"github.com/Oloruntobi1/grey/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var ErrWalletNotFound = errors.New("wallet not found")

type WalletRepository struct {
	db     db.Querier
	tracer trace.Tracer
}

func NewWalletRepository(db db.Querier) *WalletRepository {
	return &WalletRepository{
		db:     db,
		tracer: otel.Tracer("walletRepository"),
	}
}

func (r *WalletRepository) CreateWallet(ctx context.Context, walletModel *models.Wallet) (string, error) {
	ctx, span := r.tracer.Start(ctx, "walletRepo.Create")
	defer span.End()

	dbWallet, err := r.toDb(walletModel)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return "", fmt.Errorf("mapping failed: err %v", err)
	}

	walletDB, err := r.db.CreateWallet(ctx, dbWallet)
	if err != nil {
		err = fmt.Errorf("failed to add wallet in db: %w", err)
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return "", err
	}

	return walletDB.ID.String(), nil
}

func (u *WalletRepository) toDb(walletModel *models.Wallet) (db.CreateWalletParams, error) {
	wallet := db.CreateWalletParams{
		UserID: uuid.MustParse(walletModel.UserID),
		Balance: pgtype.Numeric{
			Int:   walletModel.Balance.BigInt(),
			Exp:   walletModel.Balance.Exponent(),
			Valid: true,
		},
	}

	return wallet, nil
}
