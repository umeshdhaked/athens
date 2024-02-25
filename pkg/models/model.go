package models

type TesingServer struct {
	Message string `json:"message" binding:"required"`
}
