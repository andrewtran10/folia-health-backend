package utils

import "testing"

func Test_ParseDateTime(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectError bool
	}{
		// Valid date-time strings
		{
			name:        "valid rfc3339",
			input:       "2023-01-02T15:04:05+07:00",
			expectError: false,
		},
		{
			name:        "valid rfc3339nano",
			input:       "2023-01-02T15:04:05.999999999+07:00",
			expectError: false,
		},
		{
			name:        "valid rfc3339 utc",
			input:       "2023-01-02T15:04:05Z",
			expectError: false,
		},
		{
			name:        "valid rfc3339 no timezone",
			input:       "2023-01-02T15:04:05",
			expectError: false,
		},
		{
			name:        "valid space separated",
			input:       "2023-01-02 15:04:05",
			expectError: false,
		},
		{
			name:        "valid yyyy-mm-dd",
			input:       "2023-10-05",
			expectError: false,
		},
		{
			name:        "valid mm/dd/yyyy",
			input:       "10/05/2023",
			expectError: false,
		},
		{
			name:        "valid yyyy/mm/dd",
			input:       "2023/10/05",
			expectError: false,
		},
		{
			name:        "valid yyyy-mm-dd hh:mm",
			input:       "2023-10-05 14:30",
			expectError: false,
		},

		// Invalid date-time strings
		{
			name:        "invalid month",
			input:       "2023-13-01",
			expectError: true,
		},
		{
			name:        "invalid day",
			input:       "2023-12-32",
			expectError: true,
		},
		{
			name:        "invalid day for February",
			input:       "2023-02-30",
			expectError: true,
		},
		{
			name:        "invalid month with slashes",
			input:       "2023/13/01",
			expectError: true,
		},
		{
			name:        "invalid hour",
			input:       "2023-01-01T25:00:00Z",
			expectError: true,
		},
		{
			name:        "not-a-date",
			input:       "not-a-date",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseDateTime(tc.input)
			if (err != nil) != tc.expectError {
				t.Errorf("ParseDateTime(%q) error = %v; want error: %v", tc.input, err != nil, tc.expectError)
			}
		})
	}

}
