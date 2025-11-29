package servise

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// генератор кода подтверждения
func GenerateSecureCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	code := n.Int64() + 100000
	return fmt.Sprintf("%06d", code)
}
