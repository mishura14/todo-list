package password_hash

import (
	check_hash_password "git-register-project/internal/servise/hash_password/check_hach_password"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "mishura_14_12_2010(google)"

	// 1. Хешируем пароль
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// 2. Проверяем что хеш НЕ пустой
	if hash == "" {
		t.Errorf("HashPassword returned an empty string")
	}

	// 3. Проверяем что хеш НЕ равен паролю
	if hash == password {
		t.Errorf("Hash should NOT equal original password")
	}

	// 4. Проверяем что можно верифицировать хеш
	if !check_hash_password.CheckPasswordHash(password, hash) {
		t.Errorf("Cannot verify generated hash")
	}

	// 5. Хешируем тот же пароль ещё раз
	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	// 6. Хеши должны быть РАЗНЫМИ (из-за соли в bcrypt)
	if hash1 == hash2 {
		t.Errorf("Same password should produce DIFFERENT hashes due to salt")
	}

	// 7. Проверяем уникальность хешей при многократном хешировании
	m := map[string]bool{}
	for i := 0; i < 10; i++ {
		h, err := HashPassword(password)
		if err != nil {
			t.Fatalf("HashPassword failed: %v", err)
		}
		m[h] = true
	}

	// 8. Все 10 хешей должны быть РАЗНЫМИ
	if len(m) != 10 {
		t.Errorf("Expected 10 unique hashes, got %d", len(m))
	}
}
