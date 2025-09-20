package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-version/internal/api/domain"
	"go-version/internal/api/models"
	"go-version/internal/api/store"
	"go-version/internal/api/store/mocks"
	"go-version/internal/api/utils"

	"go.uber.org/mock/gomock"
)

func TestNewReminderRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockReminderStoreInterface(ctrl)

	repo, err := NewReminderRepository(mockStore)
	if err != nil {
		t.Errorf("NewReminderRepository() returned unexpected error: %v", err)
	}
	if repo == nil {
		t.Error("NewReminderRepository() returned nil repository")
	}
	if repo != nil && repo.reminderStore == nil {
		t.Error("ReminderRepository store should not be nil")
	}
}

func TestReminderRepository_ListReminders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockReminderStoreInterface(ctrl)

	startDate := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 10, 31, 23, 59, 59, 0, time.UTC)
	search := "test"

	testCases := []struct {
		name          string
		request       *domain.ReminderListDomain
		setupMock     func()
		expectedError bool
		expectedCount int
	}{
		{
			name: "successful list with filters",
			request: &domain.ReminderListDomain{
				UserID:    "user-123",
				StartDate: &startDate,
				EndDate:   &endDate,
				Search:    &search,
			},
			setupMock: func() {
				expectedFilters := &store.ReminderListFilters{
					UserID:    "user-123",
					StartDate: &startDate,
					EndDate:   &endDate,
					Search:    &search,
				}
				mockStore.EXPECT().
					ListReminders(gomock.Any(), expectedFilters).
					Return([]models.Reminder{
						{
							Id:          "reminder-1",
							UserId:      "user-123",
							RRule:       "FREQ=DAILY;COUNT=5",
							Description: utils.StringPtr("Daily reminder"),
							StartAt:     startDate,
						},
						{
							Id:          "reminder-2",
							UserId:      "user-123",
							RRule:       "FREQ=WEEKLY;COUNT=3",
							Description: utils.StringPtr("Weekly reminder"),
							StartAt:     startDate,
						},
					}, nil).
					Times(1)
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name: "successful list without filters",
			request: &domain.ReminderListDomain{
				UserID: "user-123",
			},
			setupMock: func() {
				expectedFilters := &store.ReminderListFilters{
					UserID:    "user-123",
					StartDate: nil,
					EndDate:   nil,
					Search:    nil,
				}
				mockStore.EXPECT().
					ListReminders(gomock.Any(), expectedFilters).
					Return([]models.Reminder{{
						Id:          "reminder-1",
						UserId:      "user-123",
						RRule:       "FREQ=DAILY;COUNT=5",
						Description: utils.StringPtr("Daily reminder"),
						StartAt:     startDate,
					}}, nil).
					Times(1)
			},
			expectedError: false,
			expectedCount: 1,
		},
		{
			name: "store error",
			request: &domain.ReminderListDomain{
				UserID: "user-123",
			},
			setupMock: func() {
				mockStore.EXPECT().
					ListReminders(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			repo := &ReminderRepository{reminderStore: mockStore}

			result, err := repo.ListReminders(context.Background(), tc.request)

			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if result != nil {
					t.Errorf("Expected nil result but got: %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("Expected result but got nil")
				} else if len(result.Reminders) != tc.expectedCount {
					t.Errorf("Expected %d reminders, got %d", tc.expectedCount, len(result.Reminders))
				}
			}
		})
	}
}

func TestReminderRepository_CreateReminder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockReminderStoreInterface(ctrl)

	startTime := time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC)

	testCases := []struct {
		name          string
		request       *domain.ReminderCreateDomain
		setupMock     func()
		expectedError bool
		validateError func(error) bool
	}{
		{
			name: "successful creation",
			request: &domain.ReminderCreateDomain{
				UserID:      "user-123",
				RRule:       "FREQ=DAILY;COUNT=5",
				Description: utils.StringPtr("Daily reminder"),
				StartAt:     startTime,
			},
			setupMock: func() {
				mockStore.EXPECT().
					CreateReminder(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, reminder *models.Reminder) (*models.Reminder, error) {
						if reminder.UserId != "user-123" {
							t.Errorf("Expected UserId 'user-123', got %s", reminder.UserId)
						}
						if reminder.RRule != "FREQ=DAILY;COUNT=5" {
							t.Errorf("Expected RRule 'FREQ=DAILY;COUNT=5', got %s", reminder.RRule)
						}
						if reminder.Id == "" {
							t.Error("Expected generated UUID for ID")
						}
						if !reminder.StartAt.Equal(startTime) {
							t.Errorf("Expected StartAt %v, got %v", startTime, reminder.StartAt)
						}
						return reminder, nil
					}).
					Times(1)
			},
			expectedError: false,
		},
		{
			name: "invalid rrule with no occurrences",
			request: &domain.ReminderCreateDomain{
				UserID:      "user-123",
				RRule:       "FREQ=DAILY;UNTIL=20200101T000000Z",
				Description: utils.StringPtr("Invalid reminder"),
				StartAt:     startTime,
			},
			setupMock:     func() {},
			expectedError: true,
			validateError: func(err error) bool {
				_, ok := err.(*ErrInvalidRRule)
				return ok
			},
		},
		{
			name: "store creation error",
			request: &domain.ReminderCreateDomain{
				UserID:      "user-123",
				RRule:       "FREQ=DAILY;COUNT=5",
				Description: utils.StringPtr("Daily reminder"),
				StartAt:     startTime,
			},
			setupMock: func() {
				mockStore.EXPECT().
					CreateReminder(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database constraint violation")).
					Times(1)
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			repo := &ReminderRepository{reminderStore: mockStore}

			result, err := repo.CreateReminder(context.Background(), tc.request)

			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if tc.validateError != nil && !tc.validateError(err) {
					t.Errorf("Error validation failed for error: %v", err)
				}
				if result != nil {
					t.Errorf("Expected nil result but got: %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("Expected result but got nil")
				}
			}
		})
	}
}

