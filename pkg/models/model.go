package models

type TesingServer struct {
	Message string `json:"message" binding:"required"`
}

type RegisterUserRequest struct {
	MobileNumber string `json:"mobile" binding:"required"`
	Otp          string `json:"otp"`
	Password     string `json:"password"`
}

type LoginSuccessResponse struct {
	MobileNumber string `json:"mobile" binding:"required"`
	LoginToken   string `json:"token" binding:"required"`
}
