package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go_practice/user/controllers"
	"go_practice/user/logger"
	"go_practice/user/middlewares"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func SetupRouter() *gin.Engine {
	gin.DisableConsoleColor()
	f, _ := os.Create("gin_user.log")
	log.SetOutput(f)
	gin.DefaultWriter = io.MultiWriter(f)
	logger.Logger(f)
	r := gin.New()

	r.Use(middlewares.CORSMiddleware())

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	c := controllers.NewController()

	v1 := r.Group("/api/v1")
	{
		userGroup := v1.Group("/user")
		{
			userGroup.POST("/token", c.GenerateToken)
			userGroup.POST("/register", c.RegisterUser)
		}
	}
	return r
}
