package models

import (
	"github.com/jinzhu/gorm"
)

type Book struct {
	gorm.Model
	Title       string    `json:"title"`
	UserID      uint      `json:"userId"`
	Description string    `gorm:"size:6000" json:"description"`
	Authors     []*Author `gorm:"many2many:authors_books;" json:"author"`
}

type Books []Book

type Author struct {
	gorm.Model
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Description string  `gorm:"size:6000" json:"description"`
	Books       []*Book `gorm:"many2many:authors_books;" json:"books"`
	Contact     []*AuthorContact
}
type AuthorContact struct {
	gorm.Model
	Platform string `json:"platform"`
	URL      string `json:"url"`
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
