package utils

import "github.com/teambition/rrule-go"

func IsValidRRule(rruleStr string) bool {
	_, err := rrule.StrToRRule(rruleStr)
	return err == nil
}

func IsValidDateTime(dateTimeStr string) bool {
	_, err := ParseDateTime(dateTimeStr)
	return err == nil
}
