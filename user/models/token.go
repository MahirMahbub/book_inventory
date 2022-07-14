package models

import "github.com/jinzhu/gorm"

type Token struct {
	gorm.Model
	AccessToken  string `gorm:"not null"`
	RefreshToken string `gorm:"not null"`
	ChildID      uint   `gorm:"default:null"`
	IsActive     bool   `gorm:"default:true"`
	Email        string `gorm:"not null"`
}
type Tokens []Token

func (token *Token) GetTokenByAccessToken(accessToken string) (err error) {
	return DB.Where("access_token = ?", accessToken).First(&token).Error
}

func (token *Token) GetTokenByRefreshToken(refreshToken string) (err error) {
	return DB.Where("refresh_token = ?", refreshToken).First(&token).Error
}

func (token *Token) CreateToken() (err error) {
	return DB.Create(&token).Error
}

func (token *Token) UpdateToken(email string, data map[string]interface{}) (err error) {
	return DB.Model(&token).Where("email = ?", email).Updates(data).Error
}

func (token *Token) UpdateTokenByID(ID uint, data map[string]interface{}) (err error) {
	return DB.Model(&token).Where("email = ?", ID).Updates(data).Error
}
