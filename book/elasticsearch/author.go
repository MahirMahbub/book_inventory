package elasticsearch

import (
	"bytes"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"log"
)

func GetPaginatedAuthorSearch(context *gin.Context, from int, limit int, search string, buf bytes.Buffer, err error) (map[string]interface{}, error) {
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

	es := context.MustGet("elastic").(*elasticsearch.Client)
	var resposeBind map[string]interface{}
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
	if err := json.NewDecoder(res.Body).Decode(&resposeBind); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	return resposeBind, err
}
