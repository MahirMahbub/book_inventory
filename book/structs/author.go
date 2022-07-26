package structs

type CreateAuthorInput struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name"`
	Description string `json:"description"`
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

type AuthorPaginatedResponse struct {
	Data struct {
		TotalRecord int `json:"total_record"`
		TotalPage   int `json:"total_page"`
		Records     []struct {
			AuthorBase
			Url string `json:"url"`
		} `json:"records"`
		Offset   int `json:"offset"`
		Limit    int `json:"limit"`
		Page     int `json:"page"`
		PrevPage int `json:"prev_page"`
		NextPage int `json:"next_page"`
	} `json:"data"`
}

type AuthorBase struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
