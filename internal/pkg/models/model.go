package models

type TesingRequest struct {
	Message string `json:"message" binding:"required"`
	Jwt     string `json:"jwt" binding:"required"`
}

type TestingResponse struct {
	Message string `json:"message" binding:"required"`
}
