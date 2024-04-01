package requests

type RegisterUserRequest struct {
	MobileNumber string `json:"mobile" binding:"required"`
	Otp          string `json:"otp"`
	Password     string `json:"password"`
}

type PromoUserRequest struct {
	MobileNumber         string `json:"mobile"`
	Comment              string `json:"comment"`
	Is_already_contacted string `json:"is_already_contacted"`
}
