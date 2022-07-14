package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Username string `json:"username" gorm:"unique" gorm:"not null"`
	Email    string `json:"email" gorm:"unique" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
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

func (user *User) UpdateUserPassword(email string, password string) (err error) {
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

func (user *User) UpdateUserActive(email string) (err error) {
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	if err := DB.Where("email = ? AND is_active = ?", email, true).First(&user).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("user account is already activated")
	}
	user.IsActive = true
	if err := DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}
