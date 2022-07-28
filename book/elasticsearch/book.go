package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"log"
)

func GetPaginatedBookSearch(context *gin.Context, from int, limit int, search string, buf bytes.Buffer, err error) (map[string]interface{}, error) {
	query := map[string]interface{}{
		"collapse": map[string]interface{}{
			"field": "title.keyword",
		},
		"from": from,
		"size": limit,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					map[string]interface{}{
						"match": map[string]interface{}{
							"author_first_name": map[string]interface{}{
								"query":         search,
								"fuzziness":     "AUTO",
								"prefix_length": 1,
							},
						},
					},
					map[string]interface{}{
						"match": map[string]interface{}{
							"author_last_name": map[string]interface{}{
								"query":         search,
								"fuzziness":     "AUTO",
								"prefix_length": 1,
							},
						},
					},
					map[string]interface{}{
						"match": map[string]interface{}{
							"title": map[string]interface{}{
								"query":          search,
								"fuzziness":      "AUTO",
								"prefix_length":  1,
								"max_expansions": 200,
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

	es := context.MustGet("elastic").(*elasticsearch.Client)
	var resposeBind map[string]interface{}
	res, err := es.Search(
		es.Search.WithBody(&buf),
		es.Search.WithPretty(),
		es.Search.WithIndex("books"),
		es.Search.WithTrackTotalHits(false),
		es.Search.WithFilterPath("hits.hits._source.title",
			"hits.hits._source.book_id"),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	if err := json.NewDecoder(res.Body).Decode(&resposeBind); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	fmt.Println(resposeBind)
	return resposeBind, err
}
