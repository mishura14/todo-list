package hashbcrypt

import "golang.org/x/crypto/bcrypt"

// хеширование пароля
func HashBcrypt(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
