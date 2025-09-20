package transport

import "go-version/internal/api/repository"

type UserGetResponse struct {
	Id    *string `json:"id"`
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

func NewGetUserResult(user *repository.UserGetResult) *UserGetResponse {
	return &UserGetResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}
}
