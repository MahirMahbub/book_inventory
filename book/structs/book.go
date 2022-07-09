package structs

type CreateBookInput struct {
	Title       string `json:"title" binding:"required"`
	AuthorId    uint   `json:"authorId" binding:"required"`
	Description string `json:"description"`
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
			Id    int    `json:"id"`
			Title string `json:"title"`
			Url   string `json:"url"`
		} `json:"records"`
		Offset   int `json:"offset"`
		Limit    int `json:"limit"`
		Page     int `json:"page"`
		PrevPage int `json:"prev_page"`
		NextPage int `json:"next_page"`
	} `json:"data"`
}

//type BookResponse struct {
//	Data struct {
//		Id     int    `json:"id"`
//		Title  string `json:"title"`
//		Author string `json:"author"`
//		UserId uint   `json:"userId"`
//	} `json:"data"`
//}

type BookResponse struct {
	ID          uint                  `json:"id"`
	Title       string                `json:"title"`
	UserID      uint                  `json:"userId"`
	Description string                `json:"description"`
	Authors     []AuthorBasicResponse `json:"authors"`
}

type BookUpdateResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	UserID      uint   `json:"userId"`
	Description string `json:"description"`
}

type HyperBookResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
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
