package utils

import (
	"github.com/gin-gonic/gin"
	"go_practice/book/models"
	"go_practice/book/structs"
	"strconv"
)

func CreateBookResponse(book models.Book) structs.BookResponse {
	var bookResponse structs.BookResponse
	bookResponse.ID = book.ID
	bookResponse.Title = book.Title
	bookResponse.UserID = book.UserID
	bookResponse.Description = book.Description
	var authors []structs.AuthorBasicResponse
	for _, author := range book.Authors {
		customAuthor := CreateBasicAuthorResponse(author)
		authors = append(authors, customAuthor)
	}
	bookResponse.Authors = authors
	return bookResponse
}

func CreateBookObjectResponse(book models.Book) structs.BookUpdateResponse {
	var bookResponse structs.BookUpdateResponse
	bookResponse.ID = book.ID
	bookResponse.Title = book.Title
	bookResponse.UserID = book.UserID
	bookResponse.Description = book.Description
	return bookResponse
}

func CreateHyperBookResponse(book models.Book, url string) structs.HyperBookResponse {
	var bookResponse structs.HyperBookResponse
	bookResponse.ID = book.ID
	bookResponse.Title = book.Title
	bookResponse.Url = url + strconv.Itoa(int(book.ID))
	return bookResponse
}

func CreateHyperBookResponses(ctx *gin.Context, books []models.Book) []structs.HyperBookResponse {
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	url := scheme + "://" + ctx.Request.Host + ctx.Request.URL.Path + "/"
	var bookResponses []structs.HyperBookResponse
	for _, book := range books {
		bookResponse := CreateHyperBookResponse(book, url)
		bookResponses = append(bookResponses, bookResponse)
	}
	return bookResponses
}
