package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateTokens(email string, secret string, tokenttl time.Duration) (string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["exp"] = time.Now().Add(tokenttl).Unix()
	accessTokenString, err := accessToken.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}
