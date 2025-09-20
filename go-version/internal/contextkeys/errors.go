package contextkeys

type ErrUserNotInContext struct{}

func (e *ErrUserNotInContext) Error() string {
	return "user id not found in context"
}
