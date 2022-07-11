package controllers

import (
	"github.com/gin-gonic/gin"
	"go_practice/user/auth"
	"go_practice/user/models"
	"go_practice/user/structs"
	"go_practice/user/utils"
	"net/http"
)

// GenerateToken godoc
// @Summary      Generate Token
// @Description  post token
// @Tags         token
// @Accept       json
// @Produce      json
// @Param        input  body  structs.TokenRequest  true  "Create Token"
// @Success      200  {object}  structs.TokenResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /user/token [post]
func (c *Controller) GenerateToken(context *gin.Context) {
	var request structs.TokenRequest
	var user models.User
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.BaseErrorResponse(context, http.StatusBadRequest, err)
		return
	}

	//record := models.DB.Where("email = ?", request.Email).First(&user)
	if err := user.GetUserByEmail(request.Email); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "User not found", err)
		return
	}

	if err := user.CheckPassword(request.Password); err != nil {
		utils.CustomErrorResponse(context, http.StatusUnauthorized, "invalid credentials", err)
		return
	}
	tokenString, err := auth.GenerateJWT(user.Email, user.Username, user.ID)
	if err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "Token generation failed", err)
		return
	}
	context.JSON(http.StatusOK, gin.H{"token": "Bearer " + tokenString})
}
