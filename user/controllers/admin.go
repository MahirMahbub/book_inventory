package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_practice/user/auth"
	"go_practice/user/logger"
	"go_practice/user/models"
	"go_practice/user/utils"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// CreateAdmin godoc
// @Summary      Create admin
// @Description  create admin
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        userId  path  int true "User ID" Format(int)
// @Success      200  {object}  structs.MessageResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      401  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /admin/{userId}/create-admin [post]
// @Security BearerAuth
func (c *Controller) CreateAdmin(context *gin.Context) {
	var id int
	var user models.User
	tokenString := context.GetHeader("Authorization")
	err, claim := auth.ValidateToken(tokenString)
	if err != nil {
		utils.BaseErrorResponse(context, http.StatusUnauthorized, err, logger.INFO)
		return
	}
	if id, err = strconv.Atoi(context.Param("userId")); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "book can not be updated, invalid id", err, logger.INFO)
		fmt.Println(err.Error(), context)
		return
	}
	if claim.IsAdmin {
		if err := user.CreateAdmin(uint(id)); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
				return
			}
			utils.CustomErrorResponse(context, http.StatusForbidden, "operation is not allowed", err, logger.ERROR)
			return
		}
	} else {
		utils.CustomErrorResponse(context, http.StatusForbidden, "user is not allowed, only admin is allowed", err, logger.INFO)
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "user is assigned as admin"})
}
