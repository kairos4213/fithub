package auth

import (
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
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

func TestJWTValidation(t *testing.T) {
	userID := uuid.New()
	privKey, _ := os.ReadFile("../../private_key.pem")
	pubKey, _ := os.ReadFile("../../public_key.pem")
	validToken, _ := MakeJWT(userID, false, privKey)

	tests := map[string]struct {
		tokenString string
		publicKey   []byte
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		"valid token": {
			tokenString: validToken,
			publicKey:   pubKey,
			wantUserID:  userID,
			wantErr:     false,
		},
		"invalid token": {
			tokenString: "invalid.token.string",
			publicKey:   pubKey,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		"incorrect public key": {
			tokenString: validToken,
			publicKey:   []byte("invalid.public.key"),
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			claims, err := ValidateJWT(tc.tokenString, tc.publicKey)
			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error: %v, got: %v", tc.wantErr, err)
				return
			}
			if claims.UserID != tc.wantUserID {
				t.Errorf("expected userID: %v, got: %v", tc.wantUserID, claims.UserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := map[string]struct {
		headers   http.Header
		wantToken string
		wantErr   bool
	}{
		"valid authorization header": {
			headers: http.Header{
				"Authorization": []string{"Bearer valid_token"},
			},
			wantToken: "valid_token",
			wantErr:   false,
		},
		"missing authorization header": {
			headers:   http.Header{},
			wantToken: "",
			wantErr:   true,
		},
		"malformed authorization header": {
			headers: http.Header{
				"Authorization": []string{"InvalidBearer token"},
			},
			wantToken: "",
			wantErr:   true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bearerToken, err := GetBearerToken(tc.headers)
			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error: %v, got: %v", tc.wantErr, err)
				return
			}
			if bearerToken != tc.wantToken {
				t.Errorf("expected userID: %v, got: %v", tc.wantToken, bearerToken)
			}
		})
	}
}
