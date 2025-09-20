package repository

import (
	"context"

	"go-version/internal/api/auth"
	"go-version/internal/api/domain"
	"go-version/internal/api/models"
	"go-version/internal/api/store"

	"github.com/google/uuid"
)

type UserRepositoryInterface interface {
	GetUser(ctx context.Context, userId string) (*UserGetResult, error)
	CreateUser(ctx context.Context, user *domain.UserCreateDomain) (*UserCreateResult, error)
}

type UserRepository struct {
	store store.UserStoreInterface
}

func NewUserRepository(store store.UserStoreInterface) (*UserRepository, error) {
	return &UserRepository{store: store}, nil
}

func (r *UserRepository) GetUser(ctx context.Context, req *domain.UserGetDomain) (*UserGetResult, error) {
	user, err := r.store.GetUser(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return NewUserGetResult(user), nil
}

func (r *UserRepository) CreateUser(ctx context.Context, req *domain.UserCreateDomain) (*UserCreateResult, error) {
	newUser := &models.User{
		Id:        uuid.New().String(),
		Name:      *req.Name,
		Email:     *req.Email,
		Password:  *req.Password,
		ApiKey:    nil,
		CreatedAt: nil,
		UpdatedAt: nil,
	}

	apiKey, err := auth.GenerateToken(newUser.Id)
	if err != nil {
		return nil, err
	}

	newUser.ApiKey = &apiKey

	createdUser, err := r.store.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return NewUserCreateResult(createdUser), nil
}
