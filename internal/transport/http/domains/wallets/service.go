package wallets

import (
	"context"

	"github.com/Oloruntobi1/grey/internal/models"
)

type WalletService struct {
	userRepo WalletAdapter
}

func NewWalletService(userRepo WalletAdapter) *WalletService {
	return &WalletService{userRepo: userRepo}
}

func (s *WalletService) CreateWallet(ctx context.Context, user *models.Wallet) (string, error) {
	return s.userRepo.CreateWallet(ctx, user)
}
