package models

const (
	TableGroup = "Group"
)

type Group struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	UserID string `json:"user_id"`
	BaseModel
}
