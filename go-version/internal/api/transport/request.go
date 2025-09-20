package transport

import (
	"go-version/internal/contextkeys"
	"net/http"
	"net/url"
)

type RequestParser[T any] interface {
	*T
	Validate() error
	ParseFromBody(r *http.Request) error
	ParseFromQuery(values url.Values) error
	ParseFromURLParams(r *http.Request) error
	ParseFromContext(r *http.Request) error
}

type BaseRequest struct{}

func ParseRequest[T any, P RequestParser[T]](r *http.Request, req P) error {
	if err := req.ParseFromURLParams(r); err != nil {
		return err
	}

	if len(r.URL.RawQuery) > 0 {
		if err := req.ParseFromQuery(r.URL.Query()); err != nil {
			return err
		}
	}

	if r.Method != "GET" && r.Method != "DELETE" && r.Header.Get("Content-Type") == "application/json" {
		req.ParseFromBody(r)
	}

	if err := req.ParseFromContext(r); err != nil {
		return err
	}

	return req.Validate()
}

func getUserIdFromContext(r *http.Request) (string, error) {
	ctx := r.Context()
	userId, ok := contextkeys.UserIdFromContext(ctx)
	if !ok || userId == "" {
		return "", &contextkeys.ErrUserNotInContext{}
	}
	return userId, nil
}
