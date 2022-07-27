package controllers

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// GetElasticInfo godoc
// @Summary      Get Elastic Info
// @Description  get elastic details
// @Tags         elastic
// @Accept       json
// @Produce      json
// @Success      200  {object}  structs.ElasticJsonResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      401  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /elastic/info [get]
// @Security BearerAuth
func (c *Controller) GetElasticInfo(context *gin.Context) {
	es := context.MustGet("elastic").(*elasticsearch.Client)
	var r map[string]interface{}
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	context.JSON(http.StatusOK, gin.H{"data": r})
}
