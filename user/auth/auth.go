package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go_practice/user/logger"
	"go_practice/user/models"
	"os"
	"strconv"
	"strings"
	"time"
)

type NonAuthJWTClaim struct {
	Email string `json:"email"`
	jwt.StandardClaims
}
type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	UserId   uint   `json:"user_id"`
	IsAdmin  bool   `json:"is_admin"`
	IsActive bool   `json:"is_active"`
	jwt.StandardClaims
}

func GenerateJWT(email string, username string, id uint, isAdmin bool, isActive bool) (tokenString string, refreshTokenString string, err error) {
	var jwtKey = []byte(os.Getenv("TOKEN_SECRET"))
	//var anotherJwtKey = []byte(os.Getenv("ANOTHER_TOKEN_SECRET"))
	var refreshTokenKey = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))
	tokenTimeout, err := strconv.Atoi(os.Getenv("TOKEN_TIMEOUT"))
	duration := time.Duration(tokenTimeout) * time.Minute
	expirationTime := time.Now().Add(duration)
	claims := &JWTClaim{
		Email:    email,
		Username: username,
		UserId:   id,
		IsActive: isActive,
		IsAdmin:  isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	refreshTokenTimeout, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TIMEOUT"))
	refreshDuration := time.Duration(refreshTokenTimeout) * time.Hour
	refreshExpirationTime := time.Now().Add(refreshDuration)
	refreshClaims := &JWTClaim{
		Email:    email,
		Username: username,
		UserId:   id,
		IsActive: isActive,
		IsAdmin:  isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	tokenString, err = token.SignedString(jwtKey)
	refreshTokenString, err = refreshToken.SignedString(refreshTokenKey)
	return tokenString, refreshTokenString, err
}

func GenerateNonAuthJWT(email string) (tokenString string, err error) {
	var jwtKey = []byte(os.Getenv("ANOTHER_TOKEN_SECRET"))
	tokenTimeout, err := strconv.Atoi(os.Getenv("ANOTHER_TOKEN_TIMEOUT"))
	duration := time.Duration(tokenTimeout) * time.Hour
	expirationTime := time.Now().Add(duration)
	claims := &JWTClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err = token.SignedString(jwtKey)

	return tokenString, err
}

func RefreshValidateToken(signedToken string, jwtKey []byte) (err error, claim JWTClaim) {
	var token_ models.Token
	if len(jwtKey) == 0 {
		err = errors.New("refresh token can not be validated")
		logger.Error.Println("ERROR: secret key is not provided")
		return
	}
	if signedToken == "" {
		err = errors.New("no refresh token is provided")
		logger.Info.Println(err.Error())
		return
	}

	signedToken = strings.Split(signedToken, " ")[1]
	err = token_.GetTokenByRefreshToken(signedToken)
	if err != nil {
		err = errors.New("refresh token can not be validated")
		logger.Error.Println(err.Error())
		return
	}
	if !token_.IsActive {
		var tokenActive models.Token
		err = errors.New("old refresh token")
		err_ := tokenActive.UpdateToken(token_.Email, map[string]interface{}{"is_active": false})
		if err_ != nil {
			logger.Error.Println("the tokens can not be deactivated, potential risks")
			err_ = errors.New("refresh token can not be validated")
			return
		}
		logger.Error.Println(err.Error())
		return
	}
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
		err = errors.New("refresh token can not be validated")
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
		err = errors.New("refresh token expired")
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

func ValidateNonAuthToken(signedToken string, jwtKey []byte) (err error, claim NonAuthJWTClaim) {
	if signedToken == "" {
		err = errors.New("no token is provided")
		logger.Info.Println(err.Error())
		return
	}
	token, err := jwt.ParseWithClaims(
		signedToken,
		&NonAuthJWTClaim{},
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
	claims, ok := token.Claims.(*NonAuthJWTClaim)
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
	return err, *claims
}
