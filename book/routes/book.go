package routes

import (
	"github.com/gin-gonic/gin"
	"go_practice/book/controllers"
	"go_practice/book/middlewares"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())
	c := controllers.NewController()

	v1 := r.Group("/api/v1")
	{
		bookGroup := v1.Group("/books").Use(middlewares.Auth())
		{
			bookGroup.GET(":id", c.FindBook)
			bookGroup.GET("", c.FindBooks)
			bookGroup.POST("", c.CreateBook)
			bookGroup.DELETE(":id", c.DeleteBook)
			bookGroup.PATCH(":id", c.UpdateBook)
		}
	}
	return r
}
