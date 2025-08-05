package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const minCost = 11

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), minCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

type CustomClaims struct {
	UserID  uuid.UUID
	IsAdmin bool
	jwt.RegisteredClaims
}

func MakeJWT(userID uuid.UUID, isAdmin bool, privateKeyBytes []byte) (string, error) {
	claims := &CustomClaims{
		userID,
		isAdmin,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "fithub",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return "", err
	}

	ss, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return ss, nil
}

func ValidateJWT(tokenString string, publicKey []byte) (*CustomClaims, error) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		return pubKey, nil
	})
	if err != nil {
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
	refreshTokenBase := make([]byte, 10)
	_, err := rand.Read(refreshTokenBase)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(refreshTokenBase), nil
}
