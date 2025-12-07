package valid_password

import "unicode"

// проверка пароля на наличие цифр
func ContainsDigits(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// проверка пароля на наличие спецсимволов
func ContainsSpecialChars(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// проверка пароля на длину
func CheckPassword(s string) bool {
	return len(s) >= 8 && ContainsDigits(s) && ContainsSpecialChars(s)
}
