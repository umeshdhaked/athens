package models

type Testing struct {
	Message string `json:"message" binding:"required"`
}
