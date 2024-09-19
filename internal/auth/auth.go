package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateJwt(userId int, tokenSecret string, expiresInSeconds int) (string, error) {
	signingKey := []byte(tokenSecret)
	issuedAt := time.Now().UTC()
	registeredClaims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(ValidateExpiresInSeconds(issuedAt, expiresInSeconds)),
		Subject:   fmt.Sprintf("%d", userId),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims)
	ss, err := token.SignedString(signingKey)
	if err != nil {
		fmt.Printf("Error: createJwt - %s", err)
		return "", err
	}
	return ss, nil
}

func ValidateExpiresInSeconds(issuedAt time.Time, expiresInSeconds int) time.Time {

	if expiresInSeconds < 1 || expiresInSeconds > 86400 {
		return issuedAt.Add(24 * time.Hour)
	}

	return issuedAt.Add(time.Duration(expiresInSeconds) * time.Second)
}

func ValidateJwt(tokenString, tokenSecret string) (string, error) {
	claims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return "", err
	}

	userIdString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}

	if issuer != string("chirpy") {
		return "", errors.New("invalid issuer")
	}

	return userIdString, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no auth header included in request")
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}
	return splitAuth[1], nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}
