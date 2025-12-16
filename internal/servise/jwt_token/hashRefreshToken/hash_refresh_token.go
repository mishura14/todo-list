package hashrefreshtoken

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token)) // 32 байта
	return hex.EncodeToString(hash[:])
}
