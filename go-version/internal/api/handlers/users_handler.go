package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"go-version/internal/api/middleware"
	"go-version/internal/api/repository"
	"go-version/internal/api/transport"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) (*UserHandler, error) {
	return &UserHandler{repo: repo}, nil
}

func (h *UserHandler) RegisterRoutes(router chi.Router) {
	authMw, err := middleware.AuthMiddleware(context.Background())
	if err != nil {
		panic(err)
	}

	h.registerPublicRoutes(router)
	h.registerProtectedRoutes(router, authMw)
}

func (h *UserHandler) registerPublicRoutes(router chi.Router) {
	router.Post("/users", h.handleCreateUser)
}

func (h *UserHandler) registerProtectedRoutes(router chi.Router, authMw func(http.Handler) http.Handler) {
	router.With(authMw).Get("/users", h.handleGetUser)
}

func (h *UserHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req transport.UserGetRequest
	if err := transport.ParseRequest(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.repo.GetUser(ctx, req.ToDomain())

	response := transport.NewGetUserResult(user)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to fetch user")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func (h *UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req transport.UserCreateRequest
	if err := transport.ParseRequest(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	createdUser, err := h.repo.CreateUser(ctx, req.ToDomain())

	response := transport.NewCreateUserResult(createdUser)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
