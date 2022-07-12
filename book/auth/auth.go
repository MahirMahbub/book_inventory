package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go_practice/book/logger"
	"os"
	"strings"
	"time"
)

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	UserId   uint   `json:"user_id"`
	IsAdmin  bool   `json:"is_admin"`
	IsActive bool   `json:"is_active"`
	jwt.StandardClaims
}

func ValidateToken(signedToken string) (err error, claim JWTClaim) {
	var jwtKey = []byte(os.Getenv("TOKEN_SECRET"))
	if len(jwtKey) == 0 || jwtKey == nil {
		err = errors.New("token can not be validated")
		logger.Error.Println("secret key is not provided")
		return
	}
	if signedToken == "" {
		err = errors.New("no token is provided")
		logger.Info.Println(err.Error())
		return
	}
	signedToken = strings.Split(signedToken, " ")[1]
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logger.Info.Println("signature is invalid")
			return
		}
		err = errors.New("token can not be validated")
		logger.Info.Println(err.Error())
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		logger.Info.Println(err.Error())
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		logger.Info.Println(err.Error())
		return
	}
	if !claims.IsActive {
		err = errors.New("user is not active")
		logger.Info.Println(err.Error())
		return
	}
	return err, *claims
}
