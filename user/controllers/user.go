package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go_practice/user/logger"
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
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /user/register [post]
func (c *Controller) RegisterUser(context *gin.Context) {
	var user models.User
	if err := context.ShouldBindJSON(&user); err != nil {
		utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
		return
	}
	err := user.GetUserByEmail(user.Email)

	if err == nil {
		utils.BaseErrorResponse(context, http.StatusForbidden, errors.New("user with email already registered"), logger.ERROR)
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "user can not be registered. operation is not allowed", err, logger.ERROR)
		return
	}
	if err := user.HashPassword(user.Password); err != nil {
		utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.ERROR)
		return
	}
	if err := user.CreateUser(); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "user can not be registered", err, logger.ERROR)
		return
	}
	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email, "username": user.Username})
	context.Abort()
}

//func (c *Controller) ResetLink(context *gin.Context) {
//	var data forms.ResendCommand
//}
