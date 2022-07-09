package auth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("supersecretkey")

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	UserId   uint   `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWT(email string, username string, id uint) (tokenString string, err error) {
	expirationTime := time.Now().Add(3 * time.Hour)
	claims := &JWTClaim{
		Email:    email,
		Username: username,
		UserId:   id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return tokenString, err
}
