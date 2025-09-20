package repository

import (
	"context"
	"errors"
	"testing"

	"go-version/internal/api/domain"
	"go-version/internal/api/models"
	"go-version/internal/api/store/mocks"
	"go-version/internal/api/utils"

	"go.uber.org/mock/gomock"
)

func TestNewUserRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockUserStoreInterface(ctrl)

	repo, err := NewUserRepository(mockStore)
	if err != nil {
		t.Errorf("NewUserRepository() returned unexpected error: %v", err)
	}
	if repo == nil {
		t.Error("NewUserRepository() returned nil repository")
	}
	if repo != nil && repo.store == nil {
		t.Error("UserRepository store should not be nil")
	}
}

func TestUserRepository_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockUserStoreInterface(ctrl)

	testCases := []struct {
		name          string
		userID        string
		setupMock     func()
		expectedError bool
		expectedName  string
	}{
		{
			name:   "successful user retrieval",
			userID: "user-123",
			setupMock: func() {
				mockStore.EXPECT().
					GetUser(gomock.Any(), "user-123").
					Return(&models.User{
						Id:    "user-123",
						Name:  "John Doe",
						Email: "john@example.com",
					}, nil).
					Times(1)
			},
			expectedError: false,
			expectedName:  "John Doe",
		},
		{
			name:   "user not found",
			userID: "nonexistent-user",
			setupMock: func() {
				mockStore.EXPECT().
					GetUser(gomock.Any(), "nonexistent-user").
					Return(nil, errors.New("user not found")).
					Times(1)
			},
			expectedError: true,
		},
		{
			name:   "database error",
			userID: "user-123",
			setupMock: func() {
				mockStore.EXPECT().
					GetUser(gomock.Any(), "user-123").
					Return(nil, errors.New("database connection error")).
					Times(1)
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup the mock expectations
			tc.setupMock()

			// Create repository with mock store
			repo := &UserRepository{store: mockStore}

			req := &domain.UserGetDomain{
				UserID: tc.userID,
			}

			result, err := repo.GetUser(context.Background(), req)

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
				} else {
					if *result.Name != tc.expectedName {
						t.Errorf("Expected name %s, got %s", tc.expectedName, *result.Name)
					}
					if *result.Id != tc.userID {
						t.Errorf("Expected ID %s, got %s", tc.userID, *result.Id)
					}
				}
			}
		})
	}
}

func TestUserRepository_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockUserStoreInterface(ctrl)

	testCases := []struct {
		name          string
		inputDomain   *domain.UserCreateDomain
		setupMock     func()
		expectedError bool
	}{
		{
			name: "successful user creation",
			inputDomain: &domain.UserCreateDomain{
				Name:     utils.StringPtr("John Doe"),
				Email:    utils.StringPtr("john@example.com"),
				Password: utils.StringPtr("hashedpassword123"),
			},
			setupMock: func() {
				mockStore.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, user *models.User) (*models.User, error) {
						// Validate the input
						if user.Name != "John Doe" {
							t.Errorf("Expected name 'John Doe', got %s", user.Name)
						}
						if user.Email != "john@example.com" {
							t.Errorf("Expected email 'john@example.com', got %s", user.Email)
						}
						if user.Id == "" {
							t.Error("Expected UUID to be generated")
						}
						if user.ApiKey == nil {
							t.Error("Expected ApiKey to be set")
						}
						// Return the same user (simulating database insert)
						return user, nil
					}).
					Times(1)
			},
			expectedError: false,
		},
		{
			name: "database error during creation",
			inputDomain: &domain.UserCreateDomain{
				Name:     utils.StringPtr("Jane Doe"),
				Email:    utils.StringPtr("jane@example.com"),
				Password: utils.StringPtr("hashedpassword456"),
			},
			setupMock: func() {
				mockStore.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database constraint violation")).
					Times(1)
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			repo := &UserRepository{store: mockStore}

			result, err := repo.CreateUser(context.Background(), tc.inputDomain)

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
				}
			}
		})
	}
}
