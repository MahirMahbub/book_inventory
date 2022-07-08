package controllers

import (
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/gin-gonic/gin"
	"go_practice/book/models"
	"go_practice/book/structs"
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
func FindBooks(ctx *gin.Context) {
	var books []models.Book
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	paginator := pagination.Paging(&pagination.Param{
		DB:      models.DB.Find(&books),
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"id desc"},
		ShowSQL: true,
	}, &books)

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
func CreateBook(ctx *gin.Context) {
	var input structs.CreateBookInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var book models.Book
	book = models.Book{Title: input.Title, Author: input.Author}
	models.DB.Create(&book)

	ctx.JSON(http.StatusOK, gin.H{"data": book})
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
func FindBook(ctx *gin.Context) {
	var book models.Book
	if err := models.DB.Where("id = ?", ctx.Param("id")).First(&book).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": book})
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
func UpdateBook(ctx *gin.Context) {
	var book models.Book
	if err := models.DB.Where("id=?", ctx.Param("id")).First(&book).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	var input structs.UpdateBookInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Model(&book).Updates(input)
	ctx.JSON(http.StatusOK, gin.H{"data": book})
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
func DeleteBook(ctx *gin.Context) {
	var book models.Book
	if err := models.DB.Where("id=?", ctx.Param("id")).First(&book).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	models.DB.Delete(&book)

	ctx.JSON(http.StatusOK, gin.H{"data": true})
}