package structs

type CreateBookInput struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

type UpdateBookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}
type BooksPaginatedResponse struct {
	Data struct {
		TotalRecord int `json:"total_record"`
		TotalPage   int `json:"total_page"`
		Records     []struct {
			Id     int    `json:"id"`
			Title  string `json:"title"`
			Author string `json:"author"`
		} `json:"records"`
		Offset   int `json:"offset"`
		Limit    int `json:"limit"`
		Page     int `json:"page"`
		PrevPage int `json:"prev_page"`
		NextPage int `json:"next_page"`
	} `json:"data"`
}

type BookResponse struct {
	Data struct {
		Id     int    `json:"id"`
		Title  string `json:"title"`
		Author string `json:"author"`
	} `json:"data"`
}

type BookDeleteResponse struct {
	Data bool `json:"data"`
}
