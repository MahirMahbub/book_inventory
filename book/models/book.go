package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title       string    `json:"title"`
	UserID      uint      `json:"userId" gorm:"default: null"`
	Description string    `gorm:"size:6000" json:"description"`
	Authors     []*Author `gorm:"many2many:authors_books;" json:"author"`
}

type Books []Book

type Author struct {
	gorm.Model
	FirstName   string  `json:"first_name" gorm:"not null" gorm:"uniqueIndex:idx_first_second"`
	LastName    string  `json:"last_name" gorm:"not null" gorm:"uniqueIndex:idx_first_second"`
	Description string  `gorm:"size:6000" json:"description"`
	Books       []*Book `gorm:"many2many:authors_books;" json:"books"`
}

func (book *Book) GetUserBookByID(ID uint, userID uint) (err error) {
	return DB.Where("id = ? AND user_id = ?", ID, userID).First(&book).Error
}

func (books *Books) GetUserBooksBySelection(userID uint, selection []string) *gorm.DB {
	return DB.Where("user_id = ?", userID).Select(selection).Find(&books)
}

func (book *Book) GetUserBookWithAuthor(ID uint, userID uint) (err error) {
	return DB.Preload("Authors").Where("id = ? AND user_id = ?", ID, userID).First(&book).Error
}

//func (book *Book) PostBook
