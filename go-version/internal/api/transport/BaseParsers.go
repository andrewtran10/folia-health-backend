package transport

import (
	"net/http"
	"net/url"
)

type UserIDContext struct {
	UserID string `json:"-" db:"-"`
}

func (r *UserIDContext) ParseFromContext(req *http.Request) error {
	userID, err := getUserIdFromContext(req)
	if err != nil {
		return err
	}
	r.UserID = userID
	return nil
}

type NoContext struct{}

func (r *NoContext) ParseFromContext(req *http.Request) error {
	return nil
}

type NoRequestBody struct{}

func (r *NoRequestBody) ParseFromBody(req *http.Request) error {
	return nil
}

type NoQueryParams struct{}

func (r *NoQueryParams) ParseFromQuery(values url.Values) error {
	return nil
}

type NoURLParams struct{}

func (r *NoURLParams) ParseFromURLParams(req *http.Request) error {
	return nil
}
