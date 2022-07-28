package controllers

import (
	"bytes"
	"errors"
	"fmt"
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

// FindUserBooks godoc
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
func (c *Controller) FindUserBooks(context *gin.Context) {

	var books models.Books
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
// @Success      200  {object}  structs.BookAPIResponse
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
			utils.BaseErrorResponse(context, http.StatusNotFound, err, logger.INFO)
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
// @Success      201  {object}  structs.BookAPIResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      401  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /books [post]
// @Security BearerAuth
func (c *Controller) CreateBook(context *gin.Context) {
	var input structs.CreateBookInput
	var authors []models.Author
	var book models.Book

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

	for _, authorId := range input.AuthorIDs {
		var author models.Author

		if err := author.GetAuthorByID(authorId); err != nil {
			utils.CustomErrorResponse(context, http.StatusNotFound, "invalid Author!", err, logger.INFO)
			return
		}
		authors = append(authors, author)
	}

	book = models.Book{Title: input.Title, UserID: claim.UserId, Description: input.Description}

	if err := book.CreateUserBookWithAuthor(authors); err != nil {
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
// @Success      200      {object}  structs.BookAPIResponse
// @Failure      400      {object}  structs.ErrorResponse
// @Failure      401  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      404      {object}  structs.ErrorResponse
// @Failure      500      {object}  structs.ErrorResponse
// @Router       /books/{id} [patch]
// @Security BearerAuth
func (c *Controller) UpdateBook(context *gin.Context) {
	var id_ int
	var book models.Book
	var input structs.UpdateBookInput
	tokenString := context.GetHeader("Authorization")
	err, claim := auth.ValidateToken(tokenString)
	if err != nil {
		utils.BaseErrorResponse(context, http.StatusUnauthorized, err, logger.INFO)
		return
	}
	if id_, err = strconv.Atoi(context.Param("id")); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "book can not be updated, invalid id", err, logger.INFO)
		return
	}
	if err := book.GetUserBookByID(uint(id_), claim.UserId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
			return
		}
		utils.CustomErrorResponse(context, http.StatusForbidden, "operation is not allowed", err, logger.ERROR)
		return
	}

	if err := context.ShouldBindJSON(&input); err != nil {
		utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
		return
	}

	if err := book.UpdateBook(input); err != nil {
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
	var id_ int
	tokenString := context.GetHeader("Authorization")
	err, claim := auth.ValidateToken(tokenString)
	if err != nil {
		utils.BaseErrorResponse(context, http.StatusUnauthorized, err, logger.INFO)
		return
	}
	if id_, err = strconv.Atoi(context.Param("id")); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "book can not be updated, invalid id", err, logger.INFO)
		return
	}
	if err := book.GetUserBookByID(uint(id_), claim.UserId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
			return
		}
		utils.CustomErrorResponse(context, http.StatusForbidden, "operation is not allowed", err, logger.ERROR)
		return
	}

	if err := book.DeleteBook(); err != nil {
		utils.CustomErrorResponse(context, http.StatusForbidden, "book is not deleted", err, logger.ERROR)
		return
	}
	context.JSON(http.StatusNoContent, gin.H{"data": true})
}

//FindBooks godoc
//@Summary      Show Authors by Searching
//@Description  get paginated list of authors by search term
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
func (c *Controller) FindBooks(context *gin.Context) {

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
	fmt.Println(search, from, limit)
	r, err := elasticsearch.GetPaginatedBookSearch(context, from, limit, search, buf, err)

	authorsData := utils.CreateBookListSearchResponse(r)
	authorStructData := utils.CreateHyperBookElasticResponses(context, authorsData)
	//
	paginatedResponse := utils.CreateHyperPaginatedBookResponses(page, limit, authorStructData)

	context.JSON(
		http.StatusOK,
		gin.H{"data": paginatedResponse},
	)
}
