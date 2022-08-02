package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"go_practice/book/auth"
	elasticsearch2 "go_practice/book/elasticsearch"
	"go_practice/book/logger"
	"go_practice/book/utils"
	"log"
	"net/http"
	"strconv"
)

// GetElasticInfo godoc
// @Summary      Get Elastic Info
// @Description  get elastic details
// @Tags         elastic
// @Accept       json
// @Produce      json
// @Success      200  {object}  structs.ElasticJsonResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      401  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /elastic/info [get]
// @Security BearerAuth
func (c *Controller) GetElasticInfo(context *gin.Context) {
	es := context.MustGet("elastic").(*elasticsearch.Client)
	var r map[string]interface{}
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	context.JSON(http.StatusOK, gin.H{"data": r})
}

//SearchBooks godoc
//@Summary      Show Books by Searching for User
//@Description  get paginated list of books by search term
//@Tags         elastic
//@Accept       json
//@Produce      json
//@Param        page   query  int  false "paginate" Format(int)
//@Param        limit   query  int  false "paginate" Format(int)
//@Param        search   query  string  false "name searching" Format(string)
//@Success      200  {object}  structs.BookElasticPaginatedResponse
//@Failure      400  {object}  structs.ErrorResponse
//@Failure      401  {object}  structs.ErrorResponse
//@Failure      403  {object}  structs.ErrorResponse
//@Failure      404  {object}  structs.ErrorResponse
//@Failure      500  {object}  structs.ErrorResponse
//@Router       /elastic/books [get]
//@Security BearerAuth
func (c *Controller) SearchBooks(context *gin.Context) {

	//var err error
	var page, limit int
	var search string
	tokenString := context.GetHeader("Authorization")
	err, claim := auth.ValidateToken(tokenString)
	if err != nil {
		utils.BaseErrorResponse(context, http.StatusUnauthorized, err, logger.INFO)
		return
	}

	if page, err = strconv.Atoi(context.DefaultQuery("page", "1")); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "invalid 'page' param value type, Integer expected", err, logger.INFO)
		return
	}

	if limit, err = strconv.Atoi(context.DefaultQuery("limit", "10")); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "invalid 'limit' param value type, Integer expected", err, logger.INFO)
		return
	}
	var buf bytes.Buffer

	search = context.DefaultQuery("search", "")

	from := (page - 1) * limit
	r, err := elasticsearch2.GetPaginatedBookSearch(context, from, limit, search, buf, err)

	authorsData := utils.CreateBookListSearchResponse(r, claim.UserId)
	authorStructData := utils.CreateHyperBookElasticResponses(context, authorsData)
	paginatedResponse := utils.CreateHyperPaginatedBookResponses(page, limit, authorStructData)

	context.JSON(
		http.StatusOK,
		gin.H{"data": paginatedResponse},
	)
}

//SearchAuthors godoc
//@Summary      Show Authors by Searching for User
//@Description  get paginated list of authors by search term
//@Tags         elastic
//@Accept       json
//@Produce      json
//@Param        page   query  int  false "paginate" Format(int)
//@Param        limit   query  int  false "paginate" Format(int)
//@Param        search   query  string  false "name searching" Format(string)
//@Success      200  {object}  structs.AuthorElasticPaginatedResponse
//@Failure      400  {object}  structs.ErrorResponse
//@Failure      401  {object}  structs.ErrorResponse
//@Failure      403  {object}  structs.ErrorResponse
//@Failure      404  {object}  structs.ErrorResponse
//@Failure      500  {object}  structs.ErrorResponse
//@Router       /elastic/authors [get]
//@Security BearerAuth
func (c *Controller) SearchAuthors(context *gin.Context) {

	var err error
	var page, limit int
	var search string
	tokenString := context.GetHeader("Authorization")
	err, _ = auth.ValidateToken(tokenString)
	if err != nil {
		utils.BaseErrorResponse(context, http.StatusUnauthorized, err, logger.INFO)
		return
	}

	if page, err = strconv.Atoi(context.DefaultQuery("page", "1")); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "invalid 'page' param value type, Integer expected", err, logger.INFO)
		return
	}

	if limit, err = strconv.Atoi(context.DefaultQuery("limit", "10")); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "invalid 'limit' param value type, Integer expected", err, logger.INFO)
		return
	}
	var buf bytes.Buffer

	search = context.DefaultQuery("search", "")

	from := (page - 1) * limit
	//fmt.Println(search, from, limit)
	r, err := elasticsearch2.GetPaginatedAuthorSearch(context, from, limit, search, buf, err)

	authorsData := utils.CreateAuthorListSearchResponse(r)
	authorStructData := utils.CreateHyperAuthorResponses(context, authorsData)

	paginatedResponse := utils.CreateHyperPaginatedAuthorResponses(page, limit, authorStructData)

	context.JSON(
		http.StatusOK,
		gin.H{"data": paginatedResponse},
	)
}
