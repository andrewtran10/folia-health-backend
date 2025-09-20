package store

import "time"

type ReminderListFilters struct {
	UserID    string
	Search    *string
	StartDate *time.Time
	EndDate   *time.Time
	Limit     *int
	Offset    *int
}
