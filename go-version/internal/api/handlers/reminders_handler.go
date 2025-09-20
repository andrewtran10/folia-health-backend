package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"go-version/internal/api/middleware"
	"go-version/internal/api/repository"
	"go-version/internal/api/transport"

	"github.com/go-chi/chi/v5"
)

type ReminderHandler struct {
	repo *repository.ReminderRepository
}

func NewReminderHandler(repo *repository.ReminderRepository) (*ReminderHandler, error) {
	return &ReminderHandler{repo: repo}, nil
}

func (h *ReminderHandler) RegisterRoutes(router chi.Router) {
	authMw, err := middleware.AuthMiddleware(context.Background())
	if err != nil {
		panic(err)
	}

	h.registerPublicRoutes(router)
	h.registerProtectedRoutes(router, authMw)
}

func (h *ReminderHandler) registerPublicRoutes(router chi.Router) {
	// No public routes for reminders
}

func (h *ReminderHandler) registerProtectedRoutes(router chi.Router, authMw func(http.Handler) http.Handler) {
	router.Route("/reminders", func(r chi.Router) {
		r.Use(authMw)
		r.Post("/", h.handleCreateReminder)
		r.Get("/", h.handleListReminders)
		r.Patch("/{reminderId}", h.handleUpdateReminder)
		r.Delete("/{reminderId}", h.handleDeleteReminder)
	})
}

func (h *ReminderHandler) handleListReminders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reminderListRequest transport.ReminderListRequest
	if err := transport.ParseRequest(r, &reminderListRequest); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	reminders, err := h.repo.ListReminders(ctx, reminderListRequest.ToDomain())
	if err != nil {
		http.Error(w, "Failed to fetch reminders", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminders)

}

func (h *ReminderHandler) handleCreateReminder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req transport.ReminderCreateRequest
	if err := transport.ParseRequest(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	reminder, err := h.repo.CreateReminder(ctx, req.ToDomain())
	if err != nil {
		var invalidRRuleErr *repository.ErrInvalidRRule
		if errors.As(err, &invalidRRuleErr) {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeJSONError(w, http.StatusInternalServerError, "Unexpected error occurred")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reminder)
}

func (h *ReminderHandler) handleUpdateReminder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reminderUpdateRequest transport.ReminderUpdateRequest
	if err := transport.ParseRequest(r, &reminderUpdateRequest); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedReminder, err := h.repo.UpdateReminder(ctx, reminderUpdateRequest.ToDomain())
	if err != nil {
		var noResourceErr *repository.NoResourceFoundError
		if errors.As(err, &noResourceErr) {
			writeJSONError(w, http.StatusNotFound, err.Error())
			return
		}
		var invalidRRuleErr *repository.ErrInvalidRRule
		if errors.As(err, &invalidRRuleErr) {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "Unexpected error occurred")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedReminder)
}

func (h *ReminderHandler) handleDeleteReminder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reminderDeleteRequest transport.ReminderDeleteRequest
	if err := transport.ParseRequest(r, &reminderDeleteRequest); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.repo.DeleteReminder(ctx, reminderDeleteRequest.ToDomain())
	if err != nil {
		var noResourceErr *repository.NoResourceFoundError
		if errors.As(err, &noResourceErr) {
			writeJSONError(w, http.StatusNotFound, err.Error())
			return
		} else {
			writeJSONError(w, http.StatusInternalServerError, "Unexpected error occurred")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
