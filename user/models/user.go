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
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
}

func (user *User) GetUserByEmail(email string) (err error) {
	return DB.Where("email = ?", email).First(&user).Error
}

func (user *User) GetUserByID(ID uint) (err error) {
	return DB.Where("id = ?", ID).First(&user).Error
}

func (user *User) CreateUser() (err error) {
	return DB.Create(&user).Error
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
