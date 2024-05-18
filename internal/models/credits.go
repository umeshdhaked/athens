package models

const (
	TableCredits = "credits"
)

type Credits struct {
	ID          int64   `json:"id"`
	UserId      int64   `json:"user_id"`
	Balance     float64 `json:"balance"`
	BalanceLeft float64 `json:"balance_left"`
	BaseModel
}

func (c *Credits) TableName() string {
	return TableCredits
}

func (c *Credits) GetID() int64 {
	return c.ID
}

func (c *Credits) SetID(id int64) {
	c.ID = id
}
