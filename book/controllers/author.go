package controllers

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"go_practice/book/auth"
	"go_practice/book/elasticsearch"
	"go_practice/book/logger"
	"go_practice/book/models"
	"go_practice/book/structs"
	"go_practice/book/utils"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// CreateAuthor godoc
// @Summary      Add an author
// @Description  post author by admin
// @Tags         admin/authors
// @Accept       json
// @Produce      json
// @Param        input  body  structs.CreateAuthorInput  true  "Add authors"
// @Success      201  {object}  structs.AuthorResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      401  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /admin/authors [post]
// @Security BearerAuth
func (c *Controller) CreateAuthor(context *gin.Context) {
	var input structs.CreateAuthorInput

	if err := context.ShouldBindJSON(&input); err != nil {
		utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
		return
	}

	tokenString := context.GetHeader("Authorization")
	err, _ := auth.ValidateToken(tokenString)
	if err != nil {
		utils.BaseErrorResponse(context, http.StatusUnauthorized, err, logger.INFO)
		return
	}
	var author models.Author
	author = models.Author{FirstName: input.FirstName, LastName: input.LastName, Description: input.Description}

	if err := author.CreateBook(); err != nil {
		utils.CustomErrorResponse(context, http.StatusForbidden, "author is not created", err, logger.ERROR)
		return
	}
	authorResponse := utils.CreateAuthorObjectResponse(author)
	context.JSON(http.StatusCreated, gin.H{"data": authorResponse})
}

//FindAuthors godoc
//@Summary      Show Authors
//@Description  get authors
//@Tags         authors
//@Accept       json
//@Produce      json
//@Param        page   query  int  false "paginate" Format(int)
//@Param        limit   query  int  false "paginate" Format(int)
//@Param        search   query  string  false "name searching" Format(string)
//@Success      200  {object}  structs.AuthorPaginatedResponse
//@Failure      400  {object}  structs.ErrorResponse
//@Failure      401  {object}  structs.ErrorResponse
//@Failure      403  {object}  structs.ErrorResponse
//@Failure      404  {object}  structs.ErrorResponse
//@Failure      500  {object}  structs.ErrorResponse
//@Router       /authors [get]
//@Security BearerAuth
func (c *Controller) FindAuthors(context *gin.Context) {

	//var books models.Author
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
	r, err := elasticsearch.GetPaginatedAuthorSearch(context, from, limit, search, buf, err)

	authorsData := utils.CreateAuthorListSearchResponse(r)
	authorStructData := utils.CreateHyperAuthorResponses(context, authorsData)

	paginatedResponse := utils.CreateHyperPaginatedAuthorResponses(page, limit, authorStructData)

	context.JSON(
		http.StatusOK,
		gin.H{"data": paginatedResponse},
	)
}

// FindAuthor godoc
// @Summary      Show Author Details
// @Description  get author
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        id  path  int true "Author ID" Format(int)
// @Success      200  {object}  structs.AuthorAPIResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      401  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /authors/{id} [get]
// @Security BearerAuth
func (c *Controller) FindAuthor(context *gin.Context) {
	var author models.Author
	var id uint64
	var err error
	tokenString := context.GetHeader("Authorization")
	err, _ = auth.ValidateToken(tokenString)
	if err != nil {
		utils.BaseErrorResponse(context, http.StatusUnauthorized, err, logger.INFO)
		return
	}

	if id, err = strconv.ParseUint(context.Param("id"), 10, 32); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "invalid 'id' param value type, Integer expected", err, logger.INFO)
		return
	}

	if err := author.GetAuthorByID(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.BaseErrorResponse(context, http.StatusNotFound, err, logger.INFO)
			return
		}
		utils.CustomErrorResponse(context, http.StatusForbidden, "operation is not allowed", err, logger.ERROR)
		return
	}
	authorResponse := utils.CreateAuthorObjectResponse(author)
	context.JSON(http.StatusOK, gin.H{"data": authorResponse})
}
