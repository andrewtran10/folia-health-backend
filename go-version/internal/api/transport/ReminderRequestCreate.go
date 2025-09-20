package transport

import (
	"encoding/json"
	"go-version/internal/api/domain"
	"go-version/internal/api/utils"
	"net/http"
)

type ReminderCreateRequest struct {
	UserIDContext
	NoURLParams
	NoQueryParams

	// Request Body
	RRule       *string `json:"rrule"`
	Description *string `json:"description"`
	StartAt     *string `json:"start_at"`
}

func (r *ReminderCreateRequest) ParseFromBody(req *http.Request) error {
	return json.NewDecoder(req.Body).Decode(r)
}

func (r *ReminderCreateRequest) Validate() error {
	var errors []error
	if r.RRule == nil || *r.RRule == "" {
		errors = append(errors, &ErrRRuleRequired{})
	} else if !utils.IsValidRRule(*r.RRule) {
		errors = append(errors, &ErrInvalidRRuleFormat{})
	}

	if r.StartAt == nil || *r.StartAt == "" {
		errors = append(errors, &ErrStartAtRequired{})
	} else if !utils.IsValidDateTime(*r.StartAt) {
		errors = append(errors, &ErrInvalidStartAt{})
	}

	if len(errors) > 0 {
		return &ErrBadRequest{Errs: errors}
	}
	return nil
}

func (r *ReminderCreateRequest) ToDomain() *domain.ReminderCreateDomain {
	startAt, _ := utils.ParseDateTime(*r.StartAt)
	return &domain.ReminderCreateDomain{
		RRule:       *r.RRule,
		Description: r.Description,
		StartAt:     startAt,
	}
}
