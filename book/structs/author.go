package structs

type CreateAuthorInput struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name"`
	Description string `json:"description"`
}

type AuthorBase struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type HyperAuthorResponse struct {
	AuthorBase
	Url string `json:"url"`
}

type AuthorResponse struct {
	AuthorBase
	Description string              `json:"description"`
	Books       []HyperBookResponse `json:"books"`
}

type AuthorAPIResponse struct {
	Data struct {
		AuthorResponse
	} `json:"data"`
}

type AuthorBasicResponse struct {
	AuthorBase
	Name string `json:"name"`
}

type AuthorElasticPaginatedResponse struct {
	Data struct {
		Records  []HyperAuthorResponse `json:"records"`
		Limit    int                   `json:"limit"`
		Page     int                   `json:"page"`
		PrevPage int                   `json:"prev_page"`
		NextPage int                   `json:"next_page"`
	} `json:"data"`
}

type AuthorPaginated struct {
	Records  []HyperAuthorResponse `json:"records"`
	Limit    int                   `json:"limit"`
	Page     int                   `json:"page"`
	PrevPage int                   `json:"prev_page"`
	NextPage int                   `json:"next_page"`
}
