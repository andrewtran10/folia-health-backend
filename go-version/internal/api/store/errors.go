package store

type NoReminderFoundError struct {
	Resource string
	ID       string
}

func (e *NoReminderFoundError) Error() string {
	return "no reminder found with ID " + e.ID
}
