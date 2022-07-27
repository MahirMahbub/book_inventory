package structs

type CreateBookInput struct {
	Title       string `json:"title" binding:"required"`
	AuthorIDs   []uint `json:"authorIds" binding:"required"`
	Description string `json:"description"`
}
type T struct {
	As []int `json:"as"`
}

type UpdateBookInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}
type BooksPaginatedResponse struct {
	Data struct {
		TotalRecord int `json:"totalRecord"`
		TotalPage   int `json:"totalPage"`
		Records     []struct {
			BookBase
			Url string `json:"url"`
		} `json:"records"`
		Offset   int `json:"offset"`
		Limit    int `json:"limit"`
		Page     int `json:"page"`
		PrevPage int `json:"prevPage"`
		NextPage int `json:"nextPage"`
	} `json:"data"`
}

type BookBase struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}
type BookBaseResponse struct {
	BookBase
	UserID      uint   `json:"userId"`
	Description string `json:"description"`
}

type BookResponse struct {
	BookBaseResponse
	Authors []AuthorBasicResponse `json:"authors"`
}

type BookAPIResponse struct {
	Data struct {
		BookResponse
	} `json:"data"`
}

type BookUpdateResponse struct {
	BookBaseResponse
}

type HyperBookResponse struct {
	BookBase
	Url string `json:"url"`
}

type BookDeleteResponse struct {
	Data bool `json:"data"`
}

type ElasticJsonResponse struct {
	Data any `json:"data"`
}
