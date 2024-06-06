package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Oloruntobi1/grey/internal/models"
	"github.com/Oloruntobi1/grey/internal/transport/http/domains/wallets"
	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type WalletHandler struct {
	svc    wallets.WalletService
	logger *slog.Logger

	tracer trace.Tracer
}

func NewWalletHandler(svc wallets.WalletService, logger *slog.Logger) *WalletHandler {
	return &WalletHandler{
		svc:    svc,
		logger: logger,
		tracer: otel.Tracer("userHandler"),
	}
}

type createWalletRequest struct {
	UserID         string          `json:"user_id"`
	InitialBalance decimal.Decimal `json:"initial_balance"`
}

func (h *WalletHandler) CreateWalletHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := h.tracer.Start(r.Context(), "createWalletHandler")
		defer span.End()
		var request createWalletRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := h.svc.CreateWallet(ctx, &models.Wallet{
			UserID:  request.UserID,
			Balance: request.InitialBalance,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := ResponseWithID(ctx, id)
		responseJSON, err := json.Marshal(response)
		if err != nil {
			h.logger.ErrorContext(
				r.Context(),
				"failed_to_marshal_response",
				slog.Any("err", err),
			)
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write(responseJSON); err != nil {
			h.logger.ErrorContext(
				r.Context(),
				"failed_to_write_response",
				slog.Any("err", err),
			)
		}
	}
}
