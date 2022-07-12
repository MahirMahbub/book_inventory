package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go_practice/user/auth"
	"go_practice/user/logger"
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
// @Failure      401  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /user/token [post]
func (c *Controller) GenerateToken(context *gin.Context) {
	var request structs.TokenRequest
	var user models.User
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
		return
	}

	if err := user.GetUserByEmail(request.Email); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CustomErrorResponse(context, http.StatusBadRequest, "user is not found", err, logger.INFO)
			return
		}
		utils.CustomErrorResponse(context, http.StatusForbidden, "operation is not allowed", err, logger.ERROR)
		return
	}

	if err := user.CheckPassword(request.Password); err != nil {
		utils.CustomErrorResponse(context, http.StatusUnauthorized, "invalid credentials", err, logger.INFO)
		return
	}
	tokenString, refreshToken, err := auth.GenerateJWT(user.Email, user.Username, user.ID, user.IsAdmin, user.IsActive)
	if err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "Token generation failed", err, logger.ERROR)
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"token":        "Bearer " + tokenString,
		"refreshToken": "Bearer " + refreshToken,
	})
}

//// RefreshToken godoc
//// @Summary      Refresh Token
//// @Description  refreshes the access token
//// @Tags         user
//// @Accept       json
//// @Produce      json
//// @Param        refresh_token  query  string  true "refresh_token" Format(string)
//// @Success      200  {object}  structs.TokenResponse
//// @Failure      400  {object}  structs.ErrorResponse
//// @Failure      404  {object}  structs.ErrorResponse
//// @Failure      401  {object}  structs.ErrorResponse
//// @Failure      403  {object}  structs.ErrorResponse
//// @Failure      500  {object}  structs.ErrorResponse
//// @Router       /user/refresh-token [post]
//func (c *Controller) RefreshToken(context *gin.Context) {
//	verifyToken, err := context.GetQuery("refresh_token")
//	if !err {
//		utils.CustomErrorResponse(context, http.StatusBadRequest, "refresh token is not found", errors.New("refresh token is not found"), logger.INFO)
//		return
//	}
//	_, claim := auth.ValidateToken(verifyToken, []byte(os.Getenv("REFRESH_TOKEN_SECRET")))
//	tokenString, refreshToken, err_ := auth.GenerateJWT(claim.Email, claim.Username, claim.UserId, claim.IsAdmin, claim.IsActive)
//	if err_ != nil {
//		utils.CustomErrorResponse(context, http.StatusBadRequest, "Token generation failed", err_, logger.ERROR)
//		return
//	}
//	context.JSON(http.StatusOK, gin.H{
//		"token":        "Bearer " + tokenString,
//		"refreshToken": "Bearer " + refreshToken,
//	})
//}
