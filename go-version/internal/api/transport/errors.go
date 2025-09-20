package transport

import "strings"

type ErrBadRequest struct {
	Errs []error
}

func (e *ErrBadRequest) Error() string {
	if len(e.Errs) == 0 {
		return "Bad request"
	}
	msgs := make([]string, len(e.Errs))
	for i, err := range e.Errs {
		msgs[i] = err.Error()
	}
	return "Bad request: " + strings.Join(msgs, "; ")
}

type ErrNameRequired struct{}

func (e *ErrNameRequired) Error() string {
	return "name is required"
}

type ErrEmailRequired struct{}

func (e *ErrEmailRequired) Error() string {
	return "email is required"
}

type ErrPasswordRequired struct{}

func (e *ErrPasswordRequired) Error() string {
	return "password is required"
}

type ErrRRuleRequired struct{}

func (e *ErrRRuleRequired) Error() string {
	return "rrule is required"
}

type ErrStartAtRequired struct{}

func (e *ErrStartAtRequired) Error() string {
	return "start_at is required"
}

type ErrRRuleEmpty struct{}

func (e *ErrRRuleEmpty) Error() string {
	return "rrule cannot be empty"
}

type ErrInvalidRRuleFormat struct{}

func (e *ErrInvalidRRuleFormat) Error() string {
	return "rrule is not in a valid format"
}

type ErrStartAtEmpty struct{}

func (e *ErrStartAtEmpty) Error() string {
	return "start_at cannot be empty"
}

type ErrInvalidStartAt struct{}

func (e *ErrInvalidStartAt) Error() string {
	return "start_at is invalid"
}

type ErrNoFieldsToUpdate struct{}

func (e *ErrNoFieldsToUpdate) Error() string {
	return "request body is empty"
}

type ErrStartDateRequired struct{}

func (e *ErrStartDateRequired) Error() string {
	return "start_date is required"
}

type ErrEndDateRequired struct{}

func (e *ErrEndDateRequired) Error() string {
	return "end_date is required"
}

type ErrInvalidDateFormat struct{}

func (e *ErrInvalidDateFormat) Error() string {
	return "date is not in a valid datetime format"
}
