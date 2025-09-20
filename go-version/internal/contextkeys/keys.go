package contextkeys

import "context"

type ContextKey string

const UserIDKey ContextKey = "userId"

// UserIdFromContext retrieves the authenticated user id (sub claim).
func UserIdFromContext(ctx context.Context) (string, bool) {
	uid, ok := ctx.Value(UserIDKey).(string)
	return uid, ok
}
