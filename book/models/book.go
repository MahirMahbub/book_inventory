package models

import "github.com/jinzhu/gorm"

type Book struct {
	gorm.Model
	Title       string    `json:"title"`
	UserID      uint      `json:"userId"`
	Description string    `gorm:"size:6000" json:"description"`
	Authors     []*Author `gorm:"many2many:authors_books;" json:"author"`
}

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
