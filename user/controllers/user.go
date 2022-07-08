package controllers

import (
	"github.com/gin-gonic/gin"
	"go_practice/user/models"
	_ "go_practice/user/structs"
	"net/http"
)

// RegisterUser godoc
// @Summary      Register User
// @Description  create user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        input  body  structs.UserRequest  true  "Create User"
// @Success      200  {object}  structs.UserResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /user/register [post]
func (c *Controller) RegisterUser(context *gin.Context) {
	var user models.User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	record := models.DB.Create(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email, "username": user.Username})
}
