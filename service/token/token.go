package token

import (
	"github.com/dgrijalva/jwt-go"
	"os"
)

func ValidateToken(tokenString string) (uint64, uint64) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	claims := token.Claims.(jwt.MapClaims)
	return uint64(claims["sub"].(float64)), uint64(claims["role"].(float64))
}
