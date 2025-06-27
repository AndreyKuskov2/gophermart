package validator

import "testing"

func TestLuhnAlgorith(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		name     string
	}{
		{"79927398713", true, "valid Luhn number"},
		{"1234567812345670", true, "valid Luhn number 2"},
		{"79927398710", false, "invalid Luhn number"},
		{"abcdefg", false, "non-digit input"},
		{"", false, "empty string"},
		{"0", true, "single zero"},
		{"059", true, "valid short Luhn"},
		{"059a", false, "valid digits with letter"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LuhnAlgorith(tt.input)
			if result != tt.expected {
				t.Errorf("LuhnAlgorith(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
