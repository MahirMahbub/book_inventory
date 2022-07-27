package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	es7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"go_practice/book/auth"
	"go_practice/book/logger"
	"go_practice/book/models"
	"go_practice/book/structs"
	"go_practice/book/utils"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// CreateAuthor godoc
// @Summary      Add an author
// @Description  post author by admin
// @Tags         admin/authors
// @Accept       json
// @Produce      json
// @Param        input  body  structs.CreateAuthorInput  true  "Add authors"
// @Success      201  {object}  structs.AuthorResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      401  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /admin/authors [post]
// @Security BearerAuth
func (c *Controller) CreateAuthor(context *gin.Context) {
	var input structs.CreateAuthorInput

	if err := context.ShouldBindJSON(&input); err != nil {
		utils.BaseErrorResponse(context, http.StatusBadRequest, err, logger.INFO)
		return
	}

	tokenString := context.GetHeader("Authorization")
	err, _ := auth.ValidateToken(tokenString)
	if err != nil {
		utils.BaseErrorResponse(context, http.StatusUnauthorized, err, logger.INFO)
		return
	}
	var author models.Author
	author = models.Author{FirstName: input.FirstName, LastName: input.LastName, Description: input.Description}

	if err := models.DB.Create(&author).Error; err != nil {
		utils.CustomErrorResponse(context, http.StatusForbidden, "author is not created", err, logger.ERROR)
		return
	}
	authorResponse := utils.CreateAuthorObjectResponse(author)
	context.JSON(http.StatusCreated, gin.H{"data": authorResponse})
}

//FindAuthors godoc
//@Summary      Show Authors
//@Description  get authors
//@Tags         authors
//@Accept       json
//@Produce      json
//@Param        page   query  int  false "paginate" Format(int)
//@Param        limit   query  int  false "paginate" Format(int)
//@Param        search   query  string  false "name searching" Format(string)
//@Success      200  {object}  structs.AuthorPaginatedResponse
//@Failure      400  {object}  structs.ErrorResponse
//@Failure      401  {object}  structs.ErrorResponse
//@Failure      403  {object}  structs.ErrorResponse
//@Failure      404  {object}  structs.ErrorResponse
//@Failure      500  {object}  structs.ErrorResponse
//@Router       /authors [get]
//@Security BearerAuth
func (c *Controller) FindAuthors(context *gin.Context) {

	//var books models.Author
	var err error
	var page, limit int
	var search string
	tokenString := context.GetHeader("Authorization")
	err, _ = auth.ValidateToken(tokenString)
	if err != nil {
		utils.BaseErrorResponse(context, http.StatusUnauthorized, err, logger.INFO)
		return
	}

	if page, err = strconv.Atoi(context.DefaultQuery("page", "1")); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "invalid 'page' param value type, Integer expected", err, logger.INFO)
		return
	}

	if limit, err = strconv.Atoi(context.DefaultQuery("limit", "10")); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "invalid 'limit' param value type, Integer expected", err, logger.INFO)
		return
	}
	var buf bytes.Buffer

	search = context.DefaultQuery("search", "")

	from := (page - 1) * limit
	//fmt.Println(search, from, limit)
	query := map[string]interface{}{
		"from": from,
		"size": limit,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					map[string]interface{}{
						"match": map[string]interface{}{
							"first_name": map[string]interface{}{
								"query":         search,
								"fuzziness":     "AUTO",
								"prefix_length": 1,
							},
						},
					},
					map[string]interface{}{
						"match": map[string]interface{}{
							"last_name": map[string]interface{}{
								"query":         search,
								"fuzziness":     "AUTO",
								"prefix_length": 1,
							},
						},
					},
				},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	es := context.MustGet("elastic").(*es7.Client)
	var r map[string]interface{}
	res, err := es.Search(
		es.Search.WithBody(&buf),
		es.Search.WithPretty(),
		es.Search.WithIndex("authors"),
		es.Search.WithTrackTotalHits(false),
		es.Search.WithFilterPath("hits.hits._source.last_name",
			"hits.hits._source.first_name",
			"hits.hits._source.id"),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	var authorsData []structs.AuthorBase
	authorsData = []structs.AuthorBase{}
	if len(r) > 0 {
		dataList := r["hits"].(map[string]interface{})["hits"].([]interface{})

		for i := 0; i < len(dataList); i++ {
			var authorDetails structs.AuthorBase
			sources := dataList[i].(map[string]interface{})["_source"].(map[string]interface{})
			//fmt.Println(sources)
			authorDetails.ID = uint(sources["id"].(float64))
			authorDetails.FirstName = sources["first_name"].(string)
			authorDetails.LastName = sources["last_name"].(string)
			authorsData = append(authorsData, authorDetails)
			//fmt.Println(authorDetails.ID, authorDetails.FirstName, authorDetails.LastName)
		}
	}
	authorStructData := utils.CreateHyperAuthorResponses(context, authorsData)

	paginatedResponse := utils.CreateHyperPaginatedAuthorResponses(page, limit, authorStructData)

	context.JSON(
		http.StatusOK,
		gin.H{"data": paginatedResponse},
	)
}

// FindAuthor godoc
// @Summary      Show Author Details
// @Description  get author
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        id  path  int true "Author ID" Format(int)
// @Success      200  {object}  structs.AuthorAPIResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      401  {object}  structs.ErrorResponse
// @Failure      403  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /authors/{id} [get]
// @Security BearerAuth
func (c *Controller) FindAuthor(context *gin.Context) {
	var author models.Author
	var id uint64
	var err error
	tokenString := context.GetHeader("Authorization")
	err, _ = auth.ValidateToken(tokenString)
	if err != nil {
		utils.BaseErrorResponse(context, http.StatusUnauthorized, err, logger.INFO)
		return
	}

	if id, err = strconv.ParseUint(context.Param("id"), 10, 32); err != nil {
		utils.CustomErrorResponse(context, http.StatusBadRequest, "invalid 'id' param value type, Integer expected", err, logger.INFO)
		return
	}

	if err := author.GetAuthorByID(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.BaseErrorResponse(context, http.StatusNotFound, err, logger.INFO)
			return
		}
		utils.CustomErrorResponse(context, http.StatusForbidden, "operation is not allowed", err, logger.ERROR)
		return
	}
	authorResponse := utils.CreateAuthorObjectResponse(author)
	context.JSON(http.StatusOK, gin.H{"data": authorResponse})
}
