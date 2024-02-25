package responses

type LoginSuccessResponse struct {
	MobileNumber string `json:"mobile" binding:"required"`
	LoginToken   string `json:"token" binding:"required"`
}