func TestReminderRepository_UpdateReminder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockReminderStoreInterface(ctrl)

	originalTime := time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2023, 10, 2, 11, 0, 0, 0, time.UTC)

	existingReminder := &models.Reminder{
		Id:          "reminder-123",
		UserId:      "user-123",
		RRule:       "FREQ=DAILY;COUNT=5",
		Description: utils.StringPtr("Original reminder"),
		StartAt:     originalTime,
		CreatedAt:   &originalTime,
	}

	testCases := []struct {
		name          string
		request       *domain.ReminderUpdateDomain
		setupMock     func()
		expectedError bool
		validateError func(error) bool
	}{
		{
			name: "successful update with all fields",
			request: &domain.ReminderUpdateDomain{
				UserID:      "user-123",
				ReminderID:  "reminder-123",
				RRule:       utils.StringPtr("FREQ=WEEKLY;COUNT=3"),
				Description: utils.StringPtr("Updated reminder"),
				StartAt:     &updatedTime,
			},
			setupMock: func() {
				mockStore.EXPECT().
					GetReminderByID(gomock.Any(), "user-123", "reminder-123").
					Return(existingReminder, nil).
					Times(1)

				mockStore.EXPECT().
					UpdateReminder(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, reminder *models.Reminder) (*models.Reminder, error) {
						if reminder.RRule != "FREQ=WEEKLY;COUNT=3" {
							t.Errorf("Expected updated RRule 'FREQ=WEEKLY;COUNT=3', got %s", reminder.RRule)
						}
						if *reminder.Description != "Updated reminder" {
							t.Errorf("Expected updated description 'Updated reminder', got %s", *reminder.Description)
						}
						if !reminder.StartAt.Equal(updatedTime) {
							t.Errorf("Expected updated StartAt %v, got %v", updatedTime, reminder.StartAt)
						}
						return reminder, nil
					}).
					Times(1)
			},
			expectedError: false,
		},
		{
			name: "reminder not found",
			request: &domain.ReminderUpdateDomain{
				UserID:     "user-123",
				ReminderID: "nonexistent",
			},
			setupMock: func() {
				mockStore.EXPECT().
					GetReminderByID(gomock.Any(), "user-123", "nonexistent").
					Return(nil, errors.New("reminder not found")).
					Times(1)
			},
			expectedError: true,
			validateError: func(err error) bool {
				_, ok := err.(*NoResourceFoundError)
				return ok
			},
		},
		{
			name: "invalid rrule update",
			request: &domain.ReminderUpdateDomain{
				UserID:     "user-123",
				ReminderID: "reminder-123",
				RRule:      utils.StringPtr("FREQ=DAILY;UNTIL=20200101T000000Z"),
			},
			setupMock: func() {
				mockStore.EXPECT().
					GetReminderByID(gomock.Any(), "user-123", "reminder-123").
					Return(existingReminder, nil).
					Times(1)
			},
			expectedError: true,
			validateError: func(err error) bool {
				_, ok := err.(*ErrInvalidRRule)
				return ok
			},
		},
		{
			name: "store update error",
			request: &domain.ReminderUpdateDomain{
				UserID:      "user-123",
				ReminderID:  "reminder-123",
				Description: utils.StringPtr("Updated description"),
			},
			setupMock: func() {
				mockStore.EXPECT().
					GetReminderByID(gomock.Any(), "user-123", "reminder-123").
					Return(existingReminder, nil).
					Times(1)

				mockStore.EXPECT().
					UpdateReminder(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database update error")).
					Times(1)
			},
			expectedError: true,
			validateError: func(err error) bool {
				_, ok := err.(*NoResourceFoundError)
				return ok
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			repo := &ReminderRepository{reminderStore: mockStore}

			result, err := repo.UpdateReminder(context.Background(), tc.request)

			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if tc.validateError != nil && !tc.validateError(err) {
					t.Errorf("Error validation failed for error: %v", err)
				}
				if result != nil {
					t.Errorf("Expected nil result but got: %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("Expected result but got nil")
				}
			}
		})
	}
}

