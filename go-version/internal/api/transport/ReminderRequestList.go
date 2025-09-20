package transport

import (
	"go-version/internal/api/domain"
	"go-version/internal/api/utils"
	"net/url"
	"time"
)

type ReminderListRequest struct {
	UserIDContext
	NoRequestBody
	NoURLParams

	// Query Params
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
	Search    *string `json:"search"`
}

func (r *ReminderListRequest) ParseFromQuery(values url.Values) error {
	startDate := values.Get("start_date")
	endDate := values.Get("end_date")
	search := values.Get("search")

	if startDate != "" {
		r.StartDate = &startDate
	}
	if endDate != "" {
		r.EndDate = &endDate
	}
	if search != "" {
		r.Search = &search
	}
	return nil
}

func (r *ReminderListRequest) Validate() error {
	var errors []error
	if r.StartDate != nil && r.EndDate == nil {
		errors = append(errors, &ErrEndDateRequired{})
	}
	if r.EndDate != nil && r.StartDate == nil {
		errors = append(errors, &ErrStartDateRequired{})
	}
	if (r.StartDate != nil && !utils.IsValidDateTime(*r.StartDate)) || (r.EndDate != nil && !utils.IsValidDateTime(*r.EndDate)) {
		errors = append(errors, &ErrInvalidDateFormat{})
	}
	if len(errors) > 0 {
		return &ErrBadRequest{Errs: errors}
	}
	return nil
}

func (r *ReminderListRequest) ToDomain() *domain.ReminderListDomain {
	var startDate, endDate *time.Time
	if r.StartDate != nil {
		sd, _ := utils.ParseDateTime(*r.StartDate)
		startDate = &sd
	}
	if r.EndDate != nil {
		ed, _ := utils.ParseDateTime(*r.EndDate)
		endDate = &ed
	}
	return &domain.ReminderListDomain{
		StartDate: startDate,
		EndDate:   endDate,
		Search:    r.Search,
	}
}
