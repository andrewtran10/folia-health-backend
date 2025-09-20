package transport

import (
	"go-version/internal/api/domain"
)

type UserGetRequest struct {
	UserIDContext
	NoRequestBody
	NoQueryParams
	NoURLParams
}

func (r *UserGetRequest) Validate() error {
	return nil
}

func (r *UserGetRequest) ToDomain() *domain.UserGetDomain {
	return &domain.UserGetDomain{
		UserID: r.UserID,
	}
}
