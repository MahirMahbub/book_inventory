package routes

import (
	"github.com/gin-gonic/gin"
	"go_practice/book/controllers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/books", controllers.FindBooks)
	r.POST("/books", controllers.CreateBook)
	r.GET("/books/:id", controllers.FindBook)
	r.PATCH("/books/:id", controllers.UpdateBook)
	r.DELETE("/books/:id", controllers.DeleteBook)

	return r
}
