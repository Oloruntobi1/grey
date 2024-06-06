package handlers

import (
	"context"
	"net/http"
)

func SetupRouter(
	ctx context.Context,
	userHandler UserHandler,
	walletService WalletHandler,
	// transactionService *service.TransactionService,
) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/create-user", userHandler.CreateUserHandler(ctx))
	mux.HandleFunc("/api/create-wallet", walletService.CreateWalletHandler(ctx))
	// mux.HandleFunc("/process-transaction", processTransactionHandler(transactionService))
	return mux
}
