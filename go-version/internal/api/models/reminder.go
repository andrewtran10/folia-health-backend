package models

import (
	"time"

	"github.com/teambition/rrule-go"
)

type Reminder struct {
	Id          string      `db:"id" json:"id"`
	UserId      string      `db:"user_id" json:"user_id"`
	RRule       string      `db:"rrule" json:"rrule"`
	Description *string     `db:"description" json:"description"`
	StartAt     time.Time   `db:"start_at" json:"start_at"`
	CreatedAt   *time.Time  `db:"created_at" json:"-"`
	UpdatedAt   *time.Time  `db:"updated_at" json:"-"`
	Occurrences []time.Time `db:"-" json:"occurrences,omitempty"`
}

func (r *Reminder) PopulateMetadataFields(start, end *time.Time) {
	if start != nil && end != nil {
		occurrences, err := r.generateOccurrences(*start, *end)
		if err == nil {
			r.Occurrences = occurrences
		}
	}
}

func (r *Reminder) generateOccurrences(startDate, endDate time.Time) ([]time.Time, error) {
	rruleObj, err := rrule.StrToRRule(r.RRule)
	if err != nil {
		return nil, err
	}

	rruleObj.DTStart(r.StartAt)

	return rruleObj.Between(startDate, endDate, true), nil
}
