package utils

import (
	"testing"
)

func TestIsValidRRule(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectValid bool
	}{
		// Valid RRule strings
		{
			name:        "valid daily",
			input:       "FREQ=DAILY;COUNT=10",
			expectValid: true,
		},
		{
			name:        "valid weekly",
			input:       "FREQ=WEEKLY;BYDAY=MO,WE,FR;COUNT=5",
			expectValid: true,
		},
		{
			name:        "valid monthly",
			input:       "FREQ=MONTHLY;BYMONTHDAY=15;COUNT=3",
			expectValid: true,
		},
		{
			name:        "valid yearly",
			input:       "FREQ=YEARLY;BYMONTH=12;BYMONTHDAY=25;COUNT=2",
			expectValid: true,
		},

		// Invalid RRule strings
		{
			name:        "invalid frequency",
			input:       "FREQ=INVALID;COUNT=10",
			expectValid: false,
		},
		{
			name:        "trailing semicolon",
			input:       "FREQ=DAILY;COUNT=2;",
			expectValid: false,
		},
		{
			name:        "missing semicolon",
			input:       "FREQ=WEEKLY;BYDAY=MO,WE,FR COUNT=5;",
			expectValid: false,
		},
		{
			name:        "negative count",
			input:       "FREQ=DAILY;COUNT=-5",
			expectValid: false,
		},
		{
			name:        "zero count",
			input:       "FREQ=DAILY;COUNT=0",
			expectValid: false,
		},
		{
			name:        "invalid byday value",
			input:       "FREQ=WEEKLY;BYDAY=XX",
			expectValid: false,
		},
		{
			name:        "invalid bymonthday value",
			input:       "FREQ=MONTHLY;BYMONTHDAY=32;COUNT=3",
			expectValid: false,
		},
		{
			name:        "invalid bymonth value",
			input:       "FREQ=YEARLY;BYMONTH=13;BYMONTHDAY=25;COUNT=2",
			expectValid: false,
		},
		{
			name:        "invalid byhour value",
			input:       "FREQ=DAILY;BYHOUR=24",
			expectValid: false,
		},
		{
			name:        "invalid byminute value",
			input:       "FREQ=DAILY;BYMINUTE=60",
			expectValid: false,
		},
		{
			name:        "invalid bysecond value",
			input:       "FREQ=DAILY;BYSECOND=60",
			expectValid: false,
		},
		{
			name:        "invalid interval value",
			input:       "FREQ=DAILY;INTERVAL=0",
			expectValid: false,
		},

		{
			name:        "invalid parameter",
			input:       "FREQ=DAILY;COUNT=10;INVALID=PARAM",
			expectValid: false,
		},

		{"empty string", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			valid := IsValidRRule(tc.input)
			if valid != tc.expectValid {
				t.Errorf("IsValidRRule(%q) = %v; want %v", tc.input, valid, tc.expectValid)
			}
		})
	}
}
