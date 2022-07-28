package routes

import (
	"github.com/gin-gonic/gin"
	"go_practice/book/controllers"
)

func ElasticRoute(v1 *gin.RouterGroup, c *controllers.Controller) {
	elasticGroup := v1.Group("/elastic")
	{
		elasticGroup.GET("/info", c.GetElasticInfo)
		elasticGroup.GET("/authors", c.FindAuthors)
		elasticGroup.GET("/books", c.FindBooks)
	}
}
