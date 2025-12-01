package servise

import (
	"strconv"
	"testing"
)

func TestGenerateSecureCode_Comprehensive(t *testing.T) {
	// Генерируем много кодов для статистики
	codes := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		codes[i] = GenerateSecureCode()
	}

	tests := []struct {
		name     string
		testFunc func() bool
	}{
		{
			name: "length_is_6",
			testFunc: func() bool {
				for _, code := range codes {
					if len(code) != 6 {
						return false
					}
				}
				return true
			},
		},
		{
			name: "only_digits",
			testFunc: func() bool {
				for _, code := range codes {
					for _, ch := range code {
						if ch < '0' || ch > '9' {
							return false
						}
					}
				}
				return true
			},
		},
		{
			name: "in_range_100000_999999",
			testFunc: func() bool {
				for _, code := range codes {
					n, err := strconv.Atoi(code)
					if err != nil || n < 100000 || n > 999999 {
						return false
					}
				}
				return true
			},
		},
		{
			name: "reasonable_uniqueness",
			testFunc: func() bool {
				unique := make(map[string]bool)
				for _, code := range codes {
					unique[code] = true
				}
				// Ожидаем хотя бы 900 уникальных из 1000
				return len(unique) > 900
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.testFunc() {
				t.Errorf("%s failed", tt.name)
			}
		})
	}
}
