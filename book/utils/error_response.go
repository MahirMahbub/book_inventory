package utils

import (
	"github.com/gin-gonic/gin"
	"go_practice/book/logger"
)

func BaseErrorResponse(context *gin.Context, errorCode int, err error, logType string) {
	logger.PrintLog(logType, err)
	context.JSON(errorCode, gin.H{"error": err.Error()})
	context.Abort()
	return
}

func CustomErrorResponse(context *gin.Context, errorCode int, errString string, err error, logType string) {
	logger.PrintLog(logType, err)
	context.JSON(errorCode, gin.H{"error": errString})
	context.Abort()
	return
}
