package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	jwtExpiry         = 15 * time.Minute
	refreshTokenBytes = 32
)

type CustomClaims struct {
	UserID  uuid.UUID
	IsAdmin bool
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	// TODO: Look into optimizing argon2id Params
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func MakeJWT(userID uuid.UUID, isAdmin bool, tokenSecret string) (string, error) {
	claims := &CustomClaims{
		userID,
		isAdmin,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "fithub",
		},
	}
	signingKey := []byte(tokenSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("error parsing custom claims")
	}

	if claims.Issuer != "fithub" {
		return nil, errors.New("invalid token issuer")
	}

	return claims, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authheader := headers.Get("authorization")
	if authheader == "" {
		return "", errors.New("no authorization header found")
	}
	splitAuth := strings.Split(authheader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}
	return splitAuth[1], nil
}

func MakeRefreshToken() (string, error) {
	refreshTokenBase := make([]byte, refreshTokenBytes)
	_, err := rand.Read(refreshTokenBase)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(refreshTokenBase), nil
}
