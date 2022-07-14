package models

import "github.com/jinzhu/gorm"

type NonAuthToken struct {
	gorm.Model
	PasswordChangeToken string `gorm:"default:null"`
	UserVerifyToken     string `gorm:"not null"`
	Email               string `gorm:"not null"`
}

func (token *NonAuthToken) UpdateNonAuthToken(email string, data map[string]interface{}) (err error) {
	return DB.Model(&token).Where("email = ?", email).Updates(data).Error
}

func (token *NonAuthToken) CreateNonAuthToken() (err error) {
	return DB.Create(&token).Error
}

func (token *NonAuthToken) GetNonAuthTokenByVerifyToken(verifyToken string) (err error) {
	return DB.Where("user_verify_token = ?", verifyToken).First(&token).Error
}

func (token *NonAuthToken) GetNonAuthTokenByPasswordChangeToken(passwordToken string) (err error) {
	return DB.Where("password_change_token = ?", passwordToken).First(&token).Error
}
