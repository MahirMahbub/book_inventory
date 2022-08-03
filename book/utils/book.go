package utils

import (
	"github.com/gin-gonic/gin"
	"go_practice/book/models"
	"go_practice/book/structs"
	"strconv"
)

func CreateBookResponse(ctx *gin.Context, book models.Book, isAdmin bool) structs.BookResponse {
	var bookResponse structs.BookResponse
	bookResponse.ID = book.ID
	bookResponse.Title = book.Title
	bookResponse.UserID = book.UserID
	bookResponse.Description = book.Description
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	apiPath := "api/v1/authors/"
	if isAdmin {
		apiPath = "api/v1/admin/authors/"
	}

	url := scheme + "://" + ctx.Request.Host + "/" + apiPath
	var authors []structs.HyperAuthorResponse
	for _, author := range book.Authors {
		tempAuthor := structs.AuthorBase{ID: author.ID,
			FirstName: author.FirstName, LastName: author.LastName}
		customAuthor := CreateHyperAuthorResponse(tempAuthor, url)

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

func CreateHyperBookElasticResponse(book structs.BookBase, url string) structs.HyperBookResponse {
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

func CreateHyperBookResponsesForAuthor(ctx *gin.Context, books []*models.Book, isAdmin bool) []structs.HyperBookResponse {
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	apiPath := "api/v1/books/"
	if isAdmin {
		apiPath = "api/v1/admin/books/"
	}

	url := scheme + "://" + ctx.Request.Host + "/" + apiPath
	var bookResponses []structs.HyperBookResponse
	for _, book := range books {
		bookResponse := CreateHyperBookResponse(*book, url)
		bookResponses = append(bookResponses, bookResponse)
	}
	return bookResponses
}

func CreateHyperBookElasticResponses(ctx *gin.Context, books []structs.BookBase) []structs.HyperBookResponse {

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	apiPath := "api/v1/books/"
	url := scheme + "://" + ctx.Request.Host + "/" + apiPath
	var bookResponses []structs.HyperBookResponse
	bookResponses = []structs.HyperBookResponse{}
	for _, book := range books {
		bookResponse := CreateHyperBookElasticResponse(book, url)
		bookResponses = append(bookResponses, bookResponse)
	}
	return bookResponses

}

func CreateBookListSearchResponse(input map[string]interface{}, userId uint) []structs.BookBase {
	var booksData []structs.BookBase
	booksData = []structs.BookBase{}
	if len(input) > 0 {
		dataList := input["hits"].(map[string]interface{})["hits"].([]interface{})

		for i := 0; i < len(dataList); i++ {
			var bookDetails structs.BookBase
			sources := dataList[i].(map[string]interface{})["_source"].(map[string]interface{})
			if sources["user_id"] == nil || sources["user_id"] == userId {
				bookDetails.ID = uint(sources["book_id"].(float64))
				bookDetails.Title = sources["title"].(string)
				booksData = append(booksData, bookDetails)
			}
		}
	}
	return booksData
}

func CreateHyperPaginatedBookResponses(page int, limit int, authorStructData []structs.HyperBookResponse) structs.BookPaginated {
	var paginatedResponse structs.BookPaginated
	paginatedResponse.Page = page
	paginatedResponse.Limit = limit
	if len(authorStructData) < limit {
		paginatedResponse.NextPage = 0
	} else {
		paginatedResponse.NextPage = page + 1
	}
	if page > 1 {
		paginatedResponse.PrevPage = page - 1
	} else {
		paginatedResponse.PrevPage = 0
	}
	paginatedResponse.Records = authorStructData
	return paginatedResponse
}
