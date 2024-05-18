package models

const (
	TableCreditAudits = "credit_audits"
)

type CreditAudits struct {
	ID             int64   `json:"id"`
	Category       string  `json:"category"`
	SubCategory    string  `json:"sub_category"`
	DeductedAmount float64 `json:"deducted_amount"`
	AddedAmount    float64 `json:"added_amount"`
	CreditId       int64   `json:"credit_id"`
	UserId         int64   `json:"user_id"`
	PaymentOrderId string  `json:"payment_order_id"`
	BaseModel
}

func (c *CreditAudits) TableName() string {
	return TableCreditAudits
}

func (c *CreditAudits) GetID() int64 {
	return c.ID
}

func (c *CreditAudits) SetID(id int64) {
	c.ID = id
}
