package structs

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type TokenResponse struct {
	Token string `json:"token"`
}

type UserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserResponse struct {
	UserId   uint   `json:"userId"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