func TestReminderRepository_DeleteReminder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockReminderStoreInterface(ctrl)

	testCases := []struct {
		name          string
		request       *domain.ReminderDeleteDomain
		setupMock     func()
		expectedError bool
		validateError func(error) bool
	}{
		{
			name: "successful deletion",
			request: &domain.ReminderDeleteDomain{
				UserID:     "user-123",
				ReminderID: "reminder-123",
			},
			setupMock: func() {
				mockStore.EXPECT().
					DeleteReminder(gomock.Any(), "user-123", "reminder-123").
					Return(nil).
					Times(1)
			},
			expectedError: false,
		},
		{
			name: "reminder not found",
			request: &domain.ReminderDeleteDomain{
				UserID:     "user-123",
				ReminderID: "nonexistent",
			},
			setupMock: func() {
				mockStore.EXPECT().
					DeleteReminder(gomock.Any(), "user-123", "nonexistent").
					Return(errors.New("reminder not found")).
					Times(1)
			},
			expectedError: true,
			validateError: func(err error) bool {
				_, ok := err.(*NoResourceFoundError)
				return ok
			},
		},
		{
			name: "store deletion error",
			request: &domain.ReminderDeleteDomain{
				UserID:     "user-123",
				ReminderID: "reminder-123",
			},
			setupMock: func() {
				mockStore.EXPECT().
					DeleteReminder(gomock.Any(), "user-123", "reminder-123").
					Return(errors.New("database error")).
					Times(1)
			},
			expectedError: true,
			validateError: func(err error) bool {
				_, ok := err.(*NoResourceFoundError)
				return ok
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			repo := &ReminderRepository{reminderStore: mockStore}

			err := repo.DeleteReminder(context.Background(), tc.request)

			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if tc.validateError != nil && !tc.validateError(err) {
					t.Errorf("Error validation failed for error: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateReminderOccurrences(t *testing.T) {
	startTime := time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC)

	testCases := []struct {
		name     string
		rrule    string
		startAt  time.Time
		expected bool
	}{
		{
			name:     "valid rrule with occurrences",
			rrule:    "FREQ=DAILY;COUNT=5",
			startAt:  startTime,
			expected: true,
		},
		{
			name:     "rrule with no occurrences",
			rrule:    "FREQ=DAILY;UNTIL=20200101T000000Z",
			startAt:  startTime,
			expected: false,
		},
		{
			name:     "weekly rrule with occurrences",
			rrule:    "FREQ=WEEKLY;BYDAY=MO,WE,FR;COUNT=3",
			startAt:  startTime,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validateReminderOccurrences(tc.rrule, tc.startAt)
			if result != tc.expected {
				t.Errorf("validateReminderOccurrences(%s, %v) = %v; want %v", tc.rrule, tc.startAt, result, tc.expected)
			}
		})
	}
}
