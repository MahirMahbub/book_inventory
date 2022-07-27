package models

import (
	"go_practice/book/structs"
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

func (book *Book) GetUserBookByID(ID uint, userID uint) (err error) {
	return DB.Where("id = ? AND user_id = ?", ID, userID).First(&book).Error
}

func (books *Books) GetUserBooksBySelection(userID uint, selection []string) *gorm.DB {
	return DB.Where("user_id = ?", userID).Select(selection).Find(&books)
}

func (book *Book) GetUserBookWithAuthor(ID uint, userID uint) (err error) {
	return DB.Preload("Authors").Where("id = ? AND user_id = ?", ID, userID).First(&book).Error
}

func (book *Book) CreateUserBookWithAuthor(authors []Author) (err error) {
	return DB.Create(&book).Association("Authors").Append(authors)
}

func (book *Book) UpdateBook(input structs.UpdateBookInput) (err error) {
	return DB.Model(&book).Updates(input).Error
}

func (book *Book) DeleteBook() (err error) {
	return DB.Delete(&book).Error
}
