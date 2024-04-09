package dtos

type LoginSuccessResponse struct {
	MobileNumber string `json:"mobile" binding:"required"`
	LoginToken   string `json:"authToken" binding:"required"`
}
