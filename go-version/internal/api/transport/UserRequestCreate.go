package transport

import (
	"encoding/json"
	"go-version/internal/api/domain"
	"net/http"
)

type UserCreateRequest struct {
	NoContext
	NoQueryParams
	NoURLParams
	// Request Body
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (r *UserCreateRequest) ParseFromBody(req *http.Request) error {
	return json.NewDecoder(req.Body).Decode(r)
}

func (r *UserCreateRequest) Validate() error {
	var errors []error
	if r.Name == nil || *r.Name == "" {
		errors = append(errors, &ErrNameRequired{})
	}
	if r.Email == nil || *r.Email == "" {
		errors = append(errors, &ErrEmailRequired{})
	}
	if r.Password == nil || *r.Password == "" {
		errors = append(errors, &ErrPasswordRequired{})
	}
	if len(errors) > 0 {
		return &ErrBadRequest{Errs: errors}
	}
	return nil
}

func (r *UserCreateRequest) ToDomain() *domain.UserCreateDomain {
	return &domain.UserCreateDomain{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}
}
