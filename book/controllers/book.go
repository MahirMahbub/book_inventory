package controllers

import (
	"encoding/json"
	"errors"
	es7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"go_practice/book/auth"
	"go_practice/book/logger"
	"go_practice/book/models"
	"go_practice/book/structs"
	"go_practice/book/utils"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// FindBooks godoc
// @Summary      Show Books
// @Description  get books
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        page   query  int  false "paginate" Format(int)
// @Param        limit   query  int  false "paginate" Format(int)
// @Success      200  {object}  structs.BooksPaginatedResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      401  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /books [get]
// @Security BearerAuth
func (c *Controller) FindBooks(context *gin.Context) {

	var books models.Books
	var err error
	var page, limit int
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
	var db *gorm.DB

	if db = books.GetUserBooksBySelection(claim.UserId, []string{"id", "title"}); db.Error != nil {
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
	}, &books)
	bookResponses := utils.CreateHyperBookResponses(context, books)

	paginator.Records = bookResponses
	context.JSON(http.StatusOK, gin.H{"data": paginator})
}

// FindBook godoc
// @Summary      Show Book Details
// @Description  get book
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id  path  int true "Book ID" Format(int)
// @Success      200  {object}  structs.BookResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      401  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /books/{id} [get]
// @Security BearerAuth
func (c *Controller) FindBook(context *gin.Context) {
	var book models.Book
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

	if err := book.GetUserBookWithAuthor(uint(id), claim.UserId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
			return
		}
		utils.CustomErrorResponse(context, http.StatusForbidden, "operation is not allowed", err, logger.ERROR)
		return
	}
	bookResponse := utils.CreateBookResponse(book)
	context.JSON(http.StatusOK, gin.H{"data": bookResponse})
}

// CreateBook godoc
// @Summary      Add a book
// @Description  post book
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        input  body  structs.CreateBookInput  true  "Add books"
// @Success      201  {object}  structs.BookResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      401  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /books [post]
// @Security BearerAuth
func (c *Controller) CreateBook(context *gin.Context) {
	var input structs.CreateBookInput

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
	var authors []models.Author

	for _, authorId := range input.AuthorIDs {
		var author models.Author

		if err := models.DB.Where("id = ?", authorId).First(&author).Error; err != nil {
			utils.CustomErrorResponse(context, http.StatusNotFound, "invalid Author!", err, logger.INFO)
			return
		}
		authors = append(authors, author)
	}

	var book models.Book
	book = models.Book{Title: input.Title, UserID: claim.UserId, Description: input.Description}

	if err := models.DB.Create(&book).Association("Authors").Append(authors); err != nil {
		utils.CustomErrorResponse(context, http.StatusForbidden, "book is not created", err, logger.ERROR)
		return
	}
	bookResponse := utils.CreateBookResponse(book)
	context.JSON(http.StatusCreated, gin.H{"data": bookResponse})
}

// UpdateBook godoc
// @Summary      Update a book
// @Description  patch book
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Account ID"
// @Param        input  body  structs.UpdateBookInput  false  "Update books"
// @Success      200      {object}  structs.BookResponse
// @Failure      400      {object}  structs.ErrorResponse
// @Failure      401  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      404      {object}  structs.ErrorResponse
// @Failure      500      {object}  structs.ErrorResponse
// @Router       /books/{id} [patch]
// @Security BearerAuth
func (c *Controller) UpdateBook(context *gin.Context) {
	var book models.Book
	tokenString := context.GetHeader("Authorization")
	err, claim := auth.ValidateToken(tokenString)
	if err != nil {
		utils.BaseErrorResponse(context, http.StatusUnauthorized, err, logger.INFO)
		return
	}

	if err := models.DB.Where("id = ? AND user_id = ?", context.Param("id"), claim.UserId).First(&book).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
			return
		}
		utils.CustomErrorResponse(context, http.StatusForbidden, "operation is not allowed", err, logger.ERROR)
		return
	}

	var input structs.UpdateBookInput

	if err := context.ShouldBindJSON(&input); err != nil {
		utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
		return
	}

	if err := models.DB.Model(&book).Updates(input).Error; err != nil {
		utils.CustomErrorResponse(context, http.StatusForbidden, "book is not updated", err, logger.ERROR)
		return
	}
	bookResponse := utils.CreateBookObjectResponse(book)
	context.JSON(http.StatusOK, gin.H{"data": bookResponse})

}

// DeleteBook godoc
// @Summary      Delete a book
// @Description  delete book
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Book ID"  Format(int)
// @Success      204  {object}  structs.BookDeleteResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      401  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /books/{id} [delete]
// @Security BearerAuth
func (c *Controller) DeleteBook(context *gin.Context) {
	var book models.Book
	tokenString := context.GetHeader("Authorization")
	err, claim := auth.ValidateToken(tokenString)
	if err != nil {
		utils.BaseErrorResponse(context, http.StatusUnauthorized, err, logger.INFO)
		return
	}

	if err := models.DB.Where("id = ? AND user_id = ?", context.Param("id"), claim.UserId).First(&book).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
			return
		}
		utils.CustomErrorResponse(context, http.StatusForbidden, "operation is not allowed", err, logger.ERROR)
		return
	}

	if err := models.DB.Delete(&book).Error; err != nil {
		utils.CustomErrorResponse(context, http.StatusForbidden, "book is not deleted", err, logger.ERROR)
		return
	}

	context.JSON(http.StatusNoContent, gin.H{"data": true})
}

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
	es := context.MustGet("elastic").(*es7.Client)
	var r map[string]interface{}
	//fmt.Println(client.Index())
	//fmt.Println(es7.Count())
	//log.Println(es7.Info())
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	context.JSON(http.StatusOK, gin.H{"data": r})
}
