package auth

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTKey() []byte {
	key := os.Getenv("JWT_SUPER_SECRET_SIGNING_KEY")
	return []byte(key)
}

func GenerateToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(getJWTKey())
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return getJWTKey(), nil
	})
}
