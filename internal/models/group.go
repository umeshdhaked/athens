package models

const (
	TableGroup = "Group"
)

const (
	ColumnID          = "ID"
	ColumnName        = "Name"
	ColumnUserID      = "UserID"
	ColumnColumnNames = "ColumnNames"
)

type Group struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	UserID      string `json:"user_id"`
	ColumnNames string `json:"column_names"`
	BaseModel
}
