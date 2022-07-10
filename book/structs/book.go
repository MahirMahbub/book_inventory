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
		TotalRecord int `json:"total_record"`
		TotalPage   int `json:"total_page"`
		Records     []struct {
			BookBase
			Url string `json:"url"`
		} `json:"records"`
		Offset   int `json:"offset"`
		Limit    int `json:"limit"`
		Page     int `json:"page"`
		PrevPage int `json:"prev_page"`
		NextPage int `json:"next_page"`
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

type BookUpdateResponse struct {
	BookBaseResponse
}

type HyperBookResponse struct {
	BookBase
	Url string `json:"url"`
}
type AuthorResponse struct {
	ID          uint   `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Description string `json:"description"`
}

type AuthorBasicResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type BookDeleteResponse struct {
	Data bool `json:"data"`
}
