package domain

import "time"

type ReminderCreateDomain struct {
	UserID      string
	RRule       string
	Description *string
	StartAt     time.Time
}

type ReminderUpdateDomain struct {
	UserID      string
	ReminderID  string
	RRule       *string
	Description *string
	StartAt     *time.Time
}

type ReminderListDomain struct {
	UserID    string
	StartDate *time.Time
	EndDate   *time.Time
	Search    *string
}

type ReminderDeleteDomain struct {
	UserID     string
	ReminderID string
}
