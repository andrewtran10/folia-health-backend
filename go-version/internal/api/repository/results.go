package repository

import (
	"go-version/internal/api/models"
	"time"
)

type UserCreateResult struct {
	Id     *string `json:"id"`
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	ApiKey *string `json:"apiKey"`
}

type UserGetResult struct {
	Id    *string `json:"id"`
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type ReminderListResult struct {
	Reminders []models.Reminder `json:"reminders"`
}

type ReminderCreateResult struct {
	Id          *string    `json:"id"`
	RRule       *string    `json:"rrule"`
	Description *string    `json:"description"`
	StartAt     *time.Time `json:"startAt"`
}

type ReminderUpdateResult struct {
	Id          *string    `json:"id"`
	RRule       *string    `json:"rrule"`
	Description *string    `json:"description"`
	StartAt     *time.Time `json:"startAt"`
}

// Model -> Result converters
func NewUserCreateResult(user *models.User) *UserCreateResult {
	return &UserCreateResult{
		Id:     &user.Id,
		Name:   &user.Name,
		Email:  &user.Email,
		ApiKey: user.ApiKey,
	}
}

func NewUserGetResult(user *models.User) *UserGetResult {
	return &UserGetResult{
		Id:    &user.Id,
		Name:  &user.Name,
		Email: &user.Email,
	}
}

func NewReminderCreateResult(reminder *models.Reminder) *ReminderCreateResult {
	return &ReminderCreateResult{
		Id:          &reminder.Id,
		RRule:       &reminder.RRule,
		Description: reminder.Description,
		StartAt:     &reminder.StartAt,
	}
}

func NewReminderUpdateResult(reminder *models.Reminder) *ReminderUpdateResult {
	return &ReminderUpdateResult{
		Id:          &reminder.Id,
		RRule:       &reminder.RRule,
		Description: reminder.Description,
		StartAt:     &reminder.StartAt,
	}
}

func NewReminderListResult(reminders []models.Reminder) *ReminderListResult {
	if reminders == nil {
		reminders = []models.Reminder{}
	}
	return &ReminderListResult{
		Reminders: reminders,
	}
}
