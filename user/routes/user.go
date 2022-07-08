package routes

import (
	"github.com/gin-gonic/gin"
	"go_practice/user/controllers"
	"go_practice/user/middlewares"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	c := controllers.NewController()

	v1 := r.Group("/api/v1")
	{
		userGroup := v1.Group("/user")
		{
			userGroup.POST("/token", c.GenerateToken)
			userGroup.POST("/register", c.RegisterUser)
		}
		secured := v1.Group("/secured").Use(middlewares.Auth())
		{
			secured.GET("/ping", c.Ping)
		}
	}
	return r
}
