package transport

import (
	"net/http"

	"go-version/internal/api/domain"

	"github.com/go-chi/chi/v5"
)

type ReminderDeleteRequest struct {
	UserIDContext
	NoRequestBody
	NoQueryParams

	// URL Params
	ReminderID string `json:"-" db:"-"`
}

func (r *ReminderDeleteRequest) ParseFromURLParams(req *http.Request) error {
	r.ReminderID = chi.URLParam(req, "reminderId")
	return nil
}

func (r *ReminderDeleteRequest) Validate() error {
	return nil
}

func (r *ReminderDeleteRequest) ToDomain() *domain.ReminderDeleteDomain {
	return &domain.ReminderDeleteDomain{
		ReminderID: r.ReminderID,
	}
}
