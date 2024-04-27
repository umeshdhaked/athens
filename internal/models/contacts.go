package models

const (
	TableContacts = "contacts"
)

const (

	// Columns
	ColumnContactsID         = "id"
	ColumnContactsName       = "name"
	ColumnContactsMobile     = "mobile"
	ColumnContactsEmail      = "email"
	ColumnContactsGroupName  = "group_name"
	ColumnContactsAdditional = "additional"
)

type Contacts struct {
	ID         string                 `json:"id" dynamodbav:"id"`
	Name       string                 `json:"name" dynamodbav:"name"`
	Mobile     string                 `json:"mobile" dynamodbav:"mobile"`
	Email      string                 `json:"email" dynamodbav:"email"`
	GroupName  string                 `json:"group_name" dynamodbav:"group_name"`
	Additional map[string]interface{} `json:"additional" dynamodbav:"additional"`
	BaseModel
}
