package dto

type ResgisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" bindibg:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type ResgisterResponse struct {
	Message string `json:"message"`
	UserID  uint   `json:"user_id"`
}

type ErrorResponse struct {
	Error string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}