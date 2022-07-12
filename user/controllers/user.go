package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go_practice/user/auth"
	"go_practice/user/logger"
	"go_practice/user/models"
	_ "go_practice/user/structs"
	"go_practice/user/utils"
	"net/http"
	"os"
)

// RegisterUser godoc
// @Summary      Register User
// @Description  create user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        input  body  structs.UserRequest  true  "Create User"
// @Success      201  {object}  structs.UserResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /user/register [post]
func (c *Controller) RegisterUser(context *gin.Context) {
	var user models.User
	var resetToken string
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

	if resetToken, err = auth.GenerateNonAuthJWT(user.Email); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "user can not be registered", err, logger.ERROR)
		return
	}
	link := "http://localhost:5002/api/v1/user/verify/?verify_token=" + resetToken
	body := "Here is your reset <a href='" + link + "'>link</a>"
	html := "<strong>" + body + "</strong>"
	fmt.Println(html)
	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email, "username": user.Username})
	context.Abort()
}

// VerifyUser godoc
// @Summary      Verify User
// @Description  Activate User
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        verify_token  query  string  true "verify_user" Format(string)
// @Success      201  {object}  structs.UserResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /user/verify [post]
func (c *Controller) VerifyUser(context *gin.Context) {
	verifyToken, err := context.GetQuery("verify_token")
	var user models.User
	if !err {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "token is not found", errors.New("token is not found"), logger.ERROR)
		return
	}
	_, claim := auth.ValidateNonAuthToken(verifyToken, []byte(os.Getenv("ANOTHER_TOKEN_SECRET")))
	if err := user.GetUserByEmail(claim.Email); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CustomErrorResponse(context, http.StatusBadRequest, "user is not found", err, logger.INFO)
			return
		}
		utils.CustomErrorResponse(context, http.StatusForbidden, "operation is not allowed", err, logger.ERROR)
		return
	}

	if err := user.UpdateUserActive(claim.Email); err != nil {
		utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
		return
	}

	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email, "username": user.Username})
	context.Abort()
}
