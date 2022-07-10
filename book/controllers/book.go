package controllers

import (
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/gin-gonic/gin"
	"go_practice/book/auth"
	models "go_practice/book/models"
	"go_practice/book/structs"
	"go_practice/book/utils"
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
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /books [get]
// @Security BearerAuth
func (c *Controller) FindBooks(ctx *gin.Context) {
	var books []models.Book
	tokenString := ctx.GetHeader("Authorization")
	_, claim := auth.ValidateToken(tokenString)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	paginator := pagination.Paging(&pagination.Param{
		DB:      models.DB.Where("user_id = ?", claim.UserId).Select([]string{"id", "title"}).Find(&books),
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"id desc"},
		ShowSQL: true,
	}, &books)
	bookResponses := utils.CreateHyperBookResponses(ctx, books)

	paginator.Records = bookResponses
	ctx.JSON(http.StatusOK, gin.H{"data": paginator})

}

// CreateBook godoc
// @Summary      Add a book
// @Description  post book
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        input  body  structs.CreateBookInput  true  "Add books"
// @Success      200  {object}  structs.BookResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /books [post]
// @Security BearerAuth
func (c *Controller) CreateBook(ctx *gin.Context) {
	var input structs.CreateBookInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenString := ctx.GetHeader("Authorization")
	_, claim := auth.ValidateToken(tokenString)
	var authors []models.Author
	for _, authorId := range input.AuthorIDs {
		var author models.Author
		if err := models.DB.Where("id = ?", authorId).First(&author).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Author!"})
			return
		}
		authors = append(authors, author)
	}

	var book models.Book
	book = models.Book{Title: input.Title, UserID: claim.UserId, Description: input.Description}

	models.DB.Create(&book).Association("Authors").Append(authors)
	bookResponse := utils.CreateBookResponse(book)
	ctx.JSON(http.StatusOK, gin.H{"data": bookResponse})
}

// FindBook godoc
// @Summary      Show Books
// @Description  get books
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id  path  int true "Book ID" Format(int)
// @Success      200  {object}  structs.BookResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /books/{id} [get]
// @Security BearerAuth
func (c *Controller) FindBook(ctx *gin.Context) {
	var book models.Book
	tokenString := ctx.GetHeader("Authorization")
	_, claim := auth.ValidateToken(tokenString)
	if err := models.DB.Preload("Authors").Where("id = ? AND user_id = ?", ctx.Param("id"), claim.UserId).First(&book).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	bookResponse := utils.CreateBookResponse(book)
	ctx.JSON(http.StatusOK, gin.H{"data": bookResponse})
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
// @Failure      404      {object}  structs.ErrorResponse
// @Failure      500      {object}  structs.ErrorResponse
// @Router       /books/{id} [patch]
// @Security BearerAuth
func (c *Controller) UpdateBook(ctx *gin.Context) {
	var book models.Book
	tokenString := ctx.GetHeader("Authorization")
	_, claim := auth.ValidateToken(tokenString)
	if err := models.DB.Where("id = ? AND user_id = ?", ctx.Param("id"), claim.UserId).First(&book).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	var input structs.UpdateBookInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Model(&book).Updates(input)
	bookResponse := utils.CreateBookObjectResponse(book)
	ctx.JSON(http.StatusOK, gin.H{"data": bookResponse})
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
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /books/{id} [delete]
// @Security BearerAuth
func (c *Controller) DeleteBook(ctx *gin.Context) {
	var book models.Book
	tokenString := ctx.GetHeader("Authorization")
	_, claim := auth.ValidateToken(tokenString)
	if err := models.DB.Where("id = ? AND user_id = ?", ctx.Param("id"), claim.UserId).First(&book).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	models.DB.Delete(&book)

	ctx.JSON(http.StatusOK, gin.H{"data": true})
}
