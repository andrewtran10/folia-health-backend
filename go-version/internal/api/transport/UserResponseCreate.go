package transport

import "go-version/internal/api/repository"

type UserCreateResponse struct {
	Id     *string `json:"id"`
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	ApiKey *string `json:"api_key"`
}

func NewCreateUserResult(user *repository.UserCreateResult) *UserCreateResponse {
	return &UserCreateResponse{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		ApiKey: user.ApiKey,
	}
}
