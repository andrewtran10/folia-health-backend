package store

import (
	"context"
	"database/sql"
	"fmt"

	"go-version/internal/api/models"
)

type ReminderStoreInterface interface {
	ListReminders(ctx context.Context, filters *ReminderListFilters) ([]models.Reminder, error)
	GetReminderByID(ctx context.Context, userID, reminderID string) (*models.Reminder, error)
	CreateReminder(ctx context.Context, reminder *models.Reminder) (*models.Reminder, error)
	UpdateReminder(ctx context.Context, reminder *models.Reminder) (*models.Reminder, error)
	DeleteReminder(ctx context.Context, userID, reminderID string) error
}

type ReminderStore struct {
	db *sql.DB
}

func NewReminderStore(db *sql.DB) (*ReminderStore, error) {
	return &ReminderStore{db: db}, nil
}

func (s *ReminderStore) ListReminders(ctx context.Context, filters *ReminderListFilters) ([]models.Reminder, error) {
	query := `
		SELECT id, user_id, rrule, description, start_at, created_at, updated_at
		FROM reminders
		WHERE user_id=$1
	`

	args := []interface{}{filters.UserID}
	argIdx := 2

	if filters.Search != nil {
		query += fmt.Sprintf(` AND LOWER(description) LIKE LOWER($%d)`, argIdx)
		args = append(args, "%"+*filters.Search+"%")
		argIdx++
	}

	if filters.StartDate != nil && filters.EndDate != nil {
		query += fmt.Sprintf(` AND start_at <= $%d`, argIdx)
		args = append(args, *filters.EndDate)
		argIdx++
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []models.Reminder
	for rows.Next() {
		var reminder models.Reminder
		if err := rows.Scan(&reminder.Id, &reminder.UserId, &reminder.RRule, &reminder.Description, &reminder.StartAt, &reminder.CreatedAt, &reminder.UpdatedAt); err != nil {
			return nil, err
		}
		reminders = append(reminders, reminder)
	}

	return reminders, nil
}

func (s *ReminderStore) GetReminderByID(ctx context.Context, userID, reminderID string) (*models.Reminder, error) {
	query := `
		SELECT id, user_id, rrule, description, start_at, created_at, updated_at
		FROM reminders
		WHERE id=$1 AND user_id=$2
	`

	var reminder models.Reminder
	err := s.db.QueryRowContext(ctx, query, reminderID, userID).Scan(&reminder.Id, &reminder.UserId, &reminder.RRule, &reminder.Description, &reminder.StartAt, &reminder.CreatedAt, &reminder.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &NoReminderFoundError{
				ID: reminderID,
			}
		}
		return nil, err
	}
	return &reminder, nil
}

func (s *ReminderStore) CreateReminder(ctx context.Context, reminder *models.Reminder) (*models.Reminder, error) {
	query := `
		INSERT INTO reminders (id, user_id, rrule, description, start_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, rrule, description, start_at, created_at, updated_at
	`

	var newReminder models.Reminder
	err := s.db.QueryRowContext(ctx, query,
		reminder.Id,
		reminder.UserId,
		reminder.RRule,
		reminder.Description,
		reminder.StartAt,
	).Scan(&newReminder.Id, &newReminder.UserId, &newReminder.RRule, &newReminder.Description, &newReminder.StartAt, &newReminder.CreatedAt, &newReminder.UpdatedAt)

	return &newReminder, err
}

func (s *ReminderStore) UpdateReminder(ctx context.Context, reminder *models.Reminder) (*models.Reminder, error) {
	query := `
        UPDATE reminders 
        SET rrule = $1, description = $2, start_at = $3, updated_at = NOW()
        WHERE id = $4 AND user_id = $5
        RETURNING id, user_id, rrule, description, start_at, created_at, updated_at
    `

	var updatedReminder models.Reminder
	err := s.db.QueryRowContext(ctx, query,
		reminder.RRule,
		reminder.Description,
		reminder.StartAt,
		reminder.Id,
		reminder.UserId,
	).Scan(&updatedReminder.Id, &updatedReminder.UserId, &updatedReminder.RRule,
		&updatedReminder.Description, &updatedReminder.StartAt,
		&updatedReminder.CreatedAt, &updatedReminder.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &NoReminderFoundError{ID: reminder.Id}
		}
		return nil, err
	}

	return &updatedReminder, nil
}

func (s *ReminderStore) DeleteReminder(ctx context.Context, userID, reminderID string) error {
	query := `DELETE FROM reminders WHERE id=$1 AND user_id=$2`
	result, err := s.db.ExecContext(ctx, query, reminderID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &NoReminderFoundError{
			ID: reminderID,
		}
	}

	return nil
}
