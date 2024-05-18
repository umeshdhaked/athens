package models

const (
	TableContacts = "contacts"
)

type Contacts struct {
	ID         int64                  `json:"id" dynamodbav:"id"`
	Name       string                 `json:"name" dynamodbav:"name"`
	Mobile     string                 `json:"mobile" dynamodbav:"mobile"`
	Email      string                 `json:"email" dynamodbav:"email"`
	GroupName  string                 `json:"group_name" dynamodbav:"group_name"`
	Additional map[string]interface{} `json:"additional" dynamodbav:"additional"`
	BaseModel
}

func (c *Contacts) TableName() string {
	return TableContacts
}

func (c *Contacts) GetID() int64 {
	return c.ID
}

func (c *Contacts) SetID(id int64) {
	c.ID = id
}
