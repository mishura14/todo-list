package servise

import "golang.org/x/crypto/bcrypt"

// проверка хешированого пароля с не хешированым
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
