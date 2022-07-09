package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"strings"
	"time"
)

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	UserId   uint   `json:"user_id"`
	jwt.StandardClaims
}

func ValidateToken(signedToken string) (err error, claim JWTClaim) {
	var jwtKey = []byte(os.Getenv("TOKEN_SECRET"))
	signedToken = strings.Split(signedToken, " ")[1]
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return err, *claims
}
