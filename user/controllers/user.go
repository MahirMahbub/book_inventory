package controllers

import (
	"github.com/gin-gonic/gin"
	"go_practice/user/models"
	_ "go_practice/user/structs"
	"go_practice/user/utils"
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
		utils.BaseErrorResponse(context, http.StatusBadRequest, err)
		return
	}
	if err := user.HashPassword(user.Password); err != nil {
		utils.BaseErrorResponse(context, http.StatusBadRequest, err)
		return
	}
	if err := models.DB.Create(&user).Error; err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "User can not be registered", err)
		return
	}
	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email, "username": user.Username})
	context.Abort()
}

//func (c *Controller) ResetLink(context *gin.Context) {
//	var data forms.ResendCommand
//}
