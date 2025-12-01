package servise

import (
	"strings"
	"testing"
)

func TestCheckEmail(t *testing.T) {
	email := "test@example.com"
	result := CheckEmail(email)
	if !result {
		t.Errorf("valid email %s", email)
	}
}
func TestValidCheckEmail(t *testing.T) {
	cases := []struct {
		name string
		in   string
		exp  bool
	}{
		{
			name: "bad_email_len<3",
			in:   "a@",
			exp:  false,
		},
		{
			name: "bad_email_len>254",
			in:   strings.Repeat("a", 255) + "@example.com",
			exp:  false,
		},
		{
			name: "bad_email_no_at",
			in:   "testexample.com",
			exp:  false,
		},
		{
			name: "bad_email_no_domain",
			in:   "test@",
			exp:  false,
		},
		{
			name: "bad_email_no_local_part",
			in:   "@example.com",
			exp:  false,
		},
		{
			name: "bad_email_no_domain",
			in:   "andr@gmail.com",
			exp:  true,
		},
		{
			name: "bad_email",
			in:   "andr12@gmail.com",
			exp:  true,
		},
	}
	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			result := CheckEmail(tCase.in)
			if result != tCase.exp {
				t.Errorf("expected %v, got %v", tCase.exp, result)
			}
		})
	}
}
