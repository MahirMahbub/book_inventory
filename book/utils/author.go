package utils

import (
	"go_practice/book/models"
	"go_practice/book/structs"
)

func CreateBasicAuthorResponse(author *models.Author) structs.AuthorBasicResponse {
	var customAuthor structs.AuthorBasicResponse
	customAuthor.ID = author.ID
	customAuthor.Name = author.FirstName + " " + author.LastName
	return customAuthor
}
