package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Username string `json:"username" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	IsActive bool   `default:"false"`
	IsAdmin  bool   `default:"false"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

func (user *User) GetUserByEmail(email string) (err error) {
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	return nil
}

func (user *User) GetUserByID(id uint) (err error) {

	if err := DB.Where("ID = ?", id).First(&user).Error; err != nil {
		return err
	}
	return nil
}

func (user *User) UpdateUserPass(email string, password string) (err error) {

	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	if err := user.HashPassword(password); err != nil {
		return err
	}
	if err := DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (user *User) VerifyAccount(email string) (err error) {
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	user.IsActive = true
	if err := DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}
