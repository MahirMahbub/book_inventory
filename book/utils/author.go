package utils

import (
	"github.com/gin-gonic/gin"
	"go_practice/book/models"
	"go_practice/book/structs"
	"strconv"
)

func CreateAuthorObjectResponse(ctx *gin.Context, author models.Author, isAdmin bool) structs.AuthorResponse {
	var authorResponse structs.AuthorResponse
	authorResponse.ID = author.ID
	authorResponse.FirstName = author.FirstName
	authorResponse.LastName = author.LastName
	authorResponse.Description = author.Description
	authorResponse.Books = CreateHyperBookResponsesForAuthor(ctx, author.Books, isAdmin)
	return authorResponse
}

func CreateHyperAuthorResponses(ctx *gin.Context, authors []structs.AuthorBase, isAdmin bool) []structs.HyperAuthorResponse {
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	apiPath := "api/v1/authors/"
	if isAdmin {
		apiPath = "api/v1/admin/authors/"
	}
	url := scheme + "://" + ctx.Request.Host + apiPath
	var bookResponses []structs.HyperAuthorResponse
	bookResponses = []structs.HyperAuthorResponse{}
	for _, author := range authors {
		bookResponse := CreateHyperAuthorResponse(author, url)
		bookResponses = append(bookResponses, bookResponse)
	}
	return bookResponses
}

func CreateHyperAuthorResponse(author structs.AuthorBase, url string) structs.HyperAuthorResponse {
	var authorResponse structs.HyperAuthorResponse
	authorResponse.ID = author.ID
	authorResponse.FirstName = author.FirstName
	authorResponse.LastName = author.LastName
	authorResponse.Url = url + strconv.Itoa(int(author.ID))
	return authorResponse
}

func CreateHyperPaginatedAuthorResponses(page int, limit int, authorStructData []structs.HyperAuthorResponse) structs.AuthorPaginated {
	var paginatedResponse structs.AuthorPaginated
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

func CreateAuthorListSearchResponse(input map[string]interface{}) []structs.AuthorBase {
	var authorsData []structs.AuthorBase
	authorsData = []structs.AuthorBase{}
	if len(input) > 0 {
		dataList := input["hits"].(map[string]interface{})["hits"].([]interface{})

		for i := 0; i < len(dataList); i++ {
			var authorDetails structs.AuthorBase
			sources := dataList[i].(map[string]interface{})["_source"].(map[string]interface{})
			authorDetails.ID = uint(sources["id"].(float64))
			authorDetails.FirstName = sources["first_name"].(string)
			authorDetails.LastName = sources["last_name"].(string)
			authorsData = append(authorsData, authorDetails)
		}
	}
	return authorsData
}
