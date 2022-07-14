package structs

type TokenRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type TokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type PasswordChangeTokenRequest struct {
	Email string `json:"email" binding:"required"`
}

type UserRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username " binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type UserResponse struct {
	UserId   uint   `json:"userId"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type ActiveUserResponse struct {
	UserId   uint   `json:"userId"`
	Email    string `json:"email"`
	Username string `json:"username"`
	IsActive bool   `json:"isActive"`
}

type ChangePasswordRequest struct {
	Password string `json:"password" binding:"required"`
	Confirm  string `json:"confirm" binding:"required"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
