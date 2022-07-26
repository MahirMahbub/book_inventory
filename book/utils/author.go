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

func CreateAuthorObjectResponse(author models.Author) structs.AuthorResponse {
	var authorResponse structs.AuthorResponse
	authorResponse.ID = author.ID
	authorResponse.FirstName = author.FirstName
	authorResponse.LastName = author.LastName
	authorResponse.Description = author.Description
	return authorResponse
}
