package models

const (
	TableGroup = "group"
)

const (
	ColumnGroupID          = "id"
	ColumnGroupName        = "name"
	ColumnGroupUserID      = "userID"
	ColumnGroupColumnNames = "columnNames"
)

type Group struct {
	ID          string   `json:"id" dynamodbav:"id"`
	Name        string   `json:"name" dynamodbav:"name"`
	UserID      string   `json:"user_id" dynamodbav:"userID"`
	ColumnNames []string `json:"column_names" dynamodbav:"columnNames"`
	BaseModel            // todo implement auto basemodel values for insert and update operations
}
