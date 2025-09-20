package transport

import (
	"encoding/json"
	"go-version/internal/api/domain"
	"go-version/internal/api/utils"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type ReminderUpdateRequest struct {
	UserIDContext
	NoQueryParams

	// URL Params
	ReminderID string `json:"-" db:"-"`

	// Request Body
	RRule       *string `json:"rrule" db:"rrule"`
	Description *string `json:"description" db:"description"`
	StartAt     *string `json:"start_at" db:"start_at"`
}

func (r *ReminderUpdateRequest) ParseFromBody(req *http.Request) error {
	return json.NewDecoder(req.Body).Decode(r)
}

func (r *ReminderUpdateRequest) ParseFromURLParams(req *http.Request) error {
	r.ReminderID = chi.URLParam(req, "reminderId")
	return nil
}

func (r *ReminderUpdateRequest) Validate() error {
	var errors []error

	// if rrule supplied, must be non-empty and valid
	if r.RRule != nil && *r.RRule == "" {
		errors = append(errors, &ErrRRuleEmpty{})
	} else if r.RRule != nil && !utils.IsValidRRule(*r.RRule) {
		errors = append(errors, &ErrInvalidRRuleFormat{})
	}

	// if start_at supplied, must be valid datetime
	if r.StartAt != nil && !utils.IsValidDateTime(*r.StartAt) {
		errors = append(errors, &ErrInvalidStartAt{})
	}

	// at least one field must be supplied
	if r.RRule == nil && r.Description == nil && r.StartAt == nil {
		errors = append(errors, &ErrNoFieldsToUpdate{})
	}

	if len(errors) > 0 {
		return &ErrBadRequest{Errs: errors}
	}
	return nil

}

func (r *ReminderUpdateRequest) ToDomain() *domain.ReminderUpdateDomain {
	var startAt *time.Time
	if r.StartAt != nil {
		sa, _ := utils.ParseDateTime(*r.StartAt)
		startAt = &sa
	}
	return &domain.ReminderUpdateDomain{
		ReminderID:  r.ReminderID,
		RRule:       r.RRule,
		Description: r.Description,
		StartAt:     startAt,
	}
}
