package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Oloruntobi1/grey/internal/models"
	"github.com/Oloruntobi1/grey/internal/transport/http/domains/users"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type UserHandler struct {
	svc    users.UserService
	logger *slog.Logger

	tracer trace.Tracer
}

func NewUserHandler(svc users.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		svc:    svc,
		logger: logger,
		tracer: otel.Tracer("userHandler"),
	}
}

type createUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *UserHandler) CreateUserHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := h.tracer.Start(r.Context(), "createUserHandler")
		defer span.End()
		var request createUserRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := h.svc.CreateUser(ctx, &models.User{
			Name:  request.Name,
			Email: request.Email,
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
