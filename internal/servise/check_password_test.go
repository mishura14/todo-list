package servise

import (
	"testing"
)

func TestValidPasswordCheck(t *testing.T) {
	cases := []struct {
		name     string
		password string
		expected bool
	}{
		{"bas_not_number", "password@", false},
		{"bas_yes_number", "p@ssword123", true},
		{"bad_not_simvol", "password123", false},
		{"bad_yes_simvol", "p@ssword123", true},
		{"bad_len<8", "pas123", false},
		{"bad_len>8", "password_12345678", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := CheckPassword(c.password)
			if result != c.expected {
				t.Errorf("Expected %v, got %v", c.expected, result)
			}
		})
	}
}
