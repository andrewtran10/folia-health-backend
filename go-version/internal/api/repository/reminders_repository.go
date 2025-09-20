package repository

import (
	"context"
	"errors"
	"time"

	"go-version/internal/api/domain"
	"go-version/internal/api/models"
	"go-version/internal/api/store"

	"github.com/google/uuid"
	"github.com/teambition/rrule-go"
)

type ReminderRepositoryInterface interface {
	ListReminders(ctx context.Context, params *domain.ReminderListDomain) (*ReminderListResult, error)
	CreateReminder(ctx context.Context, params *domain.ReminderCreateDomain) (*ReminderCreateResult, error)
	UpdateReminder(ctx context.Context, params *domain.ReminderUpdateDomain) (*ReminderUpdateResult, error)
	DeleteReminder(ctx context.Context, params *domain.ReminderDeleteDomain) error
}

type ReminderRepository struct {
	reminderStore store.ReminderStoreInterface
}

func NewReminderRepository(reminderStore store.ReminderStoreInterface) (*ReminderRepository, error) {
	return &ReminderRepository{reminderStore: reminderStore}, nil
}

func (r *ReminderRepository) ListReminders(ctx context.Context, reminderListRequest *domain.ReminderListDomain) (*ReminderListResult, error) {
	filters := &store.ReminderListFilters{
		UserID:    reminderListRequest.UserID,
		Search:    reminderListRequest.Search,
		StartDate: reminderListRequest.StartDate,
		EndDate:   reminderListRequest.EndDate,
	}

	reminders, err := r.reminderStore.ListReminders(ctx, filters)
	if err != nil {
		return nil, err
	}

	for i := range reminders {
		reminders[i].PopulateMetadataFields(reminderListRequest.StartDate, reminderListRequest.EndDate)
	}

	return NewReminderListResult(reminders), nil
}

func (r *ReminderRepository) CreateReminder(ctx context.Context, req *domain.ReminderCreateDomain) (*ReminderCreateResult, error) {
	if !validateReminderOccurrences(req.RRule, req.StartAt) {
		return nil, &ErrInvalidRRule{
			Err: errors.New("no occurrences can be generated with the provided rrule and start_at"),
		}
	}

	newReminder := &models.Reminder{
		Id:          uuid.New().String(),
		UserId:      req.UserID,
		RRule:       req.RRule,
		Description: req.Description,
		StartAt:     req.StartAt,
		CreatedAt:   nil,
		UpdatedAt:   nil,
	}

	createdReminder, err := r.reminderStore.CreateReminder(ctx, newReminder)
	if err != nil {
		return nil, err
	}

	return NewReminderCreateResult(createdReminder), nil
}

func (r *ReminderRepository) UpdateReminder(ctx context.Context, req *domain.ReminderUpdateDomain) (*ReminderUpdateResult, error) {
	curReminder, err := r.reminderStore.GetReminderByID(ctx, req.UserID, req.ReminderID)
	if err != nil {
		return nil, &NoResourceFoundError{Err: err}
	}

	updates := &models.Reminder{
		Id:          curReminder.Id,
		UserId:      curReminder.UserId,
		RRule:       curReminder.RRule,
		Description: curReminder.Description,
		StartAt:     curReminder.StartAt,
		CreatedAt:   curReminder.CreatedAt,
		UpdatedAt:   nil,
	}

	if req.RRule != nil || req.StartAt != nil || req.Description != nil {
		if req.RRule != nil {
			updates.RRule = *req.RRule
		}
		if req.StartAt != nil {
			updates.StartAt = *req.StartAt
		}
		if req.Description != nil {
			updates.Description = req.Description
		}
		if !validateReminderOccurrences(updates.RRule, updates.StartAt) {
			return nil, &ErrInvalidRRule{
				Err: errors.New("no occurrences can be generated with the provided rrule and start_at"),
			}
		}
	}

	updatedReminder, err := r.reminderStore.UpdateReminder(ctx, updates)
	if err != nil {
		return nil, &NoResourceFoundError{Err: err}
	}

	return NewReminderUpdateResult(updatedReminder), nil
}

func (r *ReminderRepository) DeleteReminder(ctx context.Context, req *domain.ReminderDeleteDomain) error {
	err := r.reminderStore.DeleteReminder(ctx, req.UserID, req.ReminderID)
	if err != nil {
		return &NoResourceFoundError{Err: err}
	}
	return nil
}

func validateReminderOccurrences(rruleStr string, startAt time.Time) bool {
	rruleObj, _ := rrule.StrToRRule(rruleStr)

	rruleObj.DTStart(startAt)

	occurrences := rruleObj.All()

	return len(occurrences) > 0
}
