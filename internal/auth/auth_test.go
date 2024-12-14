package auth

import (
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	correctPassword := "passwordIsCorrect"
	incorrectPassword := "notCorrect"
	correctHash, _ := HashPassword(correctPassword)
	incorrectHash, _ := HashPassword(incorrectPassword)

	tests := map[string]struct {
		password string
		hash     string
		wantErr  bool
	}{
		"correct": {
			password: correctPassword,
			hash:     correctHash,
			wantErr:  false,
		},
		"incorrect password": {
			password: incorrectPassword,
			hash:     correctHash,
			wantErr:  true,
		},
		"incorrect hash": {
			password: correctPassword,
			hash:     incorrectHash,
			wantErr:  true,
		},
		"blank password": {
			password: "",
			hash:     correctHash,
			wantErr:  true,
		},
		"invalid hash": {
			password: correctPassword,
			hash:     "notValidHash",
			wantErr:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := CheckPasswordHash(tc.password, tc.hash)
			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}
