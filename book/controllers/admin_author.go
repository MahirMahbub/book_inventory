package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_practice/book/auth"
	"go_practice/book/logger"
	"go_practice/book/models"
	"go_practice/book/structs"
	"go_practice/book/utils"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// CreateAdminAuthor godoc
// @Summary      Add an Author by Admin
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
func (c *Controller) CreateAdminAuthor(context *gin.Context) {
	var input structs.CreateAuthorInput

	if err := context.ShouldBindJSON(&input); err != nil {
		utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
		return
	}

	tokenString := context.GetHeader("Authorization")
	err, claim := auth.ValidateToken(tokenString)
	if err != nil {
		utils.BaseErrorResponse(context, http.StatusUnauthorized, err, logger.INFO)
		return
	}
	var author models.Author
	author = models.Author{FirstName: input.FirstName, LastName: input.LastName, Description: input.Description}

	if err := author.CreateAuthor(); err != nil {
		utils.CustomErrorResponse(context, http.StatusForbidden, "author is not created", err, logger.ERROR)
		return
	}
	authorResponse := utils.CreateAuthorObjectResponse(context, author, claim.IsAdmin)
	context.JSON(http.StatusCreated, gin.H{"data": authorResponse})
}

// FindAdminAuthor godoc
// @Summary      Show Author Details to Admin
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
// @Router       /admin/authors/{id} [get]
// @Security BearerAuth
func (c *Controller) FindAdminAuthor(context *gin.Context) {
	var author models.Author
	var id uint64
	var err error
	tokenString := context.GetHeader("Authorization")
	err, claim := auth.ValidateToken(tokenString)
	if err != nil {
		utils.BaseErrorResponse(context, http.StatusUnauthorized, err, logger.INFO)
		return
	}

	if id, err = strconv.ParseUint(context.Param("id"), 10, 32); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "invalid 'id' param value type, Integer expected", err, logger.INFO)
		return
	}

	if err := author.GetAuthorWithBooks(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.BaseErrorResponse(context, http.StatusNotFound, err, logger.INFO)
			return
		}
		utils.CustomErrorResponse(context, http.StatusForbidden, "operation is not allowed", err, logger.ERROR)
		return
	}
	authorResponse := utils.CreateAuthorObjectResponse(context, author, claim.IsAdmin)
	context.JSON(http.StatusOK, gin.H{"data": authorResponse})
}

// FindAdminAuthors godoc
// @Summary      Show Authors to Admin
// @Description  gets all authors
// @Tags         admin/authors
// @Accept       json
// @Produce      json
// @Param        page   query  int  false "paginate" Format(int)
// @Param        limit   query  int  false "paginate" Format(int)
// @Success      200  {object}  structs.AuthorsPaginatedResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      401  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /admin/authors [get]
// @Security BearerAuth
func (c *Controller) FindAdminAuthors(context *gin.Context) {

	var authors models.Authors
	var err error
	var page, limit int
	var db *gorm.DB
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

	if db = authors.GetAuthorsBySelection([]string{"id", "first_name", "last_name"}); db.Error != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
			return
		}
		utils.CustomErrorResponse(context, http.StatusForbidden, "operation is not allowed", err, logger.ERROR)
		return
	}
	paginator := utils.Paging(&utils.Param{
		DB:      db,
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"id desc"},
		ShowSQL: true,
	}, &authors)

	var authorsList []structs.AuthorBase
	for i := 0; i < len(authors); i++ {
		tempAuthor := structs.AuthorBase{
			ID:        authors[i].ID,
			FirstName: authors[i].FirstName,
			LastName:  authors[i].LastName,
		}
		authorsList = append(authorsList, tempAuthor)
	}
	authorResponses := utils.CreateHyperAuthorResponses(context, authorsList, claim.IsAdmin)

	paginator.Records = authorResponses
	context.JSON(http.StatusOK, gin.H{"data": paginator})
}

// UpdateAdminAuthor godoc
// @Summary      Update an Author by Admin
// @Description  patch author
// @Tags         admin/authors
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Author ID" Format(int)
// @Param        input  body  structs.UpdateAuthorInput  false  "Update authors"
// @Success      200      {object}  structs.AuthorAPIResponse
// @Failure      400      {object}  structs.ErrorResponse
// @Failure      401  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      404      {object}  structs.ErrorResponse
// @Failure      500      {object}  structs.ErrorResponse
// @Router       /admin/authors/{id} [patch]
// @Security BearerAuth
func (c *Controller) UpdateAdminAuthor(context *gin.Context) {
	var id_ int
	var author models.Author
	var input structs.UpdateAuthorInput
	tokenString := context.GetHeader("Authorization")
	err, claim := auth.ValidateToken(tokenString)
	if err != nil {
		utils.BaseErrorResponse(context, http.StatusUnauthorized, err, logger.INFO)
		return
	}
	if id_, err = strconv.Atoi(context.Param("id")); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "author can not be updated, invalid id", err, logger.INFO)
		return
	}
	if err := author.GetAuthorWithBooks(uint(id_)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
			return
		}
		utils.CustomErrorResponse(context, http.StatusForbidden, "operation is not allowed", err, logger.ERROR)
		return
	}
	fmt.Println(author)
	if err := context.ShouldBindJSON(&input); err != nil {
		utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
		return
	}

	if err := author.UpdateAuthor(input); err != nil {
		utils.CustomErrorResponse(context, http.StatusForbidden, "author is not updated", err, logger.ERROR)
		return
	}
	bookResponse := utils.CreateAuthorObjectResponse(context, author, claim.IsAdmin)
	context.JSON(http.StatusOK, gin.H{"data": bookResponse})
}

// DeleteAdminAuthor godoc
// @Summary      Delete an Author by Admin
// @Description  delete author
// @Tags         admin/authors
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Book ID"  Format(int)
// @Success      204  {object}  structs.AuthorDeleteResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      401  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /admin/authors/{id} [delete]
// @Security BearerAuth
func (c *Controller) DeleteAdminAuthor(context *gin.Context) {
	var author models.Author
	var id_ int
	tokenString := context.GetHeader("Authorization")
	err, _ := auth.ValidateToken(tokenString)
	if err != nil {
		utils.BaseErrorResponse(context, http.StatusUnauthorized, err, logger.INFO)
		return
	}
	if id_, err = strconv.Atoi(context.Param("id")); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "author can not be updated, invalid id", err, logger.INFO)
		return
	}

	if err := author.GetAuthorByID(uint(id_)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
			return
		}
		utils.CustomErrorResponse(context, http.StatusForbidden, "operation is not allowed", err, logger.ERROR)
		return
	}

	if err := author.DeleteAuthor(); err != nil {
		utils.CustomErrorResponse(context, http.StatusForbidden, "author is not deleted", err, logger.ERROR)
		return
	}
	context.JSON(http.StatusNoContent, gin.H{"data": true})
}
