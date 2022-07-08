package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Ping(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}
