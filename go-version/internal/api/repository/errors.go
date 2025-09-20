package repository

type NoResourceFoundError struct {
	Err error
}

func (e *NoResourceFoundError) Error() string {
	return e.Err.Error()
}

type ErrInvalidRRule struct {
	Err error
}

func (e *ErrInvalidRRule) Error() string {
	return "rrule is invalid: " + e.Err.Error()
}
