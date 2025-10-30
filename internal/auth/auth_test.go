package auth

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	correctPassword := "passwordIsCorrect"
	incorrectPassword := "notCorrect"
	correctHash, _ := HashPassword(correctPassword)
	incorrectHash, _ := HashPassword(incorrectPassword)

	tests := map[string]struct {
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		"correct": {
			password:      correctPassword,
			hash:          correctHash,
			wantErr:       false,
			matchPassword: true,
		},
		"incorrect password": {
			password:      incorrectPassword,
			hash:          correctHash,
			wantErr:       false,
			matchPassword: false,
		},
		"incorrect hash": {
			password:      correctPassword,
			hash:          incorrectHash,
			wantErr:       false,
			matchPassword: false,
		},
		"blank password": {
			password:      "",
			hash:          correctHash,
			wantErr:       false,
			matchPassword: false,
		},
		"invalid hash": {
			password:      correctPassword,
			hash:          "notValidHash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			match, err := CheckPasswordHash(tc.password, tc.hash)
			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error: %v, got: %v", tc.wantErr, err)
			}
			if !tc.wantErr && match != tc.matchPassword {
				t.Fatalf("expected password match to be: %v, got: %v", tc.matchPassword, match)
			}
		})
	}
}

func TestJWTValidation(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, false, "secret")

	tests := map[string]struct {
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		"valid token": {
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		"invalid token": {
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		"incorrect secret": {
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			claims, err := ValidateJWT(tc.tokenString, tc.tokenSecret)
			if (err != nil) != tc.wantErr {
				t.Errorf("expected error: %v, got: %v", tc.wantErr, err)
				return
			}
			if claims != nil && claims.UserID != tc.wantUserID {
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
