package main

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "go_practice/book/docs"
	"go_practice/book/models"
	"go_practice/book/routes"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:5001
// @securityDefinitions.basic  BasicAuth
func main() {
	//r := gin.Default()

	//r.GET("/", func(context *gin.Context) {
	//	context.JSON(http.StatusOK, gin.H{"data": "hello world"})
	//})
	models.ConnectDatabase()
	router := routes.SetupRouter()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := router.Run(":5001")
	if err != nil {
		return
	}
	//r.GET("/books", controllers.FindBooks)
	//r.POST("/books", controllers.CreateBook)
	//r.GET("/books/:id", controllers.FindBook)
	//r.PATCH("/books/:id", controllers.UpdateBook)
	//r.DELETE("/books/:id", controllers.DeleteBook)
	//r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//
	//err := r.Run(":5001")
	//if err != nil {
	//	return
	//}
}
