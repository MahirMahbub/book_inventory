package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go_practice/user/auth"
	"go_practice/user/logger"
	"go_practice/user/models"
	"go_practice/user/services"
	"go_practice/user/structs"
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
// @Success      200  {object}  structs.MessageResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /user/register [post]
func (c *Controller) RegisterUser(context *gin.Context) {
	var user models.User
	var token models.NonAuthToken
	var resetToken string
	if err := context.ShouldBindJSON(&user); err != nil {
		utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
		return
	}
	err := user.GetUserByEmail(user.Email)

	if err == nil {
		utils.BaseErrorResponse(context, http.StatusForbidden,
			errors.New("user with email already registered"), logger.ERROR)
		return
	}
	//if !errors.Is(err, gorm.ErrRecordNotFound) {
	//	utils.CustomErrorResponse(context, http.StatusBadRequest,
	//		"user can not be registered. operation is not allowed", err, logger.ERROR)
	//	return
	//}
	if err := user.HashPassword(user.Password); err != nil {
		utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.ERROR)
		return
	}
	if err := user.CreateUser(); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest,
			"user can not be registered", err, logger.ERROR)
		return
	}

	if resetToken, err = auth.GenerateNonAuthJWT(user.Email); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest,
			"user can not be registered", err, logger.ERROR)
		return
	}
	token.UserVerifyToken = resetToken
	token.Email = user.Email
	err = token.CreateNonAuthToken()
	if err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest,
			"verify token generation failed", err, logger.ERROR)
		return
	}
	services.SendVerifyEmail(resetToken)
	context.JSON(http.StatusOK, gin.H{"message": "verify email has been send"})
}

// ResendUserVerifyEmail godoc
// @Summary      Resend Verify User Email
// @Description  Reactivate User
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        input  body  structs.TokenRequest  true  "Create Token"
// @Success      200  {object}  structs.MessageResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /user/resend-verify-token [post]
func (c *Controller) ResendUserVerifyEmail(context *gin.Context) {
	var user models.User
	var resetToken string
	var token models.NonAuthToken
	var request structs.TokenRequest
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
	err := user.GetUserByEmail(user.Email)

	if resetToken, err = auth.GenerateNonAuthJWT(user.Email); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest,
			"user can not be registered", err, logger.ERROR)
		return
	}
	err_ := token.UpdateNonAuthToken(request.Email, map[string]interface{}{"user_verify_token": resetToken})
	if err_ != nil {
		logger.Error.Println("old verify token can not be deactivated, potential risks")
		err_ = errors.New("verify token generation failed")
		return
	}
	services.SendVerifyEmail(resetToken)
	context.JSON(http.StatusOK, gin.H{"message": "verify email has been resend"})
}

// SendPasswordChangeEmail godoc
// @Summary      Send Password Change Email
// @Description  Change Password
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        input  body  structs.PasswordChangeTokenRequest  true  "Create Password Change Token"
// @Success      200  {object}  structs.MessageResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /user/send-password-change-token [post]
func (c *Controller) SendPasswordChangeEmail(context *gin.Context) {
	var user models.User
	var resetToken string
	var token models.NonAuthToken
	var request structs.PasswordChangeTokenRequest
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
	//if err := user.CheckPassword(request.Password); err != nil {
	//	utils.CustomErrorResponse(context, http.StatusUnauthorized, "invalid credentials", err, logger.INFO)
	//	return
	//}
	err := user.GetUserByEmail(user.Email)

	if resetToken, err = auth.GenerateNonAuthJWT(user.Email); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest,
			"user can not be registered", err, logger.ERROR)
		return
	}
	err_ := token.UpdateNonAuthToken(request.Email, map[string]interface{}{"password_change_token": resetToken})
	if err_ != nil {
		logger.Error.Println("old password token can not be deactivated, potential risks")
		err_ = errors.New("password change token generation failed")
		return
	}
	services.SendPasswordChangeEmail(resetToken)
	context.JSON(http.StatusOK, gin.H{"message": "password change email has been send"})
}

// VerifyUser godoc
// @Summary      Verify User
// @Description  Activate User
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        verify_token  query  string  true "verify_user" Format(string)
// @Success      200  {object}  structs.MessageResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /user/verify [put]
func (c *Controller) VerifyUser(context *gin.Context) {
	verifyToken, err := context.GetQuery("verify_token")
	var user models.User
	var token models.NonAuthToken
	if !err {
		utils.CustomErrorResponse(context, http.StatusBadRequest,
			"token is not found", errors.New("token is not found"), logger.INFO)
		return
	}
	_, claim := auth.ValidateNonAuthToken(verifyToken, []byte(os.Getenv("ANOTHER_TOKEN_SECRET")))
	if err := user.GetUserByEmail(claim.Email); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CustomErrorResponse(context, http.StatusBadRequest, "user is not found", err, logger.INFO)
			return
		}
		//utils.CustomErrorResponse(context, http.StatusForbidden, "operation is not allowed", err, logger.ERROR)
		//return
	}

	err_ := token.GetNonAuthTokenByVerifyToken(verifyToken)
	if gorm.IsRecordNotFoundError(err_) {
		err_ = errors.New("old verify token provided")
		utils.CustomErrorResponse(context, http.StatusBadRequest,
			"this link is invalid or old.", err_, logger.INFO)
		return
	}
	if err_ != nil {
		utils.CustomErrorResponse(context, http.StatusForbidden, "operation is not allowed", err_, logger.ERROR)
		return
	}

	if err := user.UpdateUserActive(claim.Email); err != nil {
		fmt.Println(err)
		utils.CustomErrorResponse(context, http.StatusForbidden, "can not active, operation is not allowed", err, logger.ERROR)
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "user is verified"})
}

// ChangePassword godoc
// @Summary      Change the user account password
// @Description  Change Password
// @Tags         user
// @Accept       json
// @Produce      json
// @Param		 verify_token query string true "Verify Token"
// @Param        input  body  structs.ChangePasswordRequest  true  "Change Password"
// @Success      200  {object}  structs.MessageResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /user/change-password [put]
func (c *Controller) ChangePassword(context *gin.Context) {
	verifyToken, err := context.GetQuery("verify_token")
	var user models.User
	var token models.NonAuthToken
	if !err {
		utils.CustomErrorResponse(context, http.StatusBadRequest,
			"token is not found", errors.New("token is not found"), logger.INFO)
		return
	}
	var request structs.ChangePasswordRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
		return
	}
	if request.Password != request.Confirm {
		err_ := errors.New("passwords do not match")
		utils.BaseErrorResponse(context, http.StatusBadRequest, err_, logger.INFO)
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

	err_ := token.GetNonAuthTokenByPasswordChangeToken(verifyToken)
	if gorm.IsRecordNotFoundError(err_) {
		err_ = errors.New("old verify token provided")
		utils.CustomErrorResponse(context, http.StatusBadRequest,
			"this link is invalid or old.", err_, logger.INFO)
		return
	}
	if err_ != nil {
		utils.CustomErrorResponse(context, http.StatusForbidden, "operation is not allowed", err_, logger.ERROR)
		return
	}

	if err := user.UpdateUserPassword(claim.Email, request.Password); err != nil {
		utils.CustomErrorResponse(context, http.StatusForbidden, "password can not be changed", err, logger.ERROR)
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "user password is changed"})
}
