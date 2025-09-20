package api

import (
	"database/sql"

	"go-version/internal/api/handlers"
	"go-version/internal/api/repository"
	"go-version/internal/api/store"

	"github.com/go-chi/chi/v5"
)

type ApiService struct {
	handlers map[string]handlers.HttpHandler
}

func NewService(db *sql.DB) *ApiService {

	handlersMap := make(map[string]handlers.HttpHandler)

	userStore, _ := store.NewUserStore(db)
	userRepository, _ := repository.NewUserRepository(userStore)
	userHandler, _ := handlers.NewUserHandler(userRepository)
	handlersMap["users"] = userHandler

	reminderStore, _ := store.NewReminderStore(db)
	reminderRepository, _ := repository.NewReminderRepository(reminderStore)
	reminderHandler, _ := handlers.NewReminderHandler(reminderRepository)
	handlersMap["reminders"] = reminderHandler

	return &ApiService{
		handlers: handlersMap,
	}

}

func (s *ApiService) GetHandlers() map[string]handlers.HttpHandler {
	return s.handlers
}

func (s *ApiService) RegisterRoutes(r chi.Router) {
	apiRouter := chi.NewRouter()

	for _, handler := range s.handlers {
		handler.RegisterRoutes(apiRouter)
	}

	r.Mount("/api", apiRouter)
}
