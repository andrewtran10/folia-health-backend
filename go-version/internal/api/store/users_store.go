package store

import (
	"context"
	"database/sql"
	"fmt"

	"go-version/internal/api/models"
)

type UserStoreInterface interface {
	GetUser(ctx context.Context, userId string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(conn *sql.DB) (*UserStore, error) {
	return &UserStore{db: conn}, nil
}

func (s *UserStore) GetUser(ctx context.Context, userId string) (*models.User, error) {
	query := `SELECT id, email, name, password, api_key, created_at, updated_at
		FROM users
		WHERE id=$1`

	var user models.User
	err := s.db.QueryRowContext(ctx, query, userId).Scan(&user.Id, &user.Email, &user.Name, &user.Password, &user.ApiKey, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (s *UserStore) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	var createdUser models.User
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO users (id, email, name, password, api_key)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, email, name, password, api_key`,
		user.Id, user.Email, user.Name, user.Password, user.ApiKey,
	).Scan(&createdUser.Id, &createdUser.Email, &createdUser.Name, &createdUser.Password, &createdUser.ApiKey)
	if err != nil {
		fmt.Println("Error inserting user:", err)
		return nil, err
	}

	return &createdUser, nil
}
