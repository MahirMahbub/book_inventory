package models

import "github.com/jinzhu/gorm"

type Token struct {
	gorm.Model
	AccessToken  string
	RefreshToken string
	ChildID      uint `gorm:"default:null"`
	IsActive     bool `gorm:"default:true"`
	Email        string
}
type Tokens []Token

func (token *Token) GetTokenByAccessToken(AccessToken string) (err error) {
	return DB.Where("access_token = ?", AccessToken).First(&token).Error
}

func (token *Token) GetTokenByRefreshToken(RefreshToken string) (err error) {
	return DB.Where("refresh_token = ?", RefreshToken).First(&token).Error
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
