package models

import (
	"go_practice/book/structs"
	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	FirstName   string  `json:"first_name" gorm:"not null" gorm:"uniqueIndex:idx_first_second"`
	LastName    string  `json:"last_name" gorm:"not null" gorm:"uniqueIndex:idx_first_second"`
	Description string  `gorm:"size:6000" json:"description"`
	Books       []*Book `gorm:"many2many:authors_books;" json:"books"`
}
type Authors []Author

func (author *Author) GetAuthorByID(ID uint) (err error) {
	return DB.Where("id = ?", ID).First(&author).Error
}

func (author *Author) GetAuthorWithBooks(ID uint) (err error) {
	return DB.Preload("Books").Where("id = ?", ID).First(&author).Error
}

func (author *Author) CreateAuthor() (err error) {
	return DB.Create(&author).Error
}

func (authors *Authors) GetAuthorsBySelection(selection []string) *gorm.DB {
	return DB.Select(selection).Find(&authors)
}

func (author *Author) UpdateAuthor(input structs.UpdateAuthorInput) (err error) {
	if input.FirstName != "" {
		author.FirstName = input.FirstName
	}

	if input.LastName != "" {
		author.LastName = input.LastName
	}

	if input.Description != "" {
		author.Description = input.Description
	}
	if err := DB.Save(&author).Error; err != nil {
		return err
	}
	return nil
}

func (author *Author) DeleteAuthor() (err error) {
	return DB.Delete(&author).Error
}
