package checkhash

import "golang.org/x/crypto/bcrypt"

// проверка хешированого пароля с не хешированым
func CheckHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
