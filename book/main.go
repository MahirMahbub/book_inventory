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
// @description     This is a sample book server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  bsse0807@iit.du.ac.bd

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description					Description for what is this security definition being used

// @host      localhost:5001
// @BasePath  /api/v1
func main() {
	models.ConnectDatabase()
	router := routes.SetupRouter()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := router.Run(":5001")
	if err != nil {
		return
	}
}
