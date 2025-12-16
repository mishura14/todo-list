package checkhash

import (
	hashbcrypt "git-register-project/internal/servise/bcryptHash/hash"
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	password := "mishura14"
	hash, err := hashbcrypt.HashBcrypt(password)
	if err != nil {
		t.Error("hashing password failed", err)
	}
	cases := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "valid password",
			password: password,
			hash:     hash,
			want:     true,
		},
		{
			name:     "invalid password",
			password: "invalid",
			hash:     hash,
			want:     false,
		},
		{
			name:     "empty password",
			password: "",
			hash:     hash,
			want:     false,
		},
		{
			name:     "empty hash",
			password: password,
			hash:     "",
			want:     false,
		},
		{
			name:     "empty password and hash",
			password: "",
			hash:     "",
			want:     false,
		},
		{
			name:     "invalid hash",
			password: password,
			hash:     "invalid",
			want:     false,
		},
		{
			name:     "invalid password and hash",
			password: "invalid",
			hash:     "invalid",
			want:     false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := CheckHash(c.password, c.hash)
			if got != c.want {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}
