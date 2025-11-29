package servise

import "unicode"

func ContainsDigits(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func ContainsSpecialChars(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func CheckPassword(s string) bool {
	return len(s) >= 8 && ContainsDigits(s) && ContainsSpecialChars(s)
}
