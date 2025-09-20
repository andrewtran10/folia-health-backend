package utils

import "github.com/teambition/rrule-go"

func IsValidRRule(rruleStr string) bool {
	rrule, err := rrule.StrToRRule(rruleStr)

	// Basic validation checks needed because rrule-go is lenient
	if rrule != nil && rrule.Options.Count < 1 {
		return false
	}
	if rrule != nil && rrule.Options.Interval < 1 {
		return false
	}

	return err == nil
}

func IsValidDateTime(dateTimeStr string) bool {
	_, err := ParseDateTime(dateTimeStr)
	return err == nil
}
