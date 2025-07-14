package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(UserID uint, secret string) (string, error){
	claims := &Claims{
		UserID: UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 *time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken (tokenString, secret string)(*Claims, error){
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token)(interface{},error){
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}