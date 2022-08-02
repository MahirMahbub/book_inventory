package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go_practice/book/auth"
	"go_practice/book/logger"
	"go_practice/book/models"
	"go_practice/book/utils"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

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

	if err := author.GetAuthorWithBooks(uint(id)); err != nil {
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
