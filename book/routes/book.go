package routes

import (
	"github.com/gin-gonic/gin"
	"go_practice/book/controllers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	c := controllers.NewController()

	v1 := r.Group("/api/v1")
	{
		bookGroup := v1.Group("/books")
		{
			bookGroup.GET(":id", c.FindBook)
			bookGroup.GET("", c.FindBooks)
			bookGroup.POST("", c.CreateBook)
			bookGroup.DELETE(":id", c.DeleteBook)
			bookGroup.PATCH(":id", c.UpdateBook)
		}
		//r.GET("/books", controllers.FindBooks)
		//r.POST("/books", controllers.CreateBook)
		//r.GET("/books/:id", controllers.FindBook)
		//r.PATCH("/books/:id", controllers.UpdateBook)
		//r.DELETE("/books/:id", controllers.DeleteBook)
	}
	return r
}
