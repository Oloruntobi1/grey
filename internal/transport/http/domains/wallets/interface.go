package wallets

import (
	"context"

	"github.com/Oloruntobi1/grey/internal/models"
)

type WalletAdapter interface {
	CreateWallet(ctx context.Context, walletModel *models.Wallet) (string, error)
}
