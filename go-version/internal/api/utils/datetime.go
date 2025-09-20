package utils

import (
	"fmt"
	"time"
)

func ParseDateTime(dateTimeStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05", // RFC no timezone
		"2006-01-02 15:04:05", // Space separated
		"2006-01-02",          // YYYY-MM-DD
		"01/02/2006",          // MM/DD/YYYY
		"2006/01/02",          // YYYY/MM/DD
		"2006-01-02 15:04",    // YYYY-MM-DD HH:MM
	}

	for _, format := range formats {
		if parsedTime, err := time.Parse(format, dateTimeStr); err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse datetime string: %s", dateTimeStr)
}
